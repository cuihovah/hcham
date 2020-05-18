package test

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"testing"
	"time"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

func Hehe(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	w.Write([]byte("hello"))
}

func TestRPC(t *testing.T) {
	router := httprouter.New()
	router.GET("/hello-world", Hehe)
	l, e := net.Listen("tcp", ":9999")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	arith := new(Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()
	fmt.Println(router)
	//go http.Serve(l, router)
	go http.Serve(l, nil)
	time.Sleep(1 * time.Second)
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println(client)
		log.Fatal("dialing:", err)
	}
	args := &Args{7,8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d ???", args.A, args.B, reply)
	resp, _ := http.Get("http://localhost:9999/hello-world")
	v, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(v))
}
