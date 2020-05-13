package uc

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
)

type UCServer struct {
	rsa      *common.RSA
	database *mgo.Database
	router   *httprouter.Router
	html     []byte
}

type UpdateData struct {
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}

func (s *UCServer) Create(w http.ResponseWriter, r *http.Request, parms httprouter.Params) {}
func (s *UCServer) Delete(w http.ResponseWriter, r *http.Request, parms httprouter.Params) {}
func (s *UCServer) Update(w http.ResponseWriter, r *http.Request, parms httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	c := s.database.C("user")
	body, _ := ioutil.ReadAll(r.Body)
	user := &UpdateData{}
	json.Unmarshal(body, user)
	tp, _ := base64.StdEncoding.DecodeString(user.OldPassword)
	opwd, _ := s.rsa.RSADecrypt(tp)
	oldPassword := common.MD5V(string(opwd))
	tp, _ = base64.StdEncoding.DecodeString(user.NewPassword)
	npwd, _ := s.rsa.RSADecrypt(tp)
	newPassword := common.MD5V(string(npwd))
	err := c.Update(bson.M{
		"_id":      parms.ByName("id"),
		"password": oldPassword,
	}, bson.M{"$set": bson.M{"password": newPassword}})
	if err != nil {
		fmt.Println(err.Error())
	}
	ret, _ := json.Marshal(user)
	w.Write(ret)
}

func (s *UCServer) Read(w http.ResponseWriter, r *http.Request, parms httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	c := s.database.C("user")
	user := &UserData{}
	c.Find(bson.M{"_id": parms.ByName("id")}).One(user)
	ret, _ := json.Marshal(user)
	w.Write(ret)
}
func (s *UCServer) Path() string {
	return "/users"
}

func (s *UCServer) RSA(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("content-type", "text/plain")
	w.Write(s.rsa.Publickey)
}

func (s *UCServer) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("content-type", "text/html")
	w.Write(s.html)
}

func (s *UCServer) Router(router *httprouter.Router) {
	router.GET("/rsa", s.RSA)
	router.GET("/", s.Index)
}

func NewServer(database *mgo.Database, file string, rsa *common.RSA) *UCServer {
	s := &UCServer{}
	s.html, _ = ioutil.ReadFile(file)
	s.database = database
	s.rsa = rsa
	return s
}
