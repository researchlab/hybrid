package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/pressly/chi"
	"github.com/researchlab/hybrid/brick"
	"github.com/researchlab/hybrid/net/http/rest/render"
	"github.com/researchlab/hybrid/orm"
)

func handlePanic(w http.ResponseWriter, req *http.Request) {
	if r := recover(); r != nil {
		fmt.Printf("panic: %v\n", r)
		debug.PrintStack()
		render.Status(req, 500)
		render.JSON(w, req, r)
	}
}

// RestController
type RestController struct {
	brick.Trigger
	DB            orm.DBService     `inject:"DB"`
	ModelRegistry orm.ModelRegistry `inject:"DB"`
	container     *brick.Container
	ormService    orm.Repository
}

func (p *RestController) Init() error {
	if p.ormService == nil {
		p.ormService = &orm.OrmService{DB: p.DB, ModelRegistry: p.ModelRegistry} //newormService(p.DB, p.ModelRegistry)
	}

	return nil
}

func (p *RestController) List(w http.ResponseWriter, req *http.Request) {
	defer handlePanic(w, req)

	/**
	  获得所有查询参数
	*/
	class := chi.URLParam(req, "class")
	where := req.URL.Query().Get("where")
	selectQuery := req.URL.Query().Get("select")

	/**
	  如果select参数不为空, 则获得要查询的字段集合
	*/
	var selectFields []string
	if selectQuery != "" {
		for _, arg := range strings.Split(selectQuery, ",") {
			selectFields = append(selectFields, arg)
		}
	}

	/**
	  如果where参数不为空, 则获得values数组
	*/
	var whereValues []interface{}
	if where != "" {
		values := req.URL.Query().Get("values")
		if values == "" {
			render.Status(req, 400)
			render.JSON(w, req, fmt.Sprintf("the values query param must be provided if the where query param is exists"))
			return
		}
		for _, arg := range strings.Split(values, ",") {
			whereValues = append(whereValues, arg)
		}
	}

	order := req.URL.Query().Get("order")
	page, err := getQueryParamInt(req, "page", 0)
	if err != nil {
		render.Status(req, 400)
		render.JSON(w, req, err.Error())
		return
	}
	pageSize, err := getQueryParamInt(req, "pageSize", 10)
	if err != nil {
		render.Status(req, 400)
		render.JSON(w, req, err.Error())
		return
	}

	data, err := p.ormService.List(class, selectFields, where, whereValues, order, page, pageSize)
	if err != nil {
		render.Status(req, 500)
		render.JSON(w, req, err.Error())
		return
	}

	render.Status(req, 200)
	render.JSON(w, req, data)

	// render.JSON(w, req, map[string]interface{}{"data": data, "page": page, "pageSize": limit, "pageCount": pageCount})
}

// get one obj. query params: assocations=a,b...
func (p *RestController) Get(w http.ResponseWriter, req *http.Request) {
	defer handlePanic(w, req)
	p.get(w, req)
}

func (p *RestController) get(w http.ResponseWriter, req *http.Request) {
	class := chi.URLParam(req, "class")
	id := chi.URLParam(req, "id")
	ass := req.URL.Query().Get("associations")

	data, err := p.ormService.Get(class, id, ass)
	if err != nil {
		render.Status(req, 500)
		render.JSON(w, req, err.Error())
		return
	}

	render.Status(req, 200)
	render.JSON(w, req, data)

}

func (p *RestController) Create(w http.ResponseWriter, req *http.Request) {
	//	p.CreateCtx(context.TODO(), w, req)
	defer handlePanic(w, req)

	class := chi.URLParam(req, "class")
	md := p.ModelRegistry.Get(class)
	if md == nil {
		render.Status(req, 400)
		render.JSON(w, req, fmt.Sprintf("class %s isn't exists", class))
		return
	}

	data := md.New()
	if err := render.Bind(req.Body, data); err != nil {
		render.Status(req, 400)
		render.JSON(w, req, err.Error())
		return
	}

	if err := p.ormService.Create(class, data); err != nil {
		render.Status(req, 500)
		render.JSON(w, req, err.Error())
		return
	}

	render.Status(req, 200)
	render.JSON(w, req, data)

	p.raise(fmt.Sprintf("%s.Create", class), data)

}

func (p *RestController) Remove(w http.ResponseWriter, req *http.Request) {
	defer handlePanic(w, req)

	class := chi.URLParam(req, "class")
	id := chi.URLParam(req, "id")

	md := p.ModelRegistry.Get(class)
	if md == nil {
		render.Status(req, 400)
		render.JSON(w, req, fmt.Sprintf("class %s isn't exists", class))
		return
	}

	data, err := p.ormService.Remove(class, id)
	if err != nil {
		render.Status(req, 500)
		render.JSON(w, req, err.Error())
		return
	}

	render.Status(req, 200)
	p.raise(fmt.Sprintf("%v.Delete", class), data)

}

func (p *RestController) Update(w http.ResponseWriter, req *http.Request) {
	defer handlePanic(w, req)

	class := chi.URLParam(req, "class")
	md := p.ModelRegistry.Get(class)
	if md == nil {
		render.Status(req, 400)
		render.JSON(w, req, fmt.Sprintf("class %s isn't exists", class))
		return
	}

	data := md.New()
	if err := render.Bind(req.Body, data); err != nil {
		render.Status(req, 400)
		render.JSON(w, req, err.Error())
		return
	}

	if err := p.ormService.Update(class, data); err != nil {
		render.Status(req, 500)
		render.JSON(w, req, err.Error())
		return
	}

	render.Status(req, 200)
	render.JSON(w, req, data)
	p.raise(fmt.Sprintf("%v.Update", class), data)

}

// InvokeServiceFunc call the  func of service
func (p *RestController) InvokeServiceFunc() func(w http.ResponseWriter, req *http.Request) {
	return func(c *brick.Container) func(w http.ResponseWriter, req *http.Request) {
		return func(w http.ResponseWriter, req *http.Request) {
			invokeService(c, w, req)
		}
	}(p.container)
}

func invokeService(c *brick.Container, w http.ResponseWriter, req *http.Request) {
	defer handlePanic(w, req)

	class := chi.URLParam(req, "class")
	methodName := chi.URLParam(req, "id")

	svc := c.GetByName(class)
	if svc == nil {
		render.Status(req, 404)
		render.JSON(w, req, fmt.Sprintf("service %s isn't exists", class))
		return
	}

	t := reflect.TypeOf(svc)
	if m, b := t.MethodByName(methodName); b != true {
		render.Status(req, 404)
		render.JSON(w, req, fmt.Sprintf("method %s.%s isn't exists.from %v", class, methodName, t))
		return
	} else {
		args := []json.RawMessage{}
		err := render.Bind(req.Body, &args)
		if err != nil && err != io.EOF {
			render.Status(req, 400)
			render.JSON(w, req, err.Error())
			return
		}

		in := []reflect.Value{reflect.ValueOf(svc)}
		funcType := m.Func.Type()
		if funcType.NumIn() > 1 {
			step := 1

			for i, arg := range args {
				argType := funcType.In(i + step)
				argValue := reflect.New(argType)
				if err := json.Unmarshal(arg, argValue.Interface()); err == nil {
					in = append(in, argValue.Elem())
				} else {
					render.Status(req, 500)
					render.JSON(w, req, err.Error())
					return
				}
			}
		}

		values := m.Func.Call(in)
		ret := []interface{}{}
		httpCode := 200
		for _, value := range values {
			rawValue := value.Interface()
			if e, ok := rawValue.(error); ok {
				httpCode = 500
				ret = append(ret, e.Error())
			} else {
				ret = append(ret, rawValue)
			}

		}
		render.Status(req, httpCode)
		render.JSON(w, req, ret)
	}
}

// InvokeServiceRawMessageFunc call the  func of service with json.RawMessage
func (p *RestController) InvokeServiceRawMessageFunc() func(w http.ResponseWriter, req *http.Request) {
	return func(c *brick.Container) func(w http.ResponseWriter, req *http.Request) {
		return func(w http.ResponseWriter, req *http.Request) {
			invokeServiceRawMessage(c, w, req)
		}
	}(p.container)
}

func invokeServiceRawMessage(c *brick.Container, w http.ResponseWriter, req *http.Request) {
	defer handlePanic(w, req)

	class := chi.URLParam(req, "class")
	methodName := chi.URLParam(req, "id")

	svc := c.GetByName(class)
	if svc == nil {
		render.Status(req, 404)
		render.JSON(w, req, fmt.Sprintf("service %s isn't exists", class))
		return
	}

	t := reflect.TypeOf(svc)
	if m, b := t.MethodByName(methodName); b != true {
		render.Status(req, 404)
		render.JSON(w, req, fmt.Sprintf("method %s.%s isn't exists.from %v", class, methodName, t))
		return

	} else {
		arg := json.RawMessage{}
		if err := render.Bind(req.Body, &arg); err != nil {
			render.Status(req, 400)
			render.JSON(w, req, err.Error())
			return
		}

		in := []reflect.Value{reflect.ValueOf(svc)}
		funcType := m.Func.Type()
		if funcType.NumIn() > 1 {
			in = append(in, reflect.ValueOf(arg))
		}

		values := m.Func.Call(in)
		ret := []interface{}{}
		httpCode := 200
		if len(values) != 2 {
			render.Status(req, httpCode)
			render.JSON(w, req, fmt.Sprintf("service must return 2 values, (string,error)"))
			return
		}

		for _, value := range values {
			rawValue := value.Interface()
			if e, ok := rawValue.(error); ok {
				httpCode = 500
				ret = append(ret, e.Error())
			} else {
				ret = append(ret, rawValue)
			}
		}
		render.Status(req, httpCode)
		render.JSON(w, req, ret)
	}
}

func (p *RestController) InvokeObj(w http.ResponseWriter, req *http.Request) {
	defer handlePanic(w, req)

	class := chi.URLParam(req, "class")
	id := chi.URLParam(req, "id")
	methodName := chi.URLParam(req, "method")

	data, err := p.ormService.Get(class, id, "")
	if err != nil {
		render.Status(req, 404)
		render.JSON(w, req, err.Error())
		return
	}

	t := reflect.TypeOf(data)
	if m, b := t.MethodByName(methodName); b != true {
		render.Status(req, 404)
		render.JSON(w, req, fmt.Sprintf("%s/%s/%s isn't exists", class, id, methodName))
		return

	} else {
		args := []interface{}{}
		if err := render.Bind(req.Body, &args); err != nil {
			render.Status(req, 400)
			render.JSON(w, req, err.Error())
			return
		}

		in := []reflect.Value{reflect.ValueOf(data)}
		for _, arg := range args {
			in = append(in, reflect.ValueOf(arg))
		}
		values := m.Func.Call(in)
		ret := []interface{}{}
		httpCode := 200
		for _, value := range values {
			rawValue := value.Interface()
			if e, ok := rawValue.(error); ok {
				httpCode = 500
				ret = append(ret, e.Error())
			} else {
				ret = append(ret, rawValue)
			}
		}
		render.Status(req, httpCode)
		render.JSON(w, req, ret)
	}
}

func (p *RestController) SetContainer(c *brick.Container) {
	//log.Printf("SetContainer:%+v\n", c)
	p.container = c
	//log.Printf("SetContainer:%+v\n", p.container)
}

func (p *RestController) raise(event string, data interface{}) {
	p.Emmit(event, data)
}
