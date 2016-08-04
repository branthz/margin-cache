package zlconfig

import (
	"fmt"
	"testing"
)

func TestToml(t *testing.T) {
	var filePath = "./example.toml"
	var section = "database"
	Param, err := Readfile(filePath, section)
	if err != nil {
		t.Error(err)
		return
	}
	//fmt.Println(Param)
	list, err := Param.GetArrayInt("ports")
	fmt.Println(list, err)
	str, err := Param.GetString("server")
	fmt.Println(str, err)
	interger, err := Param.GetInt("connection_max")
	fmt.Println(interger, err)
	bl, err := Param.GetBool("enabled")
	fmt.Println(bl, err)
	inte2, err := Param.GetInt("host")
	fmt.Println(inte2, err)
}
