package sso

import (
	"../common"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"sync"
)

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SessionData struct {
	Id   string `json:"id" bson:"_id"`
	Name string `json:"name"`
}

type SSOServ struct {
	rsa       *common.RSA
	loginHTML []byte
	//clients map[string]bool
	clients sync.Map
	db      *mgo.Database
}

func (s *SSOServ) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	buf := make([]byte, 1024*1024*10, 1024*1024*10)
	n, _ := r.Body.Read(buf)
	logindata := &LoginData{}
	err := json.Unmarshal(buf[:n], logindata)
	if err != nil {
		fmt.Println(err.Error())
	}
	pd, err := base64.StdEncoding.DecodeString(logindata.Password)
	passwd, _ := s.rsa.RSADecrypt(pd)
	spasswd := common.MD5V(string(passwd))
	if err != nil {
		fmt.Println(err.Error())
	}
	c := s.db.C("user")
	sd := &SessionData{}
	c.Find(bson.M{"_id": logindata.Username, "password": spasswd}).One(sd)
	d, _ := json.Marshal(sd)
	token := base64.StdEncoding.EncodeToString(d)
	w.Header().Set("content-type", "text/plain")
	w.Write([]byte(token))
}

func (s *SSOServ) SSOIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	header := w.Header()
	header.Set("content-type", "text/html")
	w.Write(s.loginHTML)
}

func (s *SSOServ) SSODecodeToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	values := r.URL.Query()
	token := values.Get("token")
	if token != "" {
		v, _ := base64.StdEncoding.DecodeString(token)
		header := w.Header()
		header.Set("content-type", "application/json")
		w.Write(v)
	} else {
		w.WriteHeader(400)
	}
}

func (s *SSOServ) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	buf, _ := ioutil.ReadAll(r.Body)
	cname := string(buf)
	//s.clients[cname] = true
	s.clients.Store(cname, true)
	w.WriteHeader(200)
}

func (s *SSOServ) GetClients(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	v := make([]string, 0)
	i := 0
	s.clients.Range(func(key interface{}, _ interface{}) bool {
		v = append(v, key.(string))
		i++
		return true
	})
	ret, _ := json.Marshal(v)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(200)
	w.Write(ret)
}

func (s *SSOServ) RSA(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("content-type", "text/plain")
	w.Write(s.rsa.Publickey)
}

func (s *SSOServ) Run(port string) {
	router := httprouter.New()
	router.GET("/", s.SSOIndex)
	router.POST("/login", s.Login)
	router.POST("/decode", s.SSODecodeToken)
	router.POST("/register", s.Register)
	router.GET("/clients", s.GetClients)
	router.GET("/rsa", s.RSA)
	startStr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(startStr, router)
}

func NewSSOServ(db *mgo.Database, file string, rsa *common.RSA) *SSOServ {
	s := &SSOServ{}
	s.db = db
	s.rsa = rsa
	s.loginHTML, _ = ioutil.ReadFile(file)
	//s.clients = make(map[string]bool)
	return s
}
