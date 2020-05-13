package common

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Server struct {
	router *httprouter.Router
}

func loadServer(servs ...CRUDServer) *httprouter.Router {
	router := httprouter.New()
	for _, s := range servs {
		pth1 := fmt.Sprintf("%s/:id", s.Path())
		router.GET(pth1, s.Read)
		pth2 := fmt.Sprintf("%s", s.Path())
		router.GET(pth2, s.Create)
		pth3 := fmt.Sprintf("%s/:id", s.Path())
		router.PUT(pth3, s.Update)
		pth4 := fmt.Sprintf("%s/:id", s.Path())
		router.DELETE(pth4, s.Delete)
		s.Router(router)
	}
	return router
}

func NewCRUDServer(servs ...CRUDServer) *Server {
	s := &Server{}
	s.router = loadServer(servs...)
	return s
}

func (s *Server) Run(port string) {
	startStr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(startStr, s.router)
}
