package handle

import (
	"hash/adler32"
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

type FusionError string

func (err FusionError) Error() string { return "Fusion Error: " + string(err) }

var doesNotExist = FusionError("Key does not exist ")

func SessionHash(data []byte) int {
	var sum uint32
	sum = adler32.Checksum(data)
	return int(sum) & (KEYHASHMAX - 1)
}
