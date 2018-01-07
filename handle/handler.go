package handle

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/branthz/margin-cache/common"
	"github.com/branthz/margin-cache/common/log"
	"github.com/branthz/margin-cache/hashmap"
)

var CacheSet *hashmap.Dbs

var GstartTime = time.Now().Unix()

type tclient struct {
	conn    *net.TCPConn
	wbuffer *bytes.Buffer
	le      *list.Element
	rder    *bufio.Reader
}

func newClient(pconn *net.TCPConn) *tclient {
	return &tclient{
		conn:    pconn,
		wbuffer: bytes.NewBuffer(make([]byte, 1024)),
		rder:    bufio.NewReader(pconn),
		le:      nil,
	}
}

func (tc *tclient) Clear() {
	tc.conn.Close()
	tc.wbuffer = nil
	tc.rder = nil
	termList.Remove(tc.le)
}

func TListen(db *hashmap.Dbs) {
	tcpAddr := &net.TCPAddr{
		Port: common.CFV.Outport,
	}
	tcpConn, err := net.ListenTCP("tcp4", tcpAddr)
	if err != nil {
		log.Error("%v", err)
		os.Exit(-1)
	}
	defer tcpConn.Close()
	CacheSet = db
	for {
		conn, err := tcpConn.AcceptTCP()
		if err != nil {
			log.Error("Accept failed:%v", err)
			continue
		}
		go readTrequest(conn) //long tcp connection
	}
}

func readBulk(reader *bufio.Reader, head string) ([]byte, error) {
	var err error
	var data []byte

	if head == "" {
		head, err = reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
	}
	switch head[0] {
	case ':':
		data = []byte(strings.TrimSpace(head[1:]))

	case '$':
		size, err := strconv.Atoi(strings.TrimSpace(head[1:]))
		if err != nil {
			return nil, err
		}
		if size == -1 {
			return nil, doesNotExist
		}
		lr := io.LimitReader(reader, int64(size))
		data, err = ioutil.ReadAll(lr)
		if err == nil {
			// read end of line
			_, err = reader.ReadString('\n')
		}
	default:
		return nil, FusionError("Expecting Prefix '$' or ':'")
	}

	return data, err
}

func readResponse(tc *tclient) (res []byte, err error) {
	var line string
	err = nil
	var size, expi int

	//read until the first non-whitespace line
	for {
		line, err = tc.rder.ReadString('\n')
		if len(line) == 0 || err != nil {
			log.Info("%v", err)
			return
		}
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			break
		}
	}

	if line[0] == '*' {
		size, err = strconv.Atoi(strings.TrimSpace(line[1:]))
		if err != nil {
			err = fmt.Errorf("MultiBulk reply expected a number")
			return
		}
		if size <= 0 {
			err = fmt.Errorf("cmd size less than 0")
			return
		}
		log.Debug("request parameters:")
		req := make([][]byte, size)
		for i := 0; i < size; i++ {
			req[i], err = readBulk(tc.rder, "")
			if err == doesNotExist {
				continue
			}
			if err != nil {
				return
			}
			fmt.Printf("%s-----\n", string(req[i]))
			// dont read end of line as might not have been bulk
		}
		fmt.Printf("\n")
		tc.wbuffer.Reset()

		switch string(req[0]) {
		case "PING":
			//res = fmt.Sprintf(":%d\r\n", GstartTime)
			fmt.Fprintf(tc.wbuffer, ":%d\r\n", GstartTime)

		case "DECR":
			var v int64
			v, err = CacheSet.DecrementInt64(string(req[1]), 1)
			if err != nil {
				return
			}
			//res = fmt.Sprintf(":%d\r\n", v)
			fmt.Fprintf(tc.wbuffer, ":%d\r\n", v)
		case "DECRBY":
			var v int64
			v, err = strconv.ParseInt(string(req[2]), 10, 64)
			if err != nil {
				err = fmt.Errorf("decrby expected a integer")
				return
			}
			v, err = CacheSet.IncrementInt64(string(req[1]), v)
			if err != nil {
				return
			}
			//res = fmt.Sprintf(":%d\r\n", v)
			fmt.Fprintf(tc.wbuffer, ":%d\r\n", v)

		case "INCRBY":
			var v int64
			v, err = strconv.ParseInt(string(req[2]), 10, 64)
			if err != nil {
				err = fmt.Errorf("decrby expected a integer")
				return
			}
			v, err = CacheSet.IncrementInt64(string(req[1]), v)
			if err != nil {
				return
			}
			//res = fmt.Sprintf(":%d\r\n", v)
			fmt.Fprintf(tc.wbuffer, ":%d\r\n", v)

		case "INCR":
			var v int64

			v, err = CacheSet.IncrementInt64(string(req[1]), 1)
			if err != nil {
				return
			}
			//res = fmt.Sprintf(":%d\r\n", v)
			fmt.Fprintf(tc.wbuffer, ":%d\r\n", v)

		case "SET":
			CacheSet.Set(string(req[1]), string(req[2]), hashmap.NoExpiration)
			//res = fmt.Sprintf("+\r\nok")
			tc.wbuffer.WriteString(":1\r\n")

		case "GET":
			v, ok := CacheSet.Get(string(req[1]))
			if !ok {
				err = fmt.Errorf("not find the key:%s", string(req[1]))
				return
			}
			vs := v.(string)
			//res = fmt.Sprintf("$%d\r\n%s\r\n", len(vs), vs)
			fmt.Fprintf(tc.wbuffer, "$%d\r\n%s\r\n", len(vs), vs)

		case "DEL":
			CacheSet.Delete(string(req[1]))
			tc.wbuffer.WriteString(":1\r\n")

		case "SETEX":
			if size < 3 {
				err = fmt.Errorf("parameters error")
				return
			}
			expi, err = strconv.Atoi(string(req[2]))
			if err != nil {
				err = fmt.Errorf("setex expected a time expiration")
				return
			}
			CacheSet.Set(string(req[1]), string(req[3]), time.Duration(expi*1e9))
			//res = fmt.Sprintf("+\r\nok")
			tc.wbuffer.WriteString("+\r\nok")

		case "EXISTS":
			_, ok := CacheSet.Get(string(req[1]))
			if !ok {
				tc.wbuffer.WriteString(":0\r\n")
			} else {
				tc.wbuffer.WriteString(":1\r\n")
			}
		case "HEXISTS":
			ok:=CacheSet.Hexist(string(req[1]),string(req[2]))
			if !ok {
				tc.wbuffer.WriteString(":0\r\n")
			} else {
				tc.wbuffer.WriteString(":1\r\n")
			}

		case "HSET":
			CacheSet.Hset(string(req[1]),string(req[2]),string(req[3]),hashmap.NoExpiration)
			tc.wbuffer.WriteString(":1\r\n")

		case "HSETEX":
			if size < 4 {
				err = fmt.Errorf("parameters error")
				return
			}
			expi, err = strconv.Atoi(string(req[3]))
			if err != nil {
				err = fmt.Errorf("setex expected a time expiration")
				return
			}
			CacheSet.Hset(string(req[1]),string(req[2]),string(req[3]),time.Duration(expi*1e9))
			tc.wbuffer.WriteString(":0\r\n")

		case "HMSET":
			if size%2 != 0 || size == 2 {
				err = fmt.Errorf("parameters error")
				return
			}
			CacheSet.Hmset(string(req[1]),req[2:size])
			tc.wbuffer.WriteString(":1\r\n")

		case "HGET":
			var v interface{}
			v,err=CacheSet.Hget(string(req[1]),string(req[2]))
			if err!=nil {
				return
			}
				vs := v.([]byte)
				fmt.Fprintf(tc.wbuffer, "$%d\r\n%s\r\n", len(vs), vs)

		case "HMGET":
			if size<3{
				err = fmt.Errorf("parameters error")
				return
			}
			var data [][]byte
			data,err=CacheSet.Hmget(string(req[1]),req[2:size])
			if err!=nil {
				return
			}
			//log.Debug("hmget,size:%d", size)
			fmt.Fprintf(tc.wbuffer, "*%d\r\n", size-2)
			for i:=0;i<len(data);i++{
				fmt.Fprintf(tc.wbuffer, "$%d\r\n%s\r\n", len(data[i]), data[i])
			}

		case "HDEL":
			CacheSet.Hdel(string(req[1]),string(req[2]))
			//res = fmt.Sprintf(":1\r\n")
			tc.wbuffer.WriteString(":1\r\n")

		case "HDESTROY":
			CacheSet.Hdestroy(string(req[1]))
			tc.wbuffer.WriteString(":1\r\n")


		case "KEYS":
			if size < 2 {
				err = fmt.Errorf("parameters error")
				return
			}

			var count int = 0
			_, err = fmt.Fprintf(tc.wbuffer, "*%10d\r\n", count)
			if err != nil {
				return
			}
			count, err = CacheSet.Getallkey(tc.wbuffer)
			if err != nil {
				return
			}
			countstr := fmt.Sprintf("%10d", count)
			copy(tc.wbuffer.Bytes()[1:11], []byte(countstr))

		case "HGETALL":
			if size < 2 {
				err = fmt.Errorf("parameters error")
				return
			}
			err=CacheSet.Hgetall(string(req[1]),tc.wbuffer)
			if err!=nil{
				return
			}
			//fmt.Fprintf(tc.wbuffer, "$%d\r\n%s\r\n",count,,)

		default:
			log.Warn("request not support:%s", string(req[0]))
			err = fmt.Errorf("req not support")
			return
		}
		res = tc.wbuffer.Bytes()
		//log.Debug("res:%s,length:%d", string(res), len(res))
		return
	}
	err = fmt.Errorf("req not support")
	return

	//return readBulk(tc.rder, line)
}

func readTrequest(conn *net.TCPConn) {
	log.Info("get in access tcp connection")
	tc := newClient(conn)
	tc.le = termList.PushFront(tc)
	defer tc.Clear()
	var data []byte
	var err error

	for {
		data, err = readResponse(tc)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Warn("%v", err)
			data = []byte(fmt.Sprintf("-ERR%s\r\n", err.Error()))
		}

		//log.Debug("response:%s\n", string(data))

		_, err = tc.conn.Write(data)
		if err != nil {
			log.Error("tcp write error:%v", err)
		}
		tc.rder.Reset(tc.conn)
		tc.wbuffer.Reset()
	}

	return
}

func Read(conn *net.TCPConn, data []byte) error {
	var num, n int
	var total int = len(data)
	var err error

	err = conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	if err != nil {
		return err
	}
	for {
		n, err = conn.Read(data[num:total])
		if err != nil {
			return err
		}
		num += n
		if num < total {
			continue
		} else {
			return nil
		}
	}
}

func Write(conn *net.TCPConn, data []byte) error {
	var total int = len(data)
	var num int = 0
	var err error
	var n int
	err = conn.SetWriteDeadline(time.Now().Add(time.Second * 2))
	if err != nil {
		return err
	}

	for {
		n, err = conn.Write(data[num:total])
		if err != nil {
			return err
		}
		num += n
		if num < total {
			continue
		} else {
			return nil
		}
	}
}
