package sso

import (
	"../common"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ClientIfe interface {
	SetSSO(string, string, string, string)
	GetProp() map[string]string
}

func GetCluster(cfe ClientIfe) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	prop := cfe.GetProp()
	rtf := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		resp, _ := http.Get(fmt.Sprintf("%s/clients", prop["sso"]))
		ret, _ := ioutil.ReadAll(resp.Body)
		w.Header().Set("content-type", resp.Header.Get("content-type"))
		w.WriteHeader(200)
		w.Write(ret)
	}
	return rtf
}
func Logout(cfe ClientIfe) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	prop := cfe.GetProp()
	rtf := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "cookie,set-cookie")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.WriteHeader(200)
			return
		}
		cookie := &http.Cookie{}
		cookie.Name = prop["ss-name"]
		cookie.Value = ""
		cookie.Expires = time.Now()
		cookie.Domain = prop["domain"]
		http.SetCookie(w, cookie)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "cookie,set-cookie")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(200)
	}
	return rtf
}
func Token(cfe ClientIfe) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	prop := cfe.GetProp()
	rtf := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "cookie,set-cookie")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.WriteHeader(200)
			return
		}
		values, _ := url.ParseQuery(r.URL.RawQuery)
		token := values.Get("token")
		if token != "" {
			newUrl := fmt.Sprintf("%s/decode?token=%s", prop["sso"], token)
			resp, _ := http.Post(newUrl, "application/json", nil)
			//session := &SessionData{}
			v, _ := ioutil.ReadAll(resp.Body)
			//json.Unmarshal(v, session)
			cookie := &http.Cookie{}
			cookie.Name = prop["ss-name"]
			cookie.Value = url.QueryEscape(string(v))
			cookie.Domain = prop["domain"]
			cookie.Expires = time.Now().Add(86400 * time.Second)
			http.SetCookie(w, cookie)
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "cookie,set-cookie")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.WriteHeader(200)
		} else {
			w.WriteHeader(403)
		}
	}
	return rtf
}
func Register(cfe ClientIfe) {
	prop := cfe.GetProp()
	registerUrl := fmt.Sprintf("%s/register", prop["sso"])
	go http.Post(registerUrl, "text/plain", strings.NewReader(prop["hostname"]))
}
func DecodeToken(cfe ClientIfe, token string) (*common.SessionData, error) {
	prop := cfe.GetProp()
	newUrl := fmt.Sprintf("%s/decode?token=%s", prop["sso"], token)
	resp, err := http.Post(newUrl, "application/json", nil)
	if err != nil {
		return nil, err
	}
	session := &common.SessionData{}
	v, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(v, session)
	return session, err
}
func Redirect(cfe ClientIfe, w http.ResponseWriter) {
	prop := cfe.GetProp()
	newUrl := fmt.Sprintf("%s?redirect=%s&sso=1", prop["sso"], url.QueryEscape(prop["hostname"]))
	w.Header().Set("location", newUrl)
	w.WriteHeader(302)
}
func SetCookie(cfe ClientIfe, w http.ResponseWriter, session *common.SessionData) {
	prop := cfe.GetProp()
	cookie := &http.Cookie{}
	cookie.Name = prop["ss-name"]
	sstrify, _ := json.Marshal(session)
	cookie.Value = url.QueryEscape(string(sstrify))
	cookie.Domain = prop["domain"]
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(86400 * time.Second)
	http.SetCookie(w, cookie)
	w.WriteHeader(200)
}
