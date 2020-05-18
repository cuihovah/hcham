package uc

import (
	"../common"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/url"
)

func (s *UCServer) LoginPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session, err := common.ParseSession(r, s.ssName)
	// 没登录，那就得判断token
	if err != nil || session.Id == "" {
		values, _ := url.ParseQuery(r.URL.RawQuery)
		token := values.Get("token")
		if token != "" {
			session, err := DecodeToken(s, token)
			if err != nil {
				s.Pages.RenderPage(w, "login", nil)
			} else {
				SetCookie(s, w, session)
				user := &UserDataNoPassword{}
				s.RPCSearch(session.Id, user)
				s.Pages.RenderPage(w, "index", user)
			}
		} else {
			s.Pages.RenderPage(w, "login", nil)
		}
	} else {
		// OK Done!
		user := &UserDataNoPassword{}
		s.RPCSearch(session.Id, user)
		s.Pages.RenderPage(w, "index", user)
	}
}

func (s *UCServer) RegisterPage(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	s.Pages.RenderPage(w, "register", nil)
}