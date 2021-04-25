package ledger

import (
	"fmt"
	"project/models/base"
	"project/models/model"
	"reflect"
)

const (
	//账本
	GateBook      = "GateBook"
	ClassBook      = "ClassBook"
	MemberBook    = "MemberBook" //成员账本名称
	CopyRightBook = "CopyRightBook"
	EvdBook       = "EvdBook"
	EvdAccessBook = "EvdAccessBook"
	RoleBook       = "RoleBook"
	DataExBook    ="DataExBook"
)
const (
	Key_type_current = base.Key_type_composite
)

//账本读写规则
type Ledgers struct {
	regulars map[string]*Regular
}

var G_Books *Ledgers

func (r *Ledgers) Register() {

	fmt.Printf("ledger register function handle = %p......", r)
	defer fmt.Printf("ledger register function end......")

	r.regulars = make(map[string]*Regular)

	var v interface{}
	var list []string

	v = new(model.MemberBook)
	fmt.Printf("new v = %p......", v)
	list = []string{"UserId_s", "OwnerName_s", "OwnerCode_s", "Date_s"}
	r.registerLedger(v, "MemberBook", list)

	v = new(model.CopyRightBook)
	list = []string{"CertNo_s", "DataTable_s", "SummaryHash_s", "Type_s", "UserId_s", "DigitalMark_s", "GenerateTime_s", "ValidPeriod_s",
		"Institution_s", "Code_s", "Date_s"}
	r.registerLedger(v, "CopyRightBook", list)

	v = new(model.EvdBook)
	list = []string{"Handle_s", "DataTable_s", "UserId_s", "Contend_s", "CertNo_s", "DataType_s", "StorType_s", "Indentify_s",
		"Date_s", "Remark_s"}
	r.registerLedger(v, "EvdBook", list)

	G_Books = r
	fmt.Printf("G_Books = %p \n", G_Books)
}

//注册账本
func (r *Ledgers) registerLedger(v interface{}, name string, unModFields []string) (int, string) {

	regulars := new(Regular)

	fmt.Printf("\nregister function regulars = %p,name = %s", regulars, name)
	defer fmt.Printf("register function end......")

	ret, msg := regulars.saveRegular(v, unModFields)
	//fmt.Printf("save msg = %s,%p\n", msg, regulars)
	if ret == base.StatusOK {
		r.regulars[name] = regulars
		//fmt.Printf("r is = %x\n", r.regulars)
	}

	return ret, msg
}

func  ObscureKey(key string) (encKey string){
	encKey = key
	return
}

func VisualKey(key string) (decKey string){
	decKey = key
	return
}

func ProductKey(bookName string,req map[string]interface{}) (keyCore string,err error){

	return
}

//获取账本字段类型列表
func (r *Ledgers) GetFieldsTypeList(name string, typeList *map[string]reflect.Type) (ret int, msg string) {

	ret = base.StatusOK
	msg = base.MsgSuccess
	position := ":ledger-GetFieldsTypeList"

	_, ok := r.regulars[name]
	if !ok {
		ret = base.StatusDataFormatError
		msg = base.MsgDataNotExist + position
		return ret, msg
	}

	ret, msg = r.regulars[name].getFieldsTypeList(typeList)

	return ret, msg
}

//获取账本字段类型
func (r *Ledgers) GetFieldType(name string, field string, filedType *reflect.Type) (ret int, msg string) {

	ret = base.StatusOK
	msg = base.MsgSuccess
	position := ":ledger-GetFieldType"

	_, ok := r.regulars[name]
	if !ok {
		ret = base.StatusDataFormatError
		msg = base.MsgDataNotExist + position
		return ret, msg
	}

	ret, msg = r.regulars[name].getFieldType(field, filedType)

	return ret, msg
}

func (r *Ledgers) IsLedgerExist(name string) (isExist bool) {

	isExist = true
	_, ok := r.regulars[name]
	if !ok {
		isExist = false
	}

	return
}

//获取账本字段类型列表
func (r *Ledgers) GetFieldsAttrList(name string, attrList *map[string]int) (ret int, msg string) {

	ret = base.StatusOK
	msg = base.MsgSuccess
	position := ":ledger-GetFieldsAttrList"

	_, ok := r.regulars[name]
	if !ok {
		ret = base.StatusDataFormatError
		msg = base.MsgDataNotExist + position
		return ret, msg
	}

	ret, msg = r.regulars[name].getFieldsAttrList(attrList)

	return ret, msg
}

//获取账本字段类型
func (r *Ledgers) GetFieldAttr(name string, field string, filedAttr *int) (ret int, msg string) {

	ret = base.StatusOK
	msg = base.MsgSuccess
	position := ":ledger-GetFieldAttr"

	_, ok := r.regulars[name]
	if !ok {
		ret = base.StatusDataFormatError
		msg = base.MsgDataNotExist + position
		return ret, msg
	}

	ret, msg = r.regulars[name].getFieldAttr(field, filedAttr)

	return ret, msg
}
