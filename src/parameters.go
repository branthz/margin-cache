package main

import (
	"fmt"
	"package/zlconfig"
)

type defaultOp struct {
	outPort  int
	loglevel string
}

var Param zlconfig.Result
var CFV *defaultOp

func newConfig() *defaultOp {
	op := new(defaultOp)
	op.loglevel = "ERROR"
	return op
}

func uploadParameters() error {
	var err error
	CFV = newConfig()

	CFV.loglevel, err = Param.GetString("logLevel")
	if err != nil {
		return err
	}

	CFV.outPort, err = Param.GetInt("tcpListenPort")
	if err != nil {
		return err
	}

	//get this server inner's addr
	fmt.Printf("tcp listen port:%d\n", CFV.outPort)

	return nil
}
