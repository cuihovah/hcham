package test

import (
	"fmt"
	"mime"
	"testing"
)

func TestMime(t *testing.T) {
	x := mime.TypeByExtension(".png")
	fmt.Println(x)
	fmt.Println("????????????")
	t.Log(x)
}
