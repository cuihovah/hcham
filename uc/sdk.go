package uc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/rpc/jsonrpc"
	"strings"
)

type SDK struct {
	url string
}

func (s *SDK) HttpCheckPassword(username string, password string) (bool, error) {
		destUrl := fmt.Sprintf("%s/users/check-password", s.url)
		data := &LoginData{}
		data.Username = username
		data.Password = password
		v, _ := json.Marshal(data)
		req, err := http.NewRequest("PUT", destUrl, strings.NewReader(string(v)))
		req.Header.Add("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if resp.StatusCode == 200 {
			return true, nil
		} else {
			return false, err
		}
}

func (s *SDK) HttpGet(username string) (*UserData, error) {
	destUrl := fmt.Sprintf("%s/users/%s", s.url, username)
	user := &UserData{}
	resp, err := http.Get(destUrl)
	v, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(v, user)
	user.Password = ""
	return user, err
}

func (s *SDK) RPCCheckPassword(username string, password string) (bool, error) {
	//client, err := rpc.DialHTTP("tcp", s.url)
	client, err := jsonrpc.Dial("tcp", s.url)
	if err != nil {
		return false, err
	}
	data := &LoginData{}
	data.Username = username
	data.Password = password
	var reply bool
	err = client.Call("UCServer.RPCCheckPassword", data, &reply)
	return reply, err
}

func (s *SDK) RPCGet(username string) (*UserData, error) {
	user := &UserData{}
	client, err := jsonrpc.Dial("tcp", s.url)
	if err != nil {
		return nil, err
	}
	err = client.Call("UCServer.RPCSearch", username, user)
	return user, err
}

func (s *SDK) SetServ(url string) {
	s.url = url
}
