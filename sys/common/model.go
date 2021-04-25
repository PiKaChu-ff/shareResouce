package engine

const (
	ActionSet = "set"
	ActionGet = "get"
	ActionHistory = "history"
)

type Params struct {
	Action          string
	ConditionGroup  ConditionGroup
	SourceGroup     SourceGroup
	ValueGroup      []ValueGroup
}

//条件集合
type ConditionGroup []Conditions

type Conditions []ConditionOne
//一系列的字段组成账本KEY
type KeySeries []string

//单个KEY条件
type ConditionOne struct {
	Object          KeySeries
	QueryCondition  QueryCondition
}

//条件语句：字段+操作符+值，假如Key已存在的条件则Field为"*",Operate为"?",Value为"*",Key不存在的条件则Field为"*",Operate为"!?",Value为"*",
type QueryCondition struct{
	Field   string
	Operate string
	Value   interface{}
}

//对象集合
type SourceGroup []KeySeries

//对象	假如新增对象则Field为["*"]
/*type Source struct{
	KeyObject KeySeries
	Field     []string
}*/

//对象集合
type ValueGroup []ObjectValue

//值内容 self+6 则 Operate为"+=", Value为"6"
type ObjectValue struct {
	Field   string
	Operate string
	Value   string
}