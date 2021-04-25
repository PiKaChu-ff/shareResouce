package ledger

import (
	//"fmt"
	"project/models/base"
	"reflect"
)

//handler
type Regular struct {
	typeList   map[string]reflect.Type
	fixedField map[string]int //0 非固定字段 1固定字段（不可重复写入）
}

func (h *Regular) saveRegular(v interface{}, unModefilds []string) (ret int, msg string) {

	h.typeList = make(map[string]reflect.Type)
	h.fixedField = make(map[string]int)

	ret = base.StatusOK
	msg = base.MsgSuccess
	position := ":ledger-saveRegular"

	subObject := reflect.ValueOf(v)
	subMyref := subObject.Elem()
	subTypeOfType := subMyref.Type()
	for j := 0; j < subMyref.NumField(); j++ {
		subField := subMyref.Field(j)
		subTypeField := subTypeOfType.Field(j)

		h.typeList[subTypeField.Name] = subField.Type()
		h.fixedField[subTypeField.Name] = 0
		//fmt.Printf("\nsub-%d : %s, %s\n", j, subTypeField.Name, h.typeList[subTypeField.Name])
	}

	for _, mw := range unModefilds {
		_, ok := h.typeList[mw]
		if !ok {
			ret = base.StatusParamFormatError
			msg = base.MsgParamFormatError + position
			return
		}
		h.fixedField[mw] = 1
	}

	return
}

//获取账本字段类型列表
func (h *Regular) getFieldsTypeList(typeList *map[string]reflect.Type) (ret int, msg string) {

	ret = base.StatusOK
	msg = base.MsgSuccess

	*typeList = h.typeList

	return
}

//获取账本字段类型
func (h *Regular) getFieldType(field string, filedType *reflect.Type) (ret int, msg string) {

	ret = base.StatusOK
	msg = base.MsgSuccess
	position := ":ledger-getFieldType"

	_, ok := h.typeList[field]
	if !ok {
		ret = base.StatusParamFormatError
		msg = base.MsgParamFormatError + position
		return
	}

	*filedType = h.typeList[field]

	return
}

//获取账本字段类型列表
func (h *Regular) getFieldsAttrList(attrList *map[string]int) (ret int, msg string) {

	ret = base.StatusOK
	msg = base.MsgSuccess

	*attrList = h.fixedField

	return
}

//获取账本字段类型
func (h *Regular) getFieldAttr(field string, filedAttr *int) (ret int, msg string) {

	ret = base.StatusOK
	msg = base.MsgSuccess
	position := ":ledger-getFieldAttr"

	_, ok := h.fixedField[field]
	if !ok {
		ret = base.StatusParamFormatError
		msg = base.MsgParamFormatError + position
		return
	}

	*filedAttr = h.fixedField[field]

	return
}
