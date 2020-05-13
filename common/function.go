package common

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
)

func MD5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func FileLoad(filepath string) []byte {
	privatefile, err := os.Open(filepath)
	defer privatefile.Close()
	if err != nil {
		return nil
	}
	privateKey := make([]byte, 2048)
	num, err := privatefile.Read(privateKey)
	return privateKey[:num]
}

func ParseSession(r *http.Request, name string) (*SessionData, error) {
	session, err := r.Cookie(name)
	if err != nil {
		return nil, err
	}
	str, _ := url.QueryUnescape(session.Value)
	ret := &SessionData{}
	err = json.Unmarshal([]byte(str), ret)
	return ret, err
}
