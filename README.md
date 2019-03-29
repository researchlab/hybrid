# hybrid [![GitHub release](https://img.shields.io/github/release/researchlab/hybrid.svg)](https://github.com/researchlab/hybrid/releases) [![Go Report Card](https://goreportcard.com/badge/github.com/researchlab/hybrid)](https://goreportcard.com/report/github.com/researchlab/hybrid) [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/researchlab/hybrid) 

- hybrid is an API framework of mixture RESTFUL API and RPC API 

## feature 
- Friendly support RESTFUL API Style
- Flexible support RPC API Style

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
r.Route("/api/v1/:class", func(r chi.Router) {
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
➜  ~ curl -XPOST -H "Content-Type: application/json" http://127.0.0.1:9999/api/v1/Stu -d'{"name":"mike","address":"shanghai.china.002", "sex":"male"}'
{"ID":1,"CreatedAt":"2019-03-28T10:20:06.015921+08:00","UpdatedAt":"2019-03-28T10:20:06.015921+08:00","DeletedAt":null,"Name":"mike","Address":"shanghai.china.002","Sex":"male"} 
```

- HTTP PUT / Update resourcese
```
➜  ~ curl -XPUT -H "Content-Type: application/json" http://127.0.0.1:9999/api/v1/Stu -d'{"ID":1,"CreatedAt":"2019-03-28T10:20:06+08:00","UpdatedAt":"2019-03-28T10:20:06+08:00","DeletedAt":null,"Name":"mike","Address":"shanghai.china.001","Sex":"male"}'
{"ID":1,"CreatedAt":"2019-03-28T10:20:06+08:00","UpdatedAt":"2019-03-28T10:35:43.756395+08:00","DeletedAt":null,"Name":"mike","Address":"shanghai.china.001","Sex":"male"}
```

- HTTP GET / Get One 
```
➜  ~ curl -XGET -H "Content-Type: application/json" http://127.0.0.1:9999/api/v1/Stu/1
{"ID":1,"CreatedAt":"2019-03-28T10:20:06+08:00","UpdatedAt":"2019-03-28T10:35:44+08:00","DeletedAt":null,"Name":"mike","Address":"shanghai.china.001","Sex":"male"}
```

- HTTP GET / Get All
```
➜  ~ curl -XGET -H "Content-Type: application/json" http://127.0.0.1:9999/api/v1/Stu
{"data":[{"ID":1,"CreatedAt":"2019-03-28T10:20:06+08:00","UpdatedAt":"2019-03-28T10:35:44+08:00","DeletedAt":null,"Name":"mike","Address":"shanghai.china.001","Sex":"male"}],"page":0,"pageCount":1,"pageSize":10}
```

- HTTP GET / Get With Where 
```
➜  ~ curl -XGET -H "Content-Type: application/json" http://127.0.0.1:9999/api/v1/Stu\?where\=ID\=\?%20and%20Name\=\?\&values\=1,mike\&pageSize\=20
{"data":[{"ID":1,"CreatedAt":"2019-03-28T10:20:06+08:00","UpdatedAt":"2019-03-28T10:35:44+08:00","DeletedAt":null,"Name":"mike","Address":"shanghai.china.001","Sex":"male"}],"page":0,"pageCount":1,"pageSize":20}
```

- HTTP DELETE / DELETE resourcese
```
➜  ~ curl -XDELETE -H "Content-Type: application/json" http://127.0.0.1:9999/api/v1/Stu/1
```

### RPC API Style Support 

- request StuService/SayHi Method 
```
➜  ~ curl -XPOST -H "Content-Type: application/json" http://127.0.0.1:9999/api/v1/StuService/SayHi -d'["mike"]'
["Hi, Mr.mike.\n",null]

➜  ~ curl -XPOST -H "Content-Type: application/json" http://127.0.0.1:9999/api/v1/StuService/SayHi -d'["alex"]'
["Hi, Mrs.alex.\n",null]

➜  ~ curl -XPOST -H "Content-Type: application/json" http://127.0.0.1:9999/api/v1/StuService/SayHi -d'["mike"]'
["no one named mike found.",null]
```

- RPC API Paybody DESC 
  - support array paybody and must Keep in the same order as the function parameters.
	- example. func(id uint, name string), must request with POST Protocol AND postBody like [10, "name"].


## License

hybrid is under Apache v2 License. See the [LICENSE](https://github.com/henrylee2cn/teleport/raw/v5/LICENSE) file for the full license text

