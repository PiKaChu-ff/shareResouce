package utils

import (
	"encoding/json"
	"project/models/base"
	"reflect"
	"sort"
	"strconv"
	"fmt"
)

func InArry(target string, str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, target)
	if index < len(str_array) && str_array[index] == target {
		return true
	}
	return false
}
//通用方法
func KeyValueRefresh(args []string, bytes []byte, keyType_MAP map[string]reflect.Type) ([]byte, int, string) {

	position := ":utils-KeyValueRefresh"
	var newBytes []byte
	newBytes = nil

	ret, msg := CheckKeyValue(args, keyType_MAP)

	if ret != base.StatusOK {
		return newBytes, ret, msg
	}

	var r interface{}
	err := json.Unmarshal(bytes, &r)
	if err != nil {
		return newBytes, ret, msg
	}
	gobook, ok := r.(map[string]interface{})

	if ok {
		t := keyType_MAP[args[0]]
		switch t.Name() {
		case "int":
			i, err := strconv.Atoi(args[1])
			if err == nil {
				gobook[args[0]] = i
			} else {
				ret = base.StatusBlockchainError
				msg = base.MsgBlockchainError + position
			}
		case "string":
			gobook[args[0]] = args[1]

		case "float64":
			v, err := strconv.ParseFloat(args[1], 64)
			if err == nil {
				gobook[args[0]] = v
			} else {
				ret = base.StatusBlockchainError
				msg = base.MsgBlockchainError + position
			}
		}

		if ret == base.StatusOK {
			newBytes, err = json.Marshal(gobook)
			//fmt.Printf("\nkeyValueGet:", gobook, newBytes)
			if err != nil {
				ret = base.StatusBlockchainError
				msg = base.MsgBlockchainError + position
			}
		}
	}

	return newBytes, ret, msg
}

func KeyValueGet(arg string, bytes []byte, keyType_MAP map[string]reflect.Type) (string, int, string) {

	position := ":utils-KeyValueGet"
	var newBytes string
	args := []string{arg}

	newBytes = ""
	//fmt.Println("\nbefore check", args, keyType_MAP)
	ret, msg := CheckKeyValue(args, keyType_MAP)
	if ret != base.StatusOK {
		return newBytes, ret, msg
	}
	var r interface{}
	err := json.Unmarshal(bytes, &r)
	if err != nil {
		ret = base.StatusBlockchainError
		msg = base.MsgBlockchainError + position
		return newBytes, ret, msg
	}
	gobook, ok := r.(map[string]interface{})
	//fmt.Printf("get convert:", gobook)
	if ok {
		t := keyType_MAP[arg]
		switch t.Kind() {
		case reflect.Int:
			Vfloat, _ := gobook[arg].(float64)
			//fmt.Printf("Vfloat", Vfloat)
			newBytes = strconv.Itoa(int(Vfloat))

		case reflect.Float64:
			Vfloat, _ := gobook[arg].(float64)
			//fmt.Printf("Vfloat", Vfloat)
			newBytes = strconv.FormatFloat(Vfloat, 'E', -1, 64)

		case reflect.String:
			newBytes, ok = gobook[arg].(string)
			if !ok {
				ret = base.StatusBlockchainError
				msg = base.MsgBlockchainError + position
			}
		case reflect.Slice:
			var array []interface{}
			array, ok = gobook[arg].([]interface{})
			if ok {
				result, err_ := json.Marshal(array)
				if err_ == nil {
					newBytes = string(result)
				} else {
					ret = base.StatusBlockchainError
					msg = base.MsgBlockchainError + position
				}
			} else {
				ret = base.StatusBlockchainError
				msg = base.MsgBlockchainError + position
			}
		}
	}

	return newBytes, ret, msg
}

func ContentsCheck(vi interface{}, keyType_MAP map[string]reflect.Type) (int, string, string) {

	data := ""
	ret := base.StatusOK
	msg := base.MsgSuccess
	position := ":utils-ContentsCheck"

	strRet, err := json.Marshal(vi)

	// json转map
	var gobook map[string]interface{}
	err1 := json.Unmarshal(strRet, &gobook)
	if err1 != nil{
		ret = base.StatusParamValueError
		msg = base.MsgParamValueError + position
		return ret, msg, data
	}

	L:
		for k, t := range keyType_MAP {
			v, keyOk := gobook[k]
			if !keyOk {
				fmt.Printf("\ncannot find key", k)
				ret = base.StatusParamValueError
				msg = base.MsgParamValueError + position
				break L
			}
			okType := false

			switch t.Kind() {
			case reflect.Int:
				_, okType = v.(float64)
				if !okType {
					str, okFirst := v.(string)
					if !okFirst {
						okType = false
					} else if str == "" {
						okType = true
						gobook[k] = 0
					} else {
						value, err := strconv.Atoi(str)
						if err != nil {
							okType = false
						} else {
							okType = true
							gobook[k] = value
						}
					}
				}

			case reflect.Float64:
				_, okType = v.(float64)
				if !okType {
					str, okFirst := v.(string)
					if !okFirst {
						okType = false
					} else if str == "" {
						okType = true
						gobook[k] = 0.0
					} else {
						value, err := strconv.ParseFloat(str, 64)
						if err != nil {
							okType = false
						} else {
							okType = true
							gobook[k] = value
						}
					}
				}

			case reflect.String:

				_, okType = v.(string)


			case reflect.Slice:
				_, ok := v.([]interface{})
				if ok {
					okType = true
				}

			default:
				okType = true
			}
			if !okType {
				ret = base.StatusParamValueError
				msg = base.MsgParamValueError  + position

				break L
			}
		}


	//fmt.Printf("ContentsCheck ret", ret)
	if ret == base.StatusOK {
		resultStr, marErr := json.Marshal(gobook)
		if marErr == nil {
			data = string(resultStr)
		} else {
			ret = base.StatusBlockchainError
			msg = base.MsgBlockchainError + position
		}
	}
	return ret, msg, data

}

func CheckKeyValue(args []string, keys map[string]reflect.Type) (int, string) {
	ret := base.StatusOK
	msg := base.MsgSuccess
	key := args[0]
	hasValue := false
	position := ":utils-CheckKeyValue"

	if len(args) == 0 {
		ret = base.StatusParamFormatError
		msg = base.MsgIncorrectParamNum + position
		return ret, msg
	}

	if len(args) == 2 {
		hasValue = true
	}
	checkOk := false
L:
	for k, t := range keys {
		if k != key {
			continue
		}

		checkOk = true
		if hasValue {
			switch t.Name() {
			case "int":
				_, err := strconv.Atoi(args[1])
				if err != nil {
					checkOk = false
				}
			default:

			}
		}
		break L
	}

	if !checkOk {
		ret = base.StatusParamValueError
		msg = base.MsgParamValueError + position
	}

	return ret, msg
}

