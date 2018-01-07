package cmargin

import (
	"fmt"
	"runtime"
	"testing"
	//"time"
)

var (
	client    Client
	keystr    = "hz"
	hkeystr   = "2018-01-02"
	hkeyfield = "shanghai"
	hfieldval = "rainy"
)

func init() {
	runtime.GOMAXPROCS(4)
	//client.Addr = "192.168.206.110:6380"
	client.Addr = "127.0.0.1:6380"
	client.MaxPoolSize = 5
}

func Benchmark_Hget(b *testing.B) {
	client.Hset(hkeystr, hkeyfield, []byte(hfieldval))
	for i := 0; i < b.N; i++ {
		_, err := client.Hget(hkeystr, hkeyfield)
		if err != nil {
			b.Fatal(err.Error())
			return
		}
	}
}

func TestBasic(t *testing.T) {
	var val []byte
	var err error
	//var value = `hello\r\naaaaaaaaaaaaaaa
	//	bbbb`
	var value [128]byte
	for i := 0; i < 128; i++ {
		if i%2 == 0 {
			value[i] = 13
		} else {
			value[i] = 10
		}
	}

	err = client.Set("a", value[:128])
	if err != nil {
		t.Fatal("set failed", err.Error())
	}

	if val, err = client.Get("a"); err != nil {
		t.Fatal("get failed")
	}
	if string(val) != string(value[:]) {
		t.Fatal("set not equal get")
	}

	_, err = client.Keys("*")
	if err != nil {
		t.Fatal("keys * failed", err.Error())
	}
	client.Del("a")
	if ok, _ := client.Exists("a"); ok {
		t.Fatal(" delete  failed")
	}
	//v,_:=client.Decr("age")
	//v,_:=client.Incr("age")
	var age int64 = 11
	oringin, _ := client.Incr("age")
	t.Logf("incr:%d", oringin)
	v, _ := client.Incrby("age", age)
	if v != age+oringin {
		t.Fatalf("incrby failed,%d", v)
	}

	tm, err := client.Ping()
	if err != nil {
		t.Fatal("server out of connection", err.Error())
	}
	t.Log("status:%d\n", tm)
	//hset("chun","chun")
	//hset("xia","xia")
	hmset(t)
	hgetall(t)
	hdestroy(t)
}

func hdestroy(t *testing.T) {
	err := client.Hdestroy("2016")
	if err != nil {
		t.Fatalf("hdestory return:%v", err)
	}
}

func hmset(t *testing.T) {
	src := make(map[string]string)
	src["paris"] = "sunny"
	src["beijing"] = "cloudy"
	err := client.Hmset("2016", src)
	if err != nil {
		t.Fatalf("hmset failed:%s", err.Error())
	} else {
		t.Log("hmset ok")
	}
	v, err := client.Hmget("2016", "paris", "beijing")
	if err != nil {
		t.Fatalf("hmget failed:%s", err.Error())
	} else if len(v) != 2 || string(v[0]) != "sunny" || string(v[1]) != "cloudy" {
		t.Fatalf("hmget value not match hmset")
	}
	t.Log("hmset ok")
}

func hgetall(t *testing.T) {
	result := make(map[string][]byte)
	err := client.Hgetall("2016", result)
	if err != nil {
		t.Fatalf("hgetall failed:%s", err.Error())
	}
	t.Logf("%+v", result)
}

func exists(t *testing.T) {
	client.Hset("groupa","192.168.4.1",[]byte("alive"))
	b, err := client.Hexists("groupa", "192.168.4.1")
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("bool:%v", b)
	}
}

func hdel() {
	b, err := client.Hdel("2015", "zhang")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("bool:%v\n", b)
	}
}

func del() {
	b, err := client.Del("hello")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("bool:%v\n", b)
	}
}

func hget() {
	val, err := client.Hget("2015", "zhang")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("1->hget:%s\n", string(val))
}

func hset(t *testing.T) {
	_, err := client.Hset(hkeystr, hkeyfield, []byte("rainy"))
	if err != nil {
		t.Error(err)
	}
}
