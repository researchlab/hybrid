package brick

import (
	"fmt"
	"log"
	"reflect"
)

// ObjectState  alias type of object state
type ObjectState uint32

// support 3 object states
const (
	Empty ObjectState = 0
	New               = 1
	Init              = 2
)

// ObjectRef store the information of a object instance
type ObjectRef struct {
	objType reflect.Type
	key     string
	factory Factory
	state   ObjectState
	obj     interface{}
}

func key(objType reflect.Type, name string) string {
	key := name
	if name == "" {
		key = objType.Name()
	}
	return key
}

func newObjectRef(objType reflect.Type, name string, factory Factory) *ObjectRef {

	return &ObjectRef{objType: objType, key: key(objType, name), factory: factory, state: Empty}
}

//Key unique key of object
func (p *ObjectRef) Key() string {
	return p.key
}

//Target return object itself
func (p *ObjectRef) Target() interface{} {
	return p.obj
}

//Type reflect type of object
func (p *ObjectRef) Type() reflect.Type {
	return p.objType
}

func (p *ObjectRef) new(obj interface{}) {
	p.obj = obj
	p.state = New
}

func (p *ObjectRef) init(obj interface{}) {
	p.obj = obj
	p.state = Init
}

// ContainerAware provide the container to the service
type ContainerAware interface {
	// SetContainer set the container to the service
	SetContainer(c *Container)
}

// Container manage all services
type Container struct {
	objs  map[string]*ObjectRef
	inits []interface{}
}

// NewContainer new object of container
func NewContainer() *Container {
	return &Container{objs: map[string]*ObjectRef{}, inits: []interface{}{}}
}

// Add obj type in the container
func (p *Container) Add(obj interface{}, name string, factory Factory) {
	objType := reflect.TypeOf(obj)
	if objType.Kind() == reflect.Ptr {
		st := newObjectRef(objType, name, factory)
		p.objs[st.Key()] = st
	} else {
		error := fmt.Sprintf("Contianer.Add's arg obj %s must be a pointer,but is %v", objType.Name(), objType.Kind())
		fmt.Println(error)
		panic(error)
	}

}

// Build all objs in the container
func (p *Container) Build() {
	for _, r := range p.objs {
		switch r.state {
		case Empty:
			p.initObject(r)
		case New:
			panic(fmt.Errorf("the New state of the obj is invalid. %+v,%s", r.objType, r.Key()))
		}
	}

	for index, o := range p.inits {
		if i, ok := o.(Initializer); ok {
			if err := i.Init(); err != nil {
				fmt.Printf("init %v err: %s\n", o, err.Error())
				p.disposeObjs(index)
			}
		}
	}
}

// Dispose dispose all objs in the container
func (p *Container) Dispose() {
	p.disposeObjs(len(p.inits))
}

// GetByName get object by name
func (p *Container) GetByName(name string) interface{} {
	key := key(nil, name)
	ref := p.objs[key]
	if ref == nil {
		return nil
	}

	if ref.state != Init {
		panic(fmt.Errorf("uninit obj: %+v,%s", ref.objType, ref.Key()))
	}

	return ref.Target()
}

// GetByType get object by reflect type of object
func (p *Container) GetByType(svcType reflect.Type) interface{} {
	return p.GetByName(svcType.Name())
}

func (p *Container) initObject(r *ObjectRef) {
	factory := r.factory
	var obj interface{}
	if factory != nil {
		obj = factory.New()
	} else {
		obj = reflect.New(r.Type().Elem()).Interface()
	}

	r.new(obj)
	p.inject(r)
	if c, ok := obj.(ContainerAware); ok {
		c.SetContainer(p)
	}
	if _, ok := obj.(Initializer); ok {
		p.inits = append(p.inits, obj)
	}
	r.init(obj)
	afterNewProcessor, ok := obj.(AfterNewProcessor)
	if ok {
		afterNewProcessor.AfterNew()
	}
}

func (p *Container) inject(r *ObjectRef) {
	t := r.Type().Elem()
	count := t.NumField()
	for i := 0; i < count; i++ {
		f := t.Field(i)

		injectObjName := f.Tag.Get("inject")

		if injectObjName != "" {
			injectR := p.objs[injectObjName]
			if injectR == nil {
				panic(fmt.Errorf("could not found object %s,defines in %v.%s", injectObjName, t, f.Name))
			}
			switch injectR.state {
			case New:
				panic(fmt.Errorf("cycle dependency %s,defines in %v.%s", injectObjName, t, f.Name))
			case Empty:
				p.initObject(injectR)
			}

			v := reflect.ValueOf(r.obj)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}

			fv := v.Field(i)
			//			fmt.Printf("set %v.%s\n", t, f.Name)
			fv.Set(reflect.ValueOf(injectR.Target()))
		}
	}
}

func (p *Container) disposeObjs(count int) {
	for i := count - 1; i >= 0; i-- {
		if d, ok := p.inits[i].(Disposable); ok {
			if err := d.Dispose(); err != nil {
				log.Println(err.Error())
			}
		}
	}
}
