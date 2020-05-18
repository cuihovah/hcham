package advert

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/rpc"
	"strings"
)

type SDK struct {
	url string
}

func (s *SDK) HTTPGetTextRecommends(prop map[string]int) ([]AdvertData, error) {
	destUrl := fmt.Sprintf("%s/recommend-text", s.url)
	v, _ := json.Marshal(prop)
	resp, err := http.Post(destUrl, "application/json", strings.NewReader(string(v)))
	if err != nil {
		return nil, err
	}
	var ret []AdvertData
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &ret)
	return ret, err
}

func (s *SDK) HTTPGetIMaxRecommends() ([]AdvertData, error) {
	destUrl := fmt.Sprintf("%s/imax-image-text", s.url)
	resp, err := http.Get(destUrl)
	if err != nil {
		return nil, err
	}
	var ret []AdvertData
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &ret)
	return ret, err
}

func (s *SDK) RPCGetTextRecommends(prop map[string]int) ([]AdvertData, error) {
	advert := make([]AdvertData, 0)
	client, err := rpc.DialHTTP("tcp", s.url)
	if err != nil {
		return nil, err
	}
	err = client.Call("Server.RPCGetTextRecommends", prop, &advert)
	return advert, err
}

func (s *SDK) RPCGetIMaxRecommends(num int) ([]AdvertData, error) {
	advert := make([]AdvertData, 0)
	client, err := rpc.DialHTTP("tcp", s.url)
	if err != nil {
		return nil, err
	}
	err = client.Call("Server.RPCGetIMaxRecommends", num, &advert)
	return advert, err
}

func (s *SDK) SetServ(url string) {
	s.url = url
}
