package advert

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type SDK struct {
	advertUrl string
}

func (s *SDK) GetTextRecommends(prop map[string]int) ([]AdvertData, error) {
	destUrl := fmt.Sprintf("%s/recommend-text", s.advertUrl)
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

func (s *SDK) GetIMaxRecommends() ([]AdvertData, error) {
	destUrl := fmt.Sprintf("%s/imax-image-text", s.advertUrl)
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

func (s *SDK) SetServ(_url string) {
	s.advertUrl = _url
}
