package news

import (
	"../advert"
	"../webserv"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/url"
	"time"
)

type FrontServ struct {
	router    *httprouter.Router
	database  *mgo.Database
	ssName    string
	ssoUrl    string
	domain    string
	hostname  string
	Pages     *webserv.Pages
	AdvertSDK *advert.SDK
}

type NewsItem struct {
	Id    string    `json:"id" bson:"_id"`
	Title string    `json:"title" bson:"title"`
	Image string    `json:"image" bson:"image"`
	Type  string    `json:"type" bson:"type"`
	Time  time.Time `json:"time" bson:"time"`
}

func (s *FrontServ) Run(port string) {
	s.router.GET("/", s.Index)
	s.router.GET("/news/:id", s.GetNews)
	startStr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(startStr, s.router)
}

func (s *FrontServ) SetSSO(sername string, sname string, ssourl string, domain string) {
	s.domain = domain
	s.ssoUrl = ssourl
	s.ssName = sname
	s.hostname = sername
}

func ParseAdvertSession(r *http.Request, name string) (map[string]int, error) {
	ret := make(map[string]int)
	cookie, err := r.Cookie(name)
	if err != nil {
		return ret, err
	}
	value, _ := url.QueryUnescape(cookie.Value)
	buf, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return ret, err
	}
	err = json.Unmarshal(buf, &ret)
	return ret, err
}

func (s *FrontServ) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	recommend, err := ParseAdvertSession(r, "advert-recommend-id")
	prop := s.GetProp()
	if err != nil {
		cookie := &http.Cookie{}
		cookie.Name = "advert-recommend-id"
		rt, _ := json.Marshal(recommend)
		value := base64.StdEncoding.EncodeToString(rt)
		cookie.Value = url.QueryEscape(value)
		cookie.Domain = prop["domain"]
		cookie.Expires = time.Now().Add(24 * time.Hour * 365)
		http.SetCookie(w, cookie)
	}
	news := make([]NewsItem, 100, 100)
	err = s.database.C("news").Find(bson.M{}).All(&news)
	if err != nil {
		s.Pages.RenderPage(w, "news_index", map[string]interface{}{"News": []NewsItem{}})
		return
	}
	imageNews := news[1:]
	adv, _ := s.AdvertSDK.GetTextRecommends(recommend)
	cars := make([]NewsItem, 100, 100)
	err = s.database.C("news").Find(bson.M{"type": "汽车"}).All(&cars)
	cartoon := make([]NewsItem, 100, 100)
	err = s.database.C("news").Find(bson.M{"type": "动漫"}).All(&cartoon)
	imax, _ := s.AdvertSDK.GetIMaxRecommends()
	//err = s.database.C("advert").Find(bson.M{"level": "imax"}).Limit(3).All(&imax)
	s.Pages.RenderPage(w, "news_index", map[string]interface{}{
		"News":       news,
		"ImageNews":  imageNews,
		"FirstImage": news[0],
		"Adverts":    adv,
		"Cars":       cars,
		"Cartoon":    cartoon,
		"IMax":       imax,
	})
}

func (s *FrontServ) GetNews(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	recommend, err := ParseAdvertSession(r, "advert-recommend-id")

	if err != nil {
		ck, _ := r.Cookie("advert-recommend-id")
		fmt.Println(ck.Value)
		fmt.Println(err.Error())
	}
	prop := s.GetProp()
	cookie := &http.Cookie{}
	news := &NewsData{}
	s.database.C("news").FindId(params.ByName("id")).One(&news)
	recommend[news.Type] = recommend[news.Type] + 1
	rt, _ := json.Marshal(recommend)
	value := base64.StdEncoding.EncodeToString(rt)
	cookie.Value = url.QueryEscape(value)
	cookie.Name = "advert-recommend-id"
	cookie.Domain = prop["domain"]
	cookie.Expires = time.Now().Add(24 * time.Hour * 365)
	cookie.Path = "/"
	http.SetCookie(w, cookie)
	s.Pages.RenderPage(w, "content_page", news)
}

func (s *FrontServ) GetProp() map[string]string {
	return map[string]string{
		"domain":   s.domain,
		"sso":      s.ssoUrl,
		"ss-name":  s.ssName,
		"hostname": s.hostname,
	}
}

func NewFrontServ(database *mgo.Database) *FrontServ {
	s := &FrontServ{}
	s.AdvertSDK = &advert.SDK{}
	s.database = database
	s.router = httprouter.New()
	s.Pages = webserv.NewPages()
	return s
}
