# hybrid

- hybrid is an API framework of mixture RESTFUL API and RPC API 

## feature 
- Friendly support RESTFUL API Style
- Flexible support RPC API Style

## dependencies

```
vendor
├── github.com
│   ├── go-resty
│   ├── go-sql-driver
│   ├── gopherjs
│   ├── jinzhu
│   ├── jtolds
│   ├── mattn
│   ├── pressly
│   ├── sirupsen
│   ├── smartystreets
│   └── toolkits
└── golang.org
    └── x
        ├── net
        └── sys
```

## Example 

- register modules to the global container [stu](https://github.com/researchlab/hybrid/tree/master/examples/basic)

```
func main() {
	stopSigs := make(chan os.Signal, 1)
	signal.Notify(stopSigs, syscall.SIGINT, syscall.SIGTERM)
	container := brick.NewContainer()
	container.Add(&brick.JSONConfigService{}, "config", brick.FactoryFunc(func() interface{} {
		return brick.NewJSONConfigService(configPath)
	}))
	container.Add(&model.Models{}, "Models", nil)
	container.Add(&mysql.MySQLService{}, "DB", nil)
	container.Add(&router.HttpService{}, "HttpService", nil)
	container.Add(&rest.RestController{}, "RestController", nil)
	container.Add(&stu.StuService{}, "StuService", nil)
	container.Build()
	defer container.Dispose()
	select {
	case <-stopSigs:
		log.Println("service has been stoped.")
	}
}
```

- router support 

```
r.Route("/api/objs/:class", func(r chi.Router) {
		r.Get("/", p.Controller.List)
		r.Post("/", p.Controller.Create)
		r.Put("/", p.Controller.Update)
		r.Route("/:id", func(r chi.Router) {
			r.Get("/", p.Controller.Get)
			r.Delete("/", p.Controller.Remove)
			r.Post("/", p.Controller.InvokeServiceFunc())
		})
	})
```

## Resourcese CURD 

### RESTFUL API Style Support 

- HTTP POST  / Create resourcese
```
➜  ~ curl -XPOST -H "Content-Type: application/json" http://127.0.0.1:9999/api/objs/Stu -d'{"name":"mike","address":"shanghai.china.002", "sex":"male"}'
{"ID":1,"CreatedBy":"","CreatedAt":"2018-12-11T16:12:27.756548+08:00","UpdatedBy":"","UpdatedAt":"2018-12-11T16:12:27.756548+08:00","DeletedBy":"","DeletedAt":null,"Name":"mike","Address":"shanghai.china.002","Sex":"male"}
```

- HTTP PUT / Update resourcese
```
➜  ~ curl -XPUT -H "Content-Type: application/json" http://127.0.0.1:9999/api/objs/Stu -d'{"ID":1,"CreatedBy":"","CreatedAt":"2018-12-11T16:12:27.756548+08:00","UpdatedBy":"","UpdatedAt":"2018-12-11T16:12:27.756548+08:00","DeletedBy":"","DeletedAt":null,"Name":"mike.lee","Address":"shanghai.china.002","Sex":"male"}'
```

- HTTP GET / Get One 
```
➜  ~ curl -XGET -H "Content-Type: application/json" http://127.0.0.1:9999/api/objs/Stu/1
{"ID":1,"CreatedBy":"","CreatedAt":"2018-12-11T16:12:28+08:00","UpdatedBy":"","UpdatedAt":"2018-12-11T16:15:18+08:00","DeletedBy":"","DeletedAt":null,"Name":"mike.lee","Address":"shanghai.china.002","Sex":"male"}
```

- HTTP GET / Get All
```
➜  ~ curl -XGET -H "Content-Type: application/json" http://127.0.0.1:9999/api/objs/Stu
{"data":[{"ID":1,"CreatedBy":"","CreatedAt":"2018-12-11T16:12:28+08:00","UpdatedBy":"","UpdatedAt":"2018-12-11T16:15:18+08:00","DeletedBy":"","DeletedAt":null,"Name":"mike.lee","Address":"shanghai.china.002","Sex":"male"},{"ID":2,"CreatedBy":"","CreatedAt":"2018-12-11T16:16:17+08:00","UpdatedBy":"","UpdatedAt":"2018-12-11T16:16:17+08:00","DeletedBy":"","DeletedAt":null,"Name":"alex","Address":"shanghai.china.002","Sex":"female"}],"page":0,"pageCount":1,"pageSize":10}
```

- HTTP DELETE / DELETE resourcese
```
➜  ~ curl -XDELETE -H "Content-Type: application/json" http://127.0.0.1:9999/api/objs/Stu/2
```

### RPC API Style Support 

- request StuService/SayHi Method 
```
➜  ~ curl -XPOST -H "Content-Type: application/json" http://127.0.0.1:9999/api/objs/StuService/SayHi -d'["mike"]'
["Hi, Mr.mike.\n",null]

➜  ~ curl -XPOST -H "Content-Type: application/json" http://127.0.0.1:9999/api/objs/StuService/SayHi -d'["alex"]'
["Hi, Mrs.alex.\n",null]

➜  ~ curl -XPOST -H "Content-Type: application/json" http://127.0.0.1:9999/api/objs/StuService/SayHi -d'["mike"]'
["no one named mike found.",null]
```

- RPC API Paybody DESC 
  - support array paybody and must Keep in the same order as the function parameters.
	- example. func(id uint, name string), must request with POST Protocol AND postBody like [10, "name"].
