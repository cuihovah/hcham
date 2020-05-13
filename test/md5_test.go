package test

import (
	"crypto/md5"
	"encoding/hex"
	"testing"
)

func md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func TestMd5(t *testing.T) {
	ret := md5V("123456")
	t.Log(ret)
}
