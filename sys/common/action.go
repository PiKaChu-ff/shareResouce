package engine

import (
	"errors"
	"regexp"
	"strings"
)

const(
	CountAdd = "Add"
	CountCut = "Cut"
)
//插入条件
/*func (m *Params) Where(query string, args ...string) (model *Params) {
	m.ConditionGroup.AppendCondition(query, args...)
	return m
}

//添加条件
func (c interface{})AppendCondition(condition ConditionOne) Conditions {

	c = append(c, condition)
	return c
}*/

//处理key
func GetKey(keys string) KeySeries {
	return  SplitStr(keys,".")
}



//特殊条件——KEY是否存在
func KeyIsExit(c string,exit bool) ConditionOne{
	var condition ConditionOne
	if exit {
		keySeries := c[1 : len(c)-1]

		condition.Object = GetKey(keySeries)
		condition.QueryCondition = QueryCondition{"*","?","*"}
	}else{
		keySeries := c[1 : len(c)-2]

		var condition ConditionOne
		condition.Object = GetKey(keySeries)
		condition.QueryCondition = QueryCondition{"*","!?","*"}
	}
	return condition
}

//创建单个条件，完整的条件由KEY+字段条件组成，例如key:A>10
func CreateCondition(key,field string) (condition ConditionOne){

	condition.Object = GetKey(key)

	condition.QueryCondition,_ = GetQueryCondition(field)
	return condition
}

//存在某个字符且用该字符分割
func ExitAndSplit(s,sep string){
	if strings.Index(s, sep)>0 {
		SplitByOperate(s,sep)
	}
}

//获取字段条件语句 例如：A > 10

func GetQueryCondition(s string) (queryCondition QueryCondition,err error) {
	var operateList = []string{">=","<=","==","!=",">","<"}

	for _,v:=range operateList{
		if strings.Index(s, v)>0 {
			queryCondition = SplitByOperate(s,v)
			return queryCondition,nil
		}
	}

	return QueryCondition{},errors.New("error operate")

}

//获取单组条件
///a.b.c:id=1,x.y.z:id=2,....
func GetLotCondition(v string) (conditions Conditions){
	conditionOneGroup := SplitStr(v,",")

	for _,k := range conditionOneGroup {
		//拆分条件中的key与字段条件
		conditionOne := SplitStr(k,":")

		var condition ConditionOne
		switch {
		case strings.Index(conditionOne[0], "!?")>0:
			//无字段条件，仅key条件
			//获取keys未分割状态
			condition = KeyIsExit(conditionOne[0],false)
			break
		case strings.Index(conditionOne[0], "?")>0:
			//获取keys未分割状态

			condition = KeyIsExit(conditionOne[0],true)
			break
		default:
			condition = CreateCondition(conditionOne[0],conditionOne[1])
		}
		conditions = append(conditions,condition)
	}
	return conditions
}

//将字段值中的自增自建情况处理
func GetValue(v string)(op,val string){

	//先处理特殊情况，自增自建
	//判断自增，以self+开头，或者self +(中间任意空格)
	reg := regexp.MustCompile(`^self\s*\+\s*`)
	if reg.MatchString(v){
		op = "+="
		val = reg.ReplaceAllString(v, "")
		return
	}

	reg = regexp.MustCompile(`^self\s*\-\s*`)
	//判断自减，以self-开头，或者self -(中间任意空格)
	if reg.MatchString(v){
		op = "-="
		val = reg.ReplaceAllString(v, "")
		return
	}

	//普通情况
	op = "="
	val = v
	return

}

//获取一个数据对象集合
func GetValGroup(v map[string]interface{})(objectValue []ObjectValue){

	for k,v:=range v{
		var ob ObjectValue
		ob.Field = k
		ob.Operate,ob.Value = GetValue(v.(string))
		objectValue = append(objectValue,ob)
	}

	return
}

