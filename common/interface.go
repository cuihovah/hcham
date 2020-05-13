package common

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type CRUDServer interface {
	Create(http.ResponseWriter, *http.Request, httprouter.Params)
	Delete(http.ResponseWriter, *http.Request, httprouter.Params)
	Update(http.ResponseWriter, *http.Request, httprouter.Params)
	Read(http.ResponseWriter, *http.Request, httprouter.Params)
	Path() string
	Router(*httprouter.Router)
}
