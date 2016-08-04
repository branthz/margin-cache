package log

import(
	"testing"
)

func TestLogwrite(t *testing.T){
	loger,err:=New("/tmp/log","Debug")
	if err!=nil{
		loger.Error("stdout error (%v)",err)
		t.Error("log error")
	}
	loger.Info("good morning brant!")
	loger.Debug("hello 1")
	loger.Warn("hello 2")
	var a []byte=nil
	var b int
	loger.Debug("nihao--%s--%v",a,b)
}
/*
func TestXxx(t *testing.T){
	t.Error("no pass")
}
*/
