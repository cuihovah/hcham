package main

import (
	"./advert"
	"./common"
	"./imgserv"
	"./news"
	"./uc"
	"fmt"
	"gopkg.in/mgo.v2"
	"time"
)

func main() {
	wait := make(chan bool)
	rsa := common.NewRSA()
	session, err := mgo.Dial("mongodb://cuihovah:cuihovah@pmo-hr.cuihovah.com:27017/order")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	session.SetMode(mgo.Monotonic, true)
	ucserv := uc.NewServer(session.DB("order"), "./uc/static/html/update-password.html", rsa)
	fmt.Println("starting...")
	time.Sleep(time.Second * 1)
	ucserv.Pages.SetPage("login", "./uc/views/sso-login.html")
	ucserv.Pages.SetPage("index", "./uc/views/index.html")
	ucserv.Pages.SetPage("register", "./uc/views/register.html")
	ucserv.SetSSO("http://be.cuihovah-user.com", "user-session-id", "http://be.cuihovah-user.com", ".cuihovah-user.com")
	ucserv.Start("9091")
	go ucserv.Run("8091")

	simpleserv := advert.NewServer("./advert/static", session.DB("order"))
	simpleserv.SetSSO("http://advert.cuihovah-car.com/deliver", "advert-session-id", "http://be.cuihovah-user.com", ".cuihovah-car.com")
	go simpleserv.Run("8092")

	image_serv := imgserv.NewServer("./imgserv/static")
	go image_serv.Run("8093")

	newServ := news.NewBackServ(session.DB("order"))
	newServ.Pages.SetPage("index", "./news/static/html/post.html")
	newServ.SetSSO("http://be.yeequeen.com", "news-session-id", "http://be.cuihovah-user.com", ".yeequeen.com")
	go newServ.Run("8094")

	fServ := news.NewFrontServ(session.DB("order"))
	fServ.Pages.SetPage("news_index", "./news/static/html/index.html")
	fServ.Pages.SetPage("content_page", "./news/static/html/content-page.html")
	fServ.SetSSO("http://be.yeequeen.com", "news-session-id", "http://be.cuihovah-user.com", ".yeequeen.com")
	fServ.AdvertSDK.SetServ("127.0.0.1:8092")
	go fServ.Run("8095")

	fmt.Println("started")
	<-wait
}
