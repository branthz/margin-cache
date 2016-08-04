package zlconfig

import (
	"unicode"
	//"fmt"
	"errors"
	"strconv"
	"strings"
)

const (
	typeString = iota
	typeInt
	typeArrayString
	typeArrarInt
	typeBool
	typeFloat
	typeSubSection
)

var (
	errKeyOrType = errors.New("no this key or type wrong")
)

type section struct {
	cf Result
}

type valueST struct {
	flag int
	val  interface{}
}

func newSection() *section {
	var se section
	se.cf = make(Result)
	return &se
}

func (s *section) parseKey(line string) int {
	if line[0] == '[' {
		return SECTION_END
	}

	sv := strings.Split(line, "=")
	if len(sv) != 2 {
		return DEFAULTVAL
	}

	key := strings.TrimSpace(sv[0])
	v := strings.TrimSpace(sv[1])
	var vst valueST
	if len(v) == 0{
		return DEFAULTVAL
	}

	if v[0] == '[' && v[len(v)-1] == ']' {
		//array
		v = strings.Trim(v, "[]")
		arr := strings.Split(v, ",")
		if strings.Contains(v, "\"") {
			arrClean := getCleanArrayStr(arr)
			vst.val = arrClean
			vst.flag = typeArrayString
		} else {
			arrClean, err := getCleanArrayInt(arr)
			if err != nil {
				return valueNotFormat
			}
			vst.val = arrClean
			vst.flag = typeArrarInt
		}

	} else {
		if v[0] == '"' {
			//string
			v = strings.Trim(v, "\"")
			vst.val = v
			vst.flag = typeString
		} else if unicode.IsNumber(rune(v[0])) {
			if strings.ContainsAny(v, ".") {
				//float
				valFloat, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return valueNotFormat
				}
				vst.flag = typeFloat
				vst.val = valFloat
			} else {
				//int
				valInt, err := strconv.Atoi(v)
				if err != nil {
					return valueNotFormat
				}
				vst.flag = typeInt
				vst.val = valInt
			}
		} else {
			//bool
			if strings.EqualFold(v, "true") {
				vst.flag = typeBool
				vst.val = true
			} else if strings.EqualFold(v, "false") {
				vst.flag = typeBool
				vst.val = false
			} else {
				return valueNotFormat
			}
		}
	}

	s.cf[key] = vst
	return STATUSOK
}

func getCleanArrayInt(arr []string) ([]int, error) {
	var err error
	var arrInt = make([]int, len(arr))
	for i := 0; i < len(arr); i++ {
		arr[i] = strings.TrimSpace(arr[i])
		arrInt[i], err = strconv.Atoi(arr[i])
		if err != nil {
			return nil, err
		}
	}
	return arrInt, err
}

func getCleanArrayStr(arr []string) []string {
	var arrep = make([]string, len(arr))
	for i := 0; i < len(arr); i++ {
		arrep[i] = strings.Trim(arr[i], "\" ")
	}
	return arrep
}

func (re Result) GetArrayStr(key string) ([]string, error) {
	if v, ok := re[key]; ok && v.flag == typeArrayString {
		return v.val.([]string), nil
	}
	return nil, errKeyOrType
}

func (re Result) GetBool(key string) (bool, error) {
	if v, ok := re[key]; ok && v.flag == typeBool {
		return v.val.(bool), nil
	}
	return false, errKeyOrType
}

func (re Result) GetArrayInt(key string) ([]int, error) {
	if v, ok := re[key]; ok && v.flag == typeArrarInt {
		return v.val.([]int), nil
	}
	return nil, errKeyOrType
}

func (re Result) GetString(key string) (string, error) {
	if v, ok := re[key]; ok && v.flag == typeString {
		return v.val.(string), nil
	}
	return "", errKeyOrType
}

func (re Result) GetInt(key string) (int, error) {
	if v, ok := re[key]; ok && v.flag == typeInt {
		return v.val.(int), nil
	}
	return 0, errKeyOrType
}

func (re Result) GetIntDefault(key string, val int) int {
	if v, ok := re[key]; ok && v.flag == typeInt {
		return v.val.(int)
	}
	return val
}
func (re Result) GetBoolDefault(key string, val bool) bool {
	if v, ok := re[key]; ok && v.flag == typeBool {
		return v.val.(bool)
	}
	return val
}

func (re Result) GetStringDefault(key, val string) string {
	if v, ok := re[key]; ok && v.flag == typeString {
		return v.val.(string)
	}
	return val
}
