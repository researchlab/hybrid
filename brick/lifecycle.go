package brick

// Lifecycle is:
//	1. New an object
//	2. Call AfterNew on an object when it has been created
//	3. Call Init on an object thought dependency after all objects has been created

// AfterNewProcessor called after an object has been created before init
type AfterNewProcessor interface {
	AfterNew()
}

// Initializer init a service
type Initializer interface {
	Init() error
}

// Disposable dispose a service that release all resources which used
type Disposable interface {
	Dispose() error
}
