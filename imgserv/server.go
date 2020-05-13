package imgserv

import (
	"../common"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"mime"
	"net/http"
	"path"
	"strings"
)

type Server struct {
	router *httprouter.Router
	spath  string
}

func (s *Server) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("content-type", "text/html")
	filepath := path.Join(s.spath, "html", "index.html")
	buf, _ := ioutil.ReadFile(filepath)
	w.Header().Set("content-type", "text/html")
	w.Write(buf)
}

func (s *Server) POST(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "cookie,set-cookie")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	file, header, err := r.FormFile("image")
	if err != nil {
		fmt.Println(err.Error())
	}
	data, _ := ioutil.ReadAll(file)
	fmd5 := common.MD5V(string(data))
	ext := path.Ext(header.Filename)
	filename := path.Join(s.spath, "images_base", strings.Join([]string{fmd5, ext}, ""))
	ioutil.WriteFile(filename, data, 0644)
	w.Write([]byte(fmt.Sprintf("/images/%s", strings.Join([]string{fmd5, ext}, ""))))
}

//func (s *Server) OPTIONS(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//	w.Header().Set("Access-Control-Allow-Headers", "cookie,set-cookie")
//	w.Header().Set("Access-Control-Allow-Credentials", "true")
//	w.WriteHeader(200)
//}

func (s *Server) GetImage(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	filename := param.ByName("name")
	filepath := path.Join(s.spath, "images_base", filename)
	data, _ := ioutil.ReadFile(filepath)
	mimeType := mime.TypeByExtension(path.Ext(filepath))
	w.Header().Set("content-type", mimeType)
	w.Write([]byte(data))
}

func (s *Server) Run(port string) {
	s.router.GET("/", s.Index)
	s.router.POST("/images", s.POST)
	// s.router.OPTIONS("/images", s.OPTIONS)
	s.router.GET("/images/:name", s.GetImage)
	startStr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(startStr, s.router)
}

func NewServer(staticPath string) *Server {
	s := &Server{}
	s.spath = staticPath
	s.router = httprouter.New()
	return s
}
