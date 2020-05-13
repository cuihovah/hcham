package main

import (
	"./advert"
	"./common"
	"./imgserv"
	"./news"
	"./serv"
	"./sso"
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
	ct1 := serv.NewClient()
	ct2 := serv.NewClient()
	ct3 := serv.NewClient()
	ct4 := serv.NewClient()
	ct1.Pages.SetPage("index", "./serv/static/html/client-index-v1.html")
	ct2.Pages.SetPage("index", "./serv/static/html/client-index-v2.html")
	ct3.Pages.SetPage("index", "./serv/static/html/client-index-v3.html")
	ct4.Pages.SetPage("index", "./serv/static/html/client-index-v4.html")
	/********/
	ct1.SetSSO("http://www.ser-v1.com", "SSID1", "http://www.sso.com", ".ser-v1.com")
	ct2.SetSSO("http://www.ser-v2.com", "SSID2", "http://www.sso.com", ".ser-v2.com")
	ct3.SetSSO("http://www.ser-v3.com", "SSID3", "http://www.sso.com", ".ser-v3.com")
	ct4.SetSSO("http://www.ser-v4.com", "SSID4", "http://www.sso.com", ".ser-v4.com")
	/********/
	ssSevr := sso.NewSSOServ(session.DB("order"), "./sso/static/html/sso-login.html", rsa)
	go ssSevr.Run("8080")
	fmt.Println("starting...")
	time.Sleep(time.Second * 1)
	go ct1.Run("8081")
	go ct2.Run("8082")
	go ct3.Run("8083")
	go ct4.Run("8084")
	ucserv := uc.NewServer(session.DB("order"), "./uc/static/html/update-password.html", rsa)
	s := common.NewCRUDServer(ucserv)
	go s.Run("8091")
	simpleserv := advert.NewServer("./advert/static", session.DB("order"))
	simpleserv.SetSSO("http://advert.cuihovah-car.com/deliver", "advert-session-id", "http://www.sso.com", ".cuihovah-car.com")
	go simpleserv.Run("8092")
	image_serv := imgserv.NewServer("./imgserv/static")
	go image_serv.Run("8093")
	newServ := news.NewBackServ(session.DB("order"))
	newServ.Pages.SetPage("index", "./news/static/html/post.html")
	newServ.SetSSO("http://be.yeequeen.com", "news-session-id", "http://www.sso.com", ".yeequeen.com")
	go newServ.Run("8094")

	fServ := news.NewFrontServ(session.DB("order"))
	fServ.Pages.SetPage("news_index", "./news/static/html/index.html")
	fServ.Pages.SetPage("content_page", "./news/static/html/content-page.html")
	fServ.SetSSO("http://be.yeequeen.com", "news-session-id", "http://www.sso.com", ".yeequeen.com")
	fServ.AdvertSDK.SetServ("http://advert.cuihovah-car.com")
	go fServ.Run("8095")

	fmt.Println("started")
	<-wait
}
