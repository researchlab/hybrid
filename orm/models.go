package orm

//ModelDescriptor db model descriptor
type ModelDescriptor struct {
	Type     interface{}
	New      func() interface{}
	NewSlice func() interface{}
}

// ModelRegistry Register all orm models
type ModelRegistry interface {
	Put(name string, model *ModelDescriptor)
	Get(name string) *ModelDescriptor
	Models() <-chan *ModelDescriptor
}

//ModelRegistryImpl model registry impelement struct
type ModelRegistryImpl struct {
	models map[string]*ModelDescriptor
}

//Put registry db model
func (p *ModelRegistryImpl) Put(name string, model *ModelDescriptor) {
	if p.models == nil {
		p.models = map[string]*ModelDescriptor{}
	}
	p.models[name] = model
}

//Get get db model by model struct name
func (p *ModelRegistryImpl) Get(name string) *ModelDescriptor {
	if p.models == nil {
		return nil
	}
	return p.models[name]
}

//Models models channels
func (p *ModelRegistryImpl) Models() <-chan *ModelDescriptor {
	ch := make(chan *ModelDescriptor)
	go func() {
		for _, m := range p.models {
			ch <- m
		}
		close(ch)
	}()
	return ch
}
