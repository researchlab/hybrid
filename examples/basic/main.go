package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/researchlab/hybrid/examples/basic/lib/model"
	"github.com/researchlab/hybrid/examples/basic/service/router"
	"github.com/researchlab/hybrid/examples/basic/service/stu"
	"github.com/researchlab/hybrid/net/http/rest"

	"github.com/researchlab/hybrid/brick"

	"github.com/researchlab/hybrid/orm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

var (
	c = flag.String("c", "./config/config.json", "config path")
)

func main() {
	flag.Parse()
	configPath := *c
	stopSigs := make(chan os.Signal, 1)
	signal.Notify(stopSigs, syscall.SIGINT, syscall.SIGTERM)
	container := brick.NewContainer()
	container.Add(&brick.JSONConfigService{}, "config", brick.FactoryFunc(func() interface{} {
		return brick.NewJSONConfigService(configPath)
	}))
	container.Add(&model.Models{}, "Models", nil)
	container.Add(&mysql.MySQLService{}, "DB", nil)
	container.Add(&router.HTTPService{}, "HttpService", nil)
	container.Add(&rest.Controller{}, "RestController", nil)
	container.Add(&stu.StuService{}, "StuService", nil)
	container.Build()
	defer container.Dispose()
	select {
	case <-stopSigs:
		log.Println("service has been stoped.")
	}
}
