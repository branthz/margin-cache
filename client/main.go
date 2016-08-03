package main

import (
	"fmt"
	//"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

var client Client
var recycle int = 10000000

func init() {
	//client.Addr="127.0.0.1:6380"
	client.Addr = "10.10.2.77:6380"
	client.MaxPoolSize = 100
}

func main() {
	//runtime.GOMAXPROCS(4)
	//hsetBench(2, 1000000)
	//hgetall()
	//set()
	//function()
	//hset()
	setBench()
	//getBench()
	//allkeys()
	//for i:=0;i<40;i++{
	//	go hgetBench()
	//}
	//time.Sleep(1e9*20)
}

func function() {
	//set()
	//hset()
	//hget()
	//hdel()
	//hmset()
	hmget()
}

func allkeys() {
	//set()
	//fmt.Printf("=====\n")
	var tmstart = time.Now().Unix()
	val, err := client.Keys("*")
	if err != nil {
		fmt.Println(err)
		return
	}
	tmend := time.Now().Unix()
	fmt.Printf("time diff:%d(second)---%d\n", tmend-tmstart, len(val))
}

func hgetall() {
	var result = make(map[string]interface{})
	var tmstart = time.Now().Unix()
	fmt.Printf("================\n")
	err := client.Hgetall("2015", result)
	if err != nil {
		fmt.Println(err)
	}
	tmend := time.Now().Unix()

	fmt.Printf("time diff:%d(second)--%d\n", tmend-tmstart, len(result))
	return
}

func hgetBench() {
	var tmstart = time.Now().Unix()
	var field string = "zhang"
	var count int = 10000
	for i := 0; i < count; i++ {
		//key="hello"+strconv.Itoa(i)
		_, err := client.Hget("2015", field)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	tmend := time.Now().Unix()

	fmt.Printf("repeat counts:%d,time diff:%d(second)\n", count, tmend-tmstart)
}

func hsetBench(step, count int) {
	var tmstart = time.Now().Unix()
	var key string
	var value [128]byte
	for i := 0; i < len(value); i++ {
		value[i] = byte(i)
	}
	for i := 0; i < count; i++ {
		key = "hello" + strconv.Itoa(i+step*count)
		_, err := client.Hset("2015", key, value[:])
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	tmend := time.Now().Unix()

	fmt.Printf("time diff:%d(second)\n", tmend-tmstart)
}

func setBench() {
	//set()
	//hset()
	//hdel()
	//hget()
	wg := new(sync.WaitGroup)
	var tmstart = time.Now().Unix()

	var value [128]byte
	for i := 0; i < len(value); i++ {
		value[i] = byte(i)
	}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go setBenchChild(i, 1000, wg, value[:])
	}

	wg.Wait()

	tmend := time.Now().Unix()

	fmt.Printf("time diff:%d(second)\n", tmend-tmstart)
}

var downCount int32 = 0

func getBench() {
	wg := new(sync.WaitGroup)
	var tmstart = time.Now().Unix()

	for i := 0; i < 2000; i++ {
		wg.Add(1)
		go getBenchChild(i, 2000, wg)
	}

	wg.Wait()

	tmend := time.Now().Unix()

	fmt.Printf("time diff:%d(second)\n", tmend-tmstart)
}

func getBenchChild(i int, rec int, wg *sync.WaitGroup) {
	fmt.Printf("-----------child started:%d\n", i)
	var base = recycle / rec
	var key string
	for j := i * base; j < (i+1)*base; j++ {
		key = "hello" + strconv.Itoa(j)
		client.Get(key)
	}
	wg.Done()
	downCount = atomic.AddInt32(&downCount, 1)
	fmt.Printf("-------total:%d wg done\n", atomic.LoadInt32(&downCount))
}

func setBenchChild(i int, rec int, wg *sync.WaitGroup, value []byte) {
	fmt.Printf("-----------child start\n")
	var base = recycle / rec
	var key string
	for j := i * base; j < (i+1)*base; j++ {
		key = "hello" + strconv.Itoa(j)
		client.Set(key, value)
	}
	wg.Done()
	fmt.Printf("-------wg done\n")
}

type data struct {
	a int
	b string
	c int
}

type SliceHeader struct {
	addr uintptr
	len  int
	cap  int
}

func set() {
	var key = "zhang"
	var dd data
	dd.a = 10
	dd.c = 100
	dd.b = "nihao"
	dlen := unsafe.Sizeof(dd)
	println(int(dlen))
	ss := &SliceHeader{
		addr: uintptr(unsafe.Pointer(&dd)),
		len:  int(dlen),
		cap:  int(dlen),
	}
	ssb := *(*[]byte)(unsafe.Pointer(ss))
	err := client.Set(key, ssb)
	if err != nil {
		fmt.Println(err)
	}
	val, err := client.Get("zhang")
	if err != nil {
		fmt.Println(err)
	} else {
		xx := (*SliceHeader)(unsafe.Pointer(&val))
		var ff *data = (*data)(unsafe.Pointer(xx.addr))
		fmt.Printf("=======%v\n", *ff)
	}
}

func get() {
	val, err := client.Get("hello")
	if err != nil {
		fmt.Println(err)
	} else {
		println("hello", string(val))
	}
}

func exists() {
	b, err := client.Hexists("2015", "zhang")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("bool:%v\n", b)
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

func setget() {
	var err error
	var key = "hello"
	client.Set(key, []byte("world"))
	//val, _ := client.Get(key)
	//println(key, string(val))
	//time.Sleep(1e9 * 4)
	err = client.Setex("yuan", 30, []byte("shan"))
	if err != nil {
		fmt.Println(err)
		return
	}
	val, _ := client.Get("yuan")
	println("1:yuan--->", string(val))
	/*
		time.Sleep(1e9*40)
		val,err=client.Get("yuan")
		if err!=nil{
			fmt.Println(err)
			return
		}
		println("2:yuan--->",string(val))
	*/
}

func hget() {
	val, err := client.Hget("00000000000000000000b4430daed2ca>>cmd", "zhang")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("1->hget:%s\n", string(val))
}

func hset() {
	_, err := client.Hset("2015", "zhang", []byte("asdfghjklqwertyuiop"))
	if err != nil {
		fmt.Println(err)
		return
	}
}
func hmset() {
	src := make(map[string]string)
	src["wheather"] = "sunny"
	src["location"] = "pairs"
	err := client.Hmset("2015", src)
	if err != nil {
		fmt.Printf("hmset failed:%s\n", err.Error())
	} else {
		fmt.Printf("hmset ok\n")
	}
}

func hmget() {
	v, err := client.Hmget("00000000000000000000b4430daed2ca>>cmd", "cmd", "version")
	if err != nil {
		fmt.Printf("hmget,%v\n", err)
	} else {
		fmt.Printf("hmget ok:%v\n", v)
	}

}
