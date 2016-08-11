package cmargin_test 

import (
	"fmt"
	"testing"
	"runtime"	
	//"time"
)

var client Client
/*
func main() {
	//set()
	hset()
	//hdel()
	hget()
}*/

func init(){
	runtime.GOMAXPROCS(4)
    	//client.Addr = "192.168.206.110:6380"
    	client.Addr = "10.10.2.103:6380"
	client.MaxPoolSize=5
}

func Benchmark_Hget(b *testing.B){
	hset("zhang","lei")
	for i:=0;i<b.N;i++{
		_, err := client.Hget("2015", "zhang")
        	if err != nil {
                	b.Fatal(err.Error())
                	return
        	}
	}
}

func TestBasic(t *testing.T){
	var val []byte
    	var err error
	//var value = `hello\r\naaaaaaaaaaaaaaa 
	//	bbbb`
	var value [128]byte
	for i:=0;i<128;i++{
		if i%2==0 { 
			value[i]=13
		}else{
			value[i]=10
		}
	}

    	err = client.Set("a", value[:128])
    	if err != nil {
        	t.Fatal("set failed", err.Error())
    	}

    	if val, err = client.Get("a"); err != nil {
        	t.Fatal("get failed")
    	}

	fmt.Printf("rawlength:%d, get value:%v, len:%d\n",len(value),val,len(val))
	fmt.Printf("--------------------\n")
	_,err=client.Keys("*")
	if err!=nil{
		t.Fatal("keys * failed",err.Error())
	}
	client.Del("a")
	if ok, _ := client.Exists("a"); ok {
        	t.Fatal("Should be deleted")
    	}
	//v,_:=client.Decr("age")
	//v,_:=client.Incr("age")

	v,_:=client.Incrby("age",11)
	fmt.Printf("%d",v)

	tm,err := client.Ping()
	if err != nil {
                t.Fatal("set failed", err.Error())
        }
	fmt.Printf("status:%d\n",tm)
	//hset("chun","chun")
	//hset("xia","xia")
	hmset()
	_,err=client.Hmget("2016","wheather","location")
	if err!=nil{
		fmt.Printf("hmget,%v\n",err)
	}else{
		fmt.Printf("hmset ok\n")
	}	
	hgetall()
	hdestroy()
}

func hdestroy(){
	err:=client.Hdestroy("2015")
	if err!=nil{
		fmt.Printf("hdestory return:%v",err)
	}
}

func hmset(){
	src :=make(map[string]string)
	src["wheather"]="sunny"
	src["location"]="pairs"
	err :=client.Hmset("2015",src)
	if err!=nil{
                fmt.Printf("hmset failed:%s\n", err.Error())
        }else{
		fmt.Printf("hmset ok\n")
	}
}

func hgetall(){
	result := make(map[string][]byte)
        err:=client.Hgetall("2015",result)
        if err!=nil{
                fmt.Printf("hgetall failed:%s\n", err.Error())
        }
        fmt.Printf("%+v\n",result)

}

func set() {
	var key = "hello"
	client.Set(key, []byte("world"))
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


func hget() {
	val, err := client.Hget("2015", "zhang")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("1->hget:%s\n", string(val))
}

func hset(key,val string) {
	_, err := client.Hset("2015",key, []byte(val))
	if err != nil {
		fmt.Println(err)
		return
	}
}
