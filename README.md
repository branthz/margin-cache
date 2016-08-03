# margin-cache
Summary
this is a hign-performance stand-alone cacher, it's based on C/S architecture.
margin is a stateless , the data will be lost when the process restarted, so the client need to monitor the operational status of margin.
margin will record the init time so that the client-side can maintain the heartbeat through the ping command with it, margin return pong and carries the init time.
client can judge whether server has restarted by the value of time if needed.

Support Command
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
DECR key 						  整形值减1
DECRBY key count					  整形值减count				
INCR   key 						  整形值增加1
INCRBY  key count					  整形值增加count
PING							  网络状态检测
HGETALL key						  获取hash表中所有的域/值(可直接保存到map[field][]byte中)
HMSET key map						  一次设置多个field/value，以map的方式提供输入参数
HMGET key field1 field2 ...				  一次获取多个域值
KEYS *      						  获取所有的key值，目前只支持获取所有key

----------------------------------------------------
TODO:
1.通信协议扩展：HTTP
2.扩展功能如链表/集合/发布/订阅等功能
3.缓存持久化

