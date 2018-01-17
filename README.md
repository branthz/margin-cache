# margin-cache

[![GoDoc](https://godoc.org/github.com/branthz/margin-cache?status.svg)](https://godoc.org/github.com/branthz/margin-cache)
[![Build Status](https://travis-ci.org/branthz/margin-cache.svg?branch=master)](https://travis-ci.org/branthz/margin-cache)

## Summary
margin-cache is a hign-performance stand-alone redis-like cacher, it's based on C/S architecture.
margin is a stateless , the data will lose when the progress restarted, so the client need to monitor the operational status of margin.
margin will record the init time so that the client-side can maintain the heartbeat through the ping command with it, margin return pong and carries the init time.
client can judge whether server has restarted by the value of time if needed.

## Support Command
Apis can refer to http://redis.io/commands, Now it supports the following cmd:
* set key value                                      
* get key                                                
* del key                                               
* exists key                                             
* setex key value time                                   
* hset key field value                                   
* hget key field                                            
* hsetex  key field value time                           
* hdel key field                                         
* hexists key field                                     
* hdestroy key                                            
* DECR key 						 
* DECRBY key count									
* INCR   key 						  
* INCRBY  key count					 
* PING							
* HGETALL key						  
* HMSET key map						  
* HMGET key field1 field2 ...				  
* KEYS *      						  

## Getting Started
### Installing
To start using margin-cache, install Go and run `go get`:

```sh
$ go get github.com/branthz/margin-cache
$ cd $GOPATH/src/github.com/branthz/margin-cache
$ make
$ ./marginCache -c marginCache.toml
```

## Example
this shows how to write your own code with cmargin package 
```go
package main

import (
	"github.com/branthz/margin-cache/cmargin"
	"fmt"
)
func main() {
	c:=&cmargin.Client{
		Addr:"127.0.0.1:6380",
		MaxPoolSize:5,
	}	
	value:=[]byte("value\r\n \t/* sjdf-+")
	c.Set("hello",value)
	v,err:=c.Get("hello")
	if err!=nil{
		fmt.Println(err)
        return
	}
	fmt.Printf("%q\n",string(v))
}
```
----------------------------------------------------
## TODO:
1.Support the Restful api  

