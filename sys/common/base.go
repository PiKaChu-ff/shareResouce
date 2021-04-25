package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

//将(***),(***),(***)格式的字符串拆分为[***,***,***]的切片
func SplitCondition(s string) []string{
	sep:="),("
	s = strings.TrimLeft(s,"(")
	s = strings.TrimRight(s,"),")

	arr:=strings.Split(s,sep)
	return arr
}

//dist0 := strings.FieldsFunc(str, SplitByMoreStr)
func SplitByMoreStr(r rune) bool {
	//常用分隔符可以都写上
	return r == '>' || r == '=' || r == '<'
}

func SplitStr(s string,sep string) (stringArr []string){
	stringArr = strings.Split(s,sep)
	return stringArr
}

func StringContain(s1 string,s2 string)bool{
	return strings.Contains(s1, s2)
}

func SplitByOperate(s string,sep string) QueryCondition {

	condition := SplitStr(s,sep)

	queryCondition := QueryCondition{condition[0],sep,condition[1]}

	return queryCondition
}

//切片反序列化
func UnmarshalSlice(str string) ([]map[string]interface{},error) {
	var a []map[string]interface{}

	err := json.Unmarshal([]byte(str), &a)
	if err != nil{
		fmt.Println("反序列化失败", err)
		return nil,errors.New("反序列化失败")
	}
	return  a,nil
}
