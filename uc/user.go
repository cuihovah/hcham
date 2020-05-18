package uc

import (
	"../common"
	"../webserv"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"net/url"
	"sync"
)

type UCServer struct {
	database *mgo.Database
	router   *httprouter.Router
	html     []byte
	Pages     *webserv.Pages
	ssName   string
	ssoUrl   string
	domain   string
	hostname string
	rsa       *common.RSA
	loginHTML []byte
	clients sync.Map
}

type UpdateData struct {
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}


func (s *UCServer) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logindata := &LoginData{}
	buf, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(buf, logindata)
	pd, err := base64.StdEncoding.DecodeString(logindata.Password)
	passwd, _ := s.rsa.RSADecrypt(pd)
	if err != nil {
		fmt.Println(err.Error())
	}
	logindata.Password = string(passwd)
	var ok bool
	err = s.RPCCheckPassword(*logindata, &ok)
	if ok == true {
		sd := &SessionData{}
		user := &UserDataNoPassword{}
		err := s.RPCSearch(logindata.Username, user)
		if err != nil {
			fmt.Println(err.Error())
		}
		sd.Id = user.Id
		sd.Name = user.Name
		d, _ := json.Marshal(sd)
		token := base64.StdEncoding.EncodeToString(d)
		w.Header().Set("content-type", "text/plain")
		w.Write([]byte(token))
	} else {
		w.WriteHeader(400)
	}

}
func (s *UCServer) SSOIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	header := w.Header()
	header.Set("content-type", "text/html")
	w.Write(s.loginHTML)

	session, err := common.ParseSession(r, s.ssName)
	// 没登录，那就得判断token
	if err != nil || session.Id == "" {
		values, _ := url.ParseQuery(r.URL.RawQuery)
		token := values.Get("token")
		if token != "" {
			session, err := DecodeToken(s, token)
			if err != nil {
				Redirect(s, w)
			} else {
				SetCookie(s, w, session)
				s.Pages.RenderPage(w, "index", session)
			}
		} else {
			Redirect(s, w)
		}
	} else {
		// OK Done!
		s.Pages.RenderPage(w, "index", session)
	}
}
func (s *UCServer) SSODecodeToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	values := r.URL.Query()
	token := values.Get("token")
	if token != "" {
		v, _ := base64.StdEncoding.DecodeString(token)
		header := w.Header()
		header.Set("content-type", "application/json")
		w.Write(v)
	} else {
		w.WriteHeader(400)
	}
}
func (s *UCServer) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	buf, _ := ioutil.ReadAll(r.Body)
	cname := string(buf)
	s.clients.Store(cname, true)
	w.WriteHeader(200)
}
func (s *UCServer) GetClients(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	v := make([]string, 0)
	i := 0
	s.clients.Range(func(key interface{}, _ interface{}) bool {
		v = append(v, key.(string))
		i++
		return true
	})
	ret, _ := json.Marshal(v)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(200)
	w.Write(ret)
}
func (s *UCServer) RSA(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("content-type", "text/plain")
	w.Write(s.rsa.Publickey)
}
func (s *UCServer) CheckPassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logindata := &LoginData{}
	v, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(v, logindata)
	value := common.MD5V(logindata.Password)
	userIndex := &UserIndex{}
	err := s.database.C("user_index").FindId(logindata.Username).One(userIndex)
	user := &UserData{}
	if err != nil {
		err = s.database.C("user").Find(bson.M{
			"_id": logindata.Username,
			"password": value,
		}).One(user)
	} else {
		err = s.database.C("user").Find(bson.M{
			"_id": user.Id,
			"password": value,
		}).One(user)
	}
	if err != nil {
		w.WriteHeader(400)
	} else {
		w.WriteHeader(200)
	}
}
func (s *UCServer) Create(w http.ResponseWriter, r *http.Request, parms httprouter.Params) {
	user := &UserData{}
	v, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	err = json.Unmarshal(v, user)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	var reply bool
	err = s.RPCRegister(*user, &reply)
	if err != nil || reply == false {
		w.WriteHeader(400)
		return
	}
	w.WriteHeader(201)
}
func (s *UCServer) Delete(w http.ResponseWriter, r *http.Request, parms httprouter.Params) {}
func (s *UCServer) Update(w http.ResponseWriter, r *http.Request, parms httprouter.Params) {
	user := &UserDataNoPassword{}
	v, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	err = json.Unmarshal(v, user)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	session, err := common.ParseSession(r, s.ssName)
	user.Id = session.Id
	var reply bool
	err = s.RPCUpdate(*user, &reply)
	if err != nil || reply == false {
		w.WriteHeader(400)
		return
	}
	w.WriteHeader(200)
}
func (s *UCServer) Search(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	key := param.ByName("term")
	userIndex := &UserIndex{}
	err := s.database.C("user_index").FindId(key).One(userIndex)
	if err != nil {
		user := &UserData{}
		err = s.database.C("user").FindId(key).One(user)
		if err != nil {
			w.WriteHeader(204)
		} else {
			ret, _ := json.Marshal(user)
			w.Write(ret)
		}

	} else {
		user := &UserData{}
		err = s.database.C("user").FindId(userIndex.UserId).One(user)
		if err != nil {
			w.WriteHeader(204)
		} else {
			ret, _ := json.Marshal(user)
			w.Write(ret)
		}
	}
}
func (s *UCServer) SetSSO(sername string, sname string, ssourl string, domain string) {
	s.domain = domain
	s.ssoUrl = ssourl
	s.ssName = sname
	s.hostname = sername
}
func (s *UCServer) RPC(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != "CONNECT" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "405 must CONNECT\n")
		return
	}
	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		log.Print("rpc hijacking ", r.RemoteAddr, ": ", err.Error())
		return
	}
	io.WriteString(conn, "HTTP/1.0 200 Connected to Go RPC\n\n")
	rpc.ServeConn(conn)
}
func (s *UCServer) GetProp() map[string]string {
	return map[string]string{
		"domain":   s.domain,
		"sso":      s.ssoUrl,
		"ss-name":  s.ssName,
		"hostname": s.hostname,
	}
}
func (s *UCServer) Run(port string) {
	startStr := fmt.Sprintf(":%s", port)
	s.router.GET("/rsa", s.RSA)
	//s.router.GET("/", s.SSOIndex)
	s.router.POST("/login", s.Login)
	s.router.DELETE("/logout", Logout(s))
	s.router.GET("/", s.LoginPage)
	s.router.POST("/decode", s.SSODecodeToken)
	s.router.POST("/register", s.Register)
	s.router.PUT("/user-update", s.Update)
	s.router.POST("/user-register", s.Create)
	s.router.GET("/register", s.RegisterPage)
	s.router.GET("/clients", s.GetClients)
	//s.router.GET("/", s.Index)
	// s.router.Handler("CONNECT", rpc.DefaultRPCPath, rpc.DefaultServer)

	l, _ := net.Listen("tcp", startStr)
	http.Serve(l, s.router)
}
func (s *UCServer) Start(port string) {
	rpc.Register(s)
	startStr := fmt.Sprintf(":%s", port)
	l, err := net.Listen("tcp", startStr)
	if err != nil {
		log.Printf("start sso rpc fatal: %s\n", err.Error())
	}
	go func(){
		for {
			conn, err := l.Accept()
			if err != nil {
				continue
			}
			go jsonrpc.ServeConn(conn)
		}
	}()
}

func NewServer(database *mgo.Database, file string, rsa *common.RSA) *UCServer {
	s := &UCServer{}
	s.html, _ = ioutil.ReadFile(file)
	s.database = database
	s.rsa = rsa
	s.router = httprouter.New()
	s.Pages = webserv.NewPages()
	return s
}
