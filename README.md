# margin-cache

[![GoDoc](https://godoc.org/github.com/branthz/margin-cache?status.svg)](https://godoc.org/github.com/branthz/margin-cache)
## Summary
this is a hign-performance stand-alone redis-like cacher, it's based on C/S architecture.
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
```

----------------------------------------------------
## TODO:
1.Support the Restful api
2.Support functions like list/set/publish/subscribe ...

