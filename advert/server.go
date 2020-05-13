package advert

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"
)

type Server struct {
	router   *httprouter.Router
	spath    string
	database *mgo.Database
	ssName   string
	ssoUrl   string
	domain   string
	hostname string
}

type UserData struct {
	Id      string  `json:"id" bson:"_id"`
	Name    string  `json:"name" bson:"name"`
	Balance float32 `json:"balance" bson:"balance"`
}

func (s *Server) recommend(prop map[string]int, num int) []AdvertData {
	kv := make(KeyValue, 0)
	for k, v := range prop {
		kv.Push(item{k, v})
	}
	list := make([]string, 0)
	for i := 0; i < 3; i++ {
		if kv.Len() > 0 {
			list = append(list, kv.Pop().(item).name)
		}
	}
	var adverts []AdvertData
	pipe := s.database.C("advert").Pipe([]bson.M{
		bson.M{"$match": bson.M{"type": bson.M{"$in": list}}},
		bson.M{"$sample": bson.M{"size": num}},
	})
	err := pipe.All(&adverts)
	if err != nil {
		fmt.Println(err.Error())
	}
	return adverts
}

func (s *Server) Advert(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	ssid, err := r.Cookie(s.ssName)
	header := w.Header()
	if err != nil || ssid.Value == "" {
		values, _ := url.ParseQuery(r.URL.RawQuery)
		token := values.Get("token")
		if token != "" {
			newUrl := fmt.Sprintf("%s/decode?token=%s", s.ssoUrl, token)
			resp, err := http.Post(newUrl, "application/json", nil)
			session := &UserData{}
			v, _ := ioutil.ReadAll(resp.Body)
			err = json.Unmarshal(v, session)
			if err != nil {
				str := url.QueryEscape(s.hostname)
				newUrl := fmt.Sprintf("%s?redirect=%s", s.ssoUrl, str)
				header.Set("location", newUrl)
				w.WriteHeader(302)
			} else {
				w.Header().Set("content-type", "text/html")
				filepath := path.Join(s.spath, "html", "deliver.html")
				buf, _ := ioutil.ReadFile(filepath)
				str := string(buf)
				t := template.New("name")
				t, err := t.Parse(str)
				if err != nil {
					fmt.Println(err.Error())
				}
				user := &UserData{}
				s.database.C("order_user").FindId(user.Id).One(user)
				cookie := &http.Cookie{}
				cookie.Name = s.ssName
				sstrify, _ := json.Marshal(session)
				cookie.Value = url.QueryEscape(string(sstrify))
				cookie.Domain = s.domain
				cookie.Path = "/"
				cookie.Expires = time.Now().Add(86400 * time.Second)
				http.SetCookie(w, cookie)
				w.WriteHeader(200)
				w.Header().Set("content-type", "text/html")
				t.Execute(w, user)
				//t.Execute(w, map[string]interface{}{"balance": user.Balance})
			}
		} else {
			str := url.QueryEscape(s.hostname)
			newUrl := fmt.Sprintf("%s?redirect=%s&sso=1", s.ssoUrl, str)
			header.Set("location", newUrl)
			w.WriteHeader(302)
		}
	} else {
		w.Header().Set("content-type", "text/html")
		filepath := path.Join(s.spath, "html", "deliver.html")
		buf, _ := ioutil.ReadFile(filepath)
		str := string(buf)
		t := template.New("name")
		t, err := t.Parse(str)
		if err != nil {
			fmt.Println(err.Error())
		}
		user := &UserData{}
		session, _ := r.Cookie(s.ssName)
		ss, _ := url.QueryUnescape(session.Value)
		json.Unmarshal([]byte(ss), user)
		err = s.database.C("order_user").FindId(user.Id).One(user)
		if err != nil {
			fmt.Println(err.Error())
		}
		w.Header().Set("content-type", "text/html")
		t.Execute(w, user)
		//t.Execute(w, map[string]interface{}{"balance": user.Balance})
	}
}

func (s *Server) Index(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	w.Header().Set("content-type", "text/html")
	filepath := path.Join(s.spath, "html", "index.html")
	buf, _ := ioutil.ReadFile(filepath)
	str := string(buf)
	t := template.New("index")
	t, err := t.Parse(str)
	if err != nil {
		fmt.Println(err.Error())
	}
	//var advert []AdvertData
	advert := &AdvertData{}
	pipe := s.database.C("advert").Pipe([]bson.M{
		bson.M{"$sample": bson.M{"size": 1}},
	})
	err = pipe.One(&advert)
	if err != nil {
		fmt.Println(err.Error())
	}
	t.Execute(w, advert)
}

func (s *Server) RecommendTextList(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	var adverts []AdvertData
	prop := make(map[string]int)
	buf, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(buf, &prop)
	if len(prop) == 0 {
		pipe := s.database.C("advert").Pipe([]bson.M{
			bson.M{"$match": bson.M{"contents": bson.M{"$ne": ""}}},
			bson.M{"$sample": bson.M{"size": 10}},
		})
		err = pipe.All(&adverts)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		adverts = s.recommend(prop, 5)
	}
	buf, err = json.Marshal(adverts)
	if err != nil {
		fmt.Println(err.Error())
	}
	w.Write(buf)
}

func (s *Server) RecommendIMax(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	var adverts []AdvertData
	pipe := s.database.C("advert").Pipe([]bson.M{
		bson.M{"$match": bson.M{"level": "imax"}},
		bson.M{"$sample": bson.M{"size": 10}},
	})
	err := pipe.All(&adverts)
	if err != nil {
		fmt.Println(err.Error())
	}
	buf, err := json.Marshal(adverts)
	if err != nil {
		fmt.Println(err.Error())
	}
	w.Write(buf)
}

func (s *Server) JS(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	w.Header().Set("content-type", "text/html")
	filepath := path.Join(s.spath, "js", param.ByName("path"))
	buf, _ := ioutil.ReadFile(filepath)
	w.Header().Set("content-type", "application/javascript;charset=utf-8")
	w.Write(buf)
}

func (s *Server) Pay(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	ssid, err := r.Cookie(s.ssName)
	session := &UserData{}
	ssid.Value, err = url.QueryUnescape(ssid.Value)
	json.Unmarshal([]byte(ssid.Value), session)
	err = s.database.C("order_user").UpdateId(session.Id, bson.M{"$inc": bson.M{"balance": 10}})
	if err != nil {
		fmt.Println(err.Error())
	}
	w.Header().Set("content-type", "text/plain")
	err = s.database.C("order_user").FindId(session.Id).One(session)
	retval := fmt.Sprintf("%.2f", session.Balance)
	w.Write([]byte(retval))
	w.WriteHeader(200)
}

func (s *Server) Deliver(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	user := &UserData{}
	advert := &AdvertData{}
	session, err := r.Cookie(s.ssName)
	if err != nil {
		fmt.Println(err.Error())
	}
	ss, _ := url.QueryUnescape(session.Value)
	buf, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal([]byte(ss), user)
	json.Unmarshal(buf, advert)
	advert.UserId = user.Id
	advert.Id = bson.NewObjectId().Hex()
	s.database.C("advert").Insert(advert)
	w.WriteHeader(201)
}

func (s *Server) Run(port string) {
	s.router.GET("/js/:path", s.JS)
	s.router.GET("/deliver", s.Advert)
	s.router.GET("/", s.Index)
	s.router.POST("/pay", s.Pay)
	s.router.POST("/adverts", s.Deliver)
	s.router.POST("/recommend-text", s.RecommendTextList)
	s.router.GET("/imax-image-text", s.RecommendIMax)
	startStr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(startStr, s.router)
}

func (s *Server) SetSSO(sername string, sname string, ssourl string, domain string) {
	s.domain = domain
	s.ssoUrl = ssourl
	s.ssName = sname
	s.hostname = sername
}

func NewServer(staticPath string, database *mgo.Database) *Server {
	s := &Server{}
	s.spath = staticPath
	s.database = database
	s.router = httprouter.New()
	return s
}
