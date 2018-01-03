package handle

import (
	"hash/adler32"
	"sync"

	"github.com/branthz/margin-cache/hashmap"
)

// for access request msg format--------------------------------------------------
const (
	ADD_STRING = iota + 1
	ADD_STRING_RES
	DEL_KEY
	DEL_KEY_RES
	QUERY_KEY
	QUERY_KEY_RES
	KEEPALIVE
	KEEPALIVE_RES
	HELLO
	HELLO_RES
	ADD_DB
	ADD_DB_RES
)

const (
	statusOk   = 0
	decodeFail = -1
	addErr     = -2
	queryFail  = -3
	KEYHASHMAX = 1024 * 8
)

type hmapST struct {
	mp map[string]*hashmap.Cache
	mu sync.RWMutex
}

var hcacher = newHmapst()

func hmapRead(k string) (v *hashmap.Cache, ok bool) {
	hcacher.mu.RLock()
	v, ok = hcacher.mp[k]
	hcacher.mu.RUnlock()
	return
}

func hmapWrite(k string, v *hashmap.Cache) {
	hcacher.mu.Lock()
	hcacher.mp[k] = v
	hcacher.mu.Unlock()
	return
}

func hmapdel(k string) {
	hcacher.mu.Lock()
	delete(hcacher.mp, k)
	hcacher.mu.Unlock()
	return
}

func newHmapst() *hmapST {
	st := new(hmapST)
	st.mp = make(map[string]*hashmap.Cache)
	return st
}

type FusionError string

func (err FusionError) Error() string { return "Fusion Error: " + string(err) }

var doesNotExist = FusionError("Key does not exist ")

func SessionHash(data []byte) int {
	var sum uint32
	sum = adler32.Checksum(data)
	return int(sum) & (KEYHASHMAX - 1)
}
