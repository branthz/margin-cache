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
this not like exprie key in redis;it expire the field.
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

## Single Host Performance

Intel(R) Xeon(R) CPU E5620  @ 2.40GHz(8 cores)  
use the redis-benchmark tool
```
# redis-benchmark -h 127.0.0.1 -p 6380 -t set -r 1000000 -n 1000000 -d 2048
====== SET ======
  1000000 requests completed in 12.75 seconds
  50 parallel clients
  2048 bytes payload
  keep alive: 1
78449.84 requests per second

#redis-benchmark -h 127.0.0.1 -p 6380 -t  get  -r 10000 -n 5000000 -d 2048 -c 500
====== GET ======
  5000000 requests completed in 68.11 seconds
  500 parallel clients
  2048 bytes payload
  keep alive: 1
73406.35 requests per second
```

## Getting Started
### Installing
To start using margin-cache, install Go and run `go get`:

```sh
$ go get github.com/branthz/margin-cache
$ cd $GOPATH/src/github.com/branthz/margin-cache
$ make
$ ./marginCache -c marginCache.toml
```
This will bring up margin-cache listening on port 6380 for client communication 

### Proxy
you can also use tcp proxy which applied with consitent-hash in front of multiple margin-caches.
proxy will detect the cache's status if one is down, proxy will rehashing all the keys in backends and migrating them automatically. 
![image](https://raw.githubusercontent.com/branthz/resource/master/pic/proxy-margin.png)
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

