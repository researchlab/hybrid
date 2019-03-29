package brick

// Factory create a service
type Factory interface {
	// New a service instance
	New() interface{}
}

// NewFunc New return a new object
type NewFunc func() interface{}

// FactoryFunc return a factory func
func FactoryFunc(f NewFunc) Factory {
	return &FactoryFuncWrap{f: f}
}

// FactoryFuncWrap wrap a factory
type FactoryFuncWrap struct {
	f NewFunc
}

// New New return a factory object of func
func (p *FactoryFuncWrap) New() interface{} {
	return p.f()
}
