package engine

import (
	"errors"
)

var sep = ","
//初始化[][]byte
//condition [[xxx,yyy,zzz],[],[],[]]
func Construct(param []string)(*Params,error){

	var model = new(Params)

	if param[0] == ActionSet && len(param[3])==0{
		//更新操作缺少值参数
		return nil,errors.New("ValueGroup Is Necessary")
	}
	for i,v :=range param{
		//
		switch i {
		case 0:
			//处理Action
			model.Action = v
			break
		case 1:
			//处理条件
			if len(v)==0 {
				break
			}
			model.GetConditions(v)
			break
		case 2:
			//处理对象（KEY）
			model.GetSourceGroup(v)
			break
		default:
			//处理数据对象
			model.GetValueGroup(v)
		}
	}
	
	return model,nil

}

//创建条件
func (m *Params)GetConditions(whereString string){
	//分割条件
	//获取条件组xxx,xxx,xxx yyy,yyy,yyy
	conditionArr := SplitCondition(whereString)

	for _,v := range conditionArr{
		//获取一组条件内的每个条件xxx xxx xxx
		m.ConditionGroup = append(m.ConditionGroup,GetLotCondition(v))
	}
}

//创建对象
func (m *Params)GetSourceGroup(v string){

	sourceArr := SplitStr(v,",")

	for _,k := range sourceArr {

		source := GetKey(k)
		m.SourceGroup = append(m.SourceGroup,source)
	}
}

//创建数据
func (m *Params)GetValueGroup(v string){

	valueArr,_ := UnmarshalSlice(v)

	for _,v := range valueArr {
		m.ValueGroup = append(m.ValueGroup,GetValGroup(v))
	}
}