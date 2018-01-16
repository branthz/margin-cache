package proxy 

import(
	"testing"
	"strings"
	"bufio"
)

func TestHttpHead(t *testing.T){
	const msg = "GET / HTTP/1.1\r\nHost: bar.com\r\n\r\n"
	m:=strings.NewReader(msg)	
	res:=httpHostHeader(bufio.NewReader(m))	
	t.Logf("get :%s",res)
}

func TestTcpHead(t *testing.T){
	const msg = "*3\r\n$3\r\nset\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
	m:=strings.NewReader(msg)	
	res:=getTcpKey(bufio.NewReader(m))
	t.Logf("get :%s",res)
}
