package news

import (
	"../common"
	"../uc"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/url"
	//"html/template"
	//"../sso"
	"../webserv"
	"io/ioutil"
	"net/http"
	//"net/url"
	//"path"
	"time"
)

func (s *BackServ) PostNews(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	session, err := common.ParseSession(r, s.ssName)
	if err != nil {
		fmt.Println(err.Error())
	}
	news := &NewsData{}
	v, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	json.Unmarshal(v, news)
	news.UserId = session.Id
	news.Time = time.Now().Local().UTC()
	news.Id = bson.NewObjectId().Hex()
	s.database.C("news").Insert(news)
	w.WriteHeader(201)
}
func (s *BackServ) Run(port string) {
	s.router.GET("/", s.PostPage)
	s.router.POST("/news", s.PostNews)
	startStr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(startStr, s.router)
}
func (s *BackServ) SetSSO(sername string, sname string, ssourl string, domain string) {
	s.domain = domain
	s.ssoUrl = ssourl
	s.ssName = sname
	s.hostname = sername
}
func (s *BackServ) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session, err := common.ParseSession(r, s.ssName)
	// 没登录，那就得判断token
	if err != nil || session.Id == "" {
		values, _ := url.ParseQuery(r.URL.RawQuery)
		token := values.Get("token")
		if token != "" {
			session, err := uc.DecodeToken(s, token)
			if err != nil {
				uc.Redirect(s, w)
			} else {
				uc.SetCookie(s, w, session)
				s.Pages.RenderPage(w, "index", session)
			}
		} else {
			uc.Redirect(s, w)
		}
	} else {
		// OK Done!
		s.Pages.RenderPage(w, "index", session)
	}
}

func (s *BackServ) GetProp() map[string]string {
	return map[string]string{
		"domain":   s.domain,
		"sso":      s.ssoUrl,
		"ss-name":  s.ssName,
		"hostname": s.hostname,
	}
}
func (s *BackServ) PostPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session, err := common.ParseSession(r, s.ssName)
	// 没登录，那就得判断token
	if err != nil || session.Id == "" {
		values, _ := url.ParseQuery(r.URL.RawQuery)
		token := values.Get("token")
		if token != "" {
			session, err := uc.DecodeToken(s, token)
			if err != nil {
				uc.Redirect(s, w)
			} else {
				uc.SetCookie(s, w, session)
				s.Pages.RenderPage(w, "index", session)
			}
		} else {
			uc.Redirect(s, w)
		}
	} else {
		// OK Done!
		s.Pages.RenderPage(w, "index", session)
	}
}

func NewBackServ(database *mgo.Database) *BackServ {
	s := &BackServ{}
	s.database = database
	s.router = httprouter.New()
	s.Pages = webserv.NewPages()
	return s
}
