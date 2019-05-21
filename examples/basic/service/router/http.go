package router

import (
	"log"
	"net/http"
	"time"

	"github.com/researchlab/hybrid/net/http/rest"

	"github.com/pressly/chi"
	"github.com/researchlab/hybrid/brick"
)

// HTTPService ...
type HTTPService struct {
	Controller *rest.Controller `inject:"RestController"`
	Config     brick.Config     `inject:"config"`
	server     *http.Server
}

// Init init http router
func (p *HTTPService) Init() error {
	log.Println("Init Router Http Service ...")
	defer func() {
		log.Println("Init Router Http Service Finished.")
	}()
	r := chi.NewRouter()

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

	port := p.Config.GetMapString("http", "port", "")
	if len(port) == 0 {
		panic("config invalid, server port invalid.")
	}
	log.Println("stu service start from port:", port)
	p.server = &http.Server{
		Addr:           ":" + port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Fatal(p.server.ListenAndServe())
	}()

	return nil
}

// Dispose http server close
func (p *HTTPService) Dispose() error {
	return p.server.Close()
}
