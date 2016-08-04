//package zlconfig
//copywrite @zhanglei
package zlconfig

import (
	//"errors"
	"io/ioutil"
	"os"
	"strings"
)

const (
	SECTION_END    = 1
	DEFAULTVAL     = 2
	STATUSOK       = 3
	KEYNOTINT      = 4
	valueNotFormat = 5
)

type Result map[string]valueST

func Readfile(path string, sectionName string) (Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return parseFile(string(content), sectionName)
}

func parseFile(data string, flag string) (Result, error) {
	s := newSection()
	var line string
	var ret int
	//var startLineNum,endLineNum int
	var sectionPick bool = false

	lines := strings.Split(data, "\n")
	for i := 0; i < len(lines); i++ {
		line = strings.Trim(lines[i], "\r\t\n ")
		if line == "" || line[0] == '#' {
			continue
		}

		if line[0] == '[' && line[len(line)-1] == ']' {
			if string(line[1:len(line)-1]) == flag {
				//startLineNum=i+1
				sectionPick = true
			}
			continue
		}

		if !sectionPick {
			continue
		}

		//parse key-value
		ret = s.parseKey(line)
		if ret == SECTION_END {
			return s.cf, nil
		}

	}

	return s.cf, nil
}
