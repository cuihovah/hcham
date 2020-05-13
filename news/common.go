package news

import (
	"../webserv"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"time"
)

type BackServ struct {
	router   *httprouter.Router
	database *mgo.Database
	ssName   string
	ssoUrl   string
	domain   string
	hostname string
	Pages    *webserv.Pages
}

type NewsData struct {
	Id       string    `json:"id" bson:"_id"`
	Title    string    `json:"title" bson:"title"`
	Contents string    `json:"contents" bson:"contents"`
	Image    string    `json:"image" bson:"image"`
	UserId   string    `json:"user_id" bson:"user_id"`
	Type     string    `json:"type" bson:"type"`
	Time     time.Time `json:"time" bson:"time"`
}
