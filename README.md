# margin-cache
a hign-performance stand-alone cacher
此缓存系统为无状态缓存，当系统或进程重启后数据会丢失，因此使用者需要自己监控缓存系统的运行状态。
系统启动后会记录当前的启动时间（整形），client端可通过ping命令维持与server端的心跳，server端返回pong并携带启动时间。
client端如果需要的话可通过时间值判断系统是否重启过。

友情提醒：
使用本缓存的业务需要注意：如果多业务共享一个缓存时建议key值设置防冲突，采用如下方式，
eg:
projectName.keyName
projectCode.keyName
----------------------------------------------------
2015-12-15
support:
set key value                                             设置key/value
get key                                                   获取key值
del key                                                   删除key
exists key                                                查询key值存在
setex key value time                                      设置带超时key/value
hset key field value                                      设置hash表的域/值
hget key field                                            获取hash表中域值
hsetex  key field value time                              设置hash表的域/值（带超时）
hdel key field                                            删除hash表中域
hexists key field                                         判断hash表中域存在
hdestroy key                                              删除hash表
---------------------------------------------------
2015-12-18
support more:
DECR key 						  整形值减1
DECRBY key count					  整形值减count				
INCR   key 						  整形值增加1
INCRBY  key count					  整形值增加count
PING							  网络状态检测
---------------------------------------------------
2015-12-20
support more:
HGETALL key						  获取hash表中所有的域/值(可直接保存到map[field][]byte中)
HMSET key map						  一次设置多个field/value，以map的方式提供输入参数
HMGET key field1 field2 ...				  一次获取多个域值
--------------------------------------------------
2015-1-9
support more:
KEYS *      						  获取所有的key值，目前只支持获取所有key

----------------------------------------------------
TODO:
1.通信协议扩展：HTTP
2.扩展功能如链表/集合/发布/订阅等功能
3.缓存持久化

