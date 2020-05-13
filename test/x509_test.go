package test

import (
	"../sso"
	"encoding/base64"
	//"fmt"
	"fmt"
	"testing"
)

var bits int

func TestRSADecode(t *testing.T) {
	var data []byte
	var err error
	ssSevr := sso.NewSSOServ("./sso/static/html/update-password.html")
	data, err = ssSevr.RSAEncrypt([]byte("123"))
	if err != nil {
		fmt.Println("错误", err)
	}
	fmt.Println("加密：", base64.StdEncoding.EncodeToString(data))
	origData, err := ssSevr.RSADecrypt(data) //解密
	if err != nil {
		fmt.Println("错误", err)
	}
	fmt.Println("解密:", string(origData))
}
