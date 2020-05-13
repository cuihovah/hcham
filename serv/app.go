package serv

import (
	"../common"
	"../sso"
	"../webserv"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/url"
)

type ClientServ struct {
	ssoUrl   string
	domain   string
	Pages    *webserv.Pages
	ssName   string
	hostname string
}

func (c *ClientServ) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session, err := common.ParseSession(r, c.ssName)
	// 没登录，那就得判断token
	if err != nil || session.Id == "" {
		values, _ := url.ParseQuery(r.URL.RawQuery)
		token := values.Get("token")
		if token != "" {
			session, err := sso.DecodeToken(c, token)
			if err != nil {
				sso.Redirect(c, w)
			} else {
				sso.SetCookie(c, w, session)
				c.Pages.RenderPage(w, "index", session)
			}
		} else {
			sso.Redirect(c, w)
		}
	} else {
		// OK Done!
		c.Pages.RenderPage(w, "index", session)
	}
}
func (c *ClientServ) Run(port string) {
	router := httprouter.New()
	router.GET("/logout", sso.Logout(c))
	router.GET("/", c.Index)
	router.POST("/", sso.Token(c))
	router.OPTIONS("/", sso.Token(c))
	router.GET("/clients", sso.GetCluster(c))
	startStr := fmt.Sprintf(":%s", port)
	sso.Register(c)
	http.ListenAndServe(startStr, router)
}
func (c *ClientServ) SetSSO(sername string, sname string, ssourl string, domain string) {
	c.domain = domain
	c.ssoUrl = ssourl
	c.ssName = sname
	c.hostname = sername
}
func (c *ClientServ) GetProp() map[string]string {
	return map[string]string{
		"domain":   c.domain,
		"sso":      c.ssoUrl,
		"ss-name":  c.ssName,
		"hostname": c.hostname,
	}
}

func NewClient() *ClientServ {
	c := &ClientServ{}
	c.Pages = webserv.NewPages()
	return c
}
