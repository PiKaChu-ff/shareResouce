package base

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"project/utils/simplejson"
	"reflect"
)

//code (code 400 is not allowed to use)
const (
	StatusOK                = 200  //sys
	StatusParamFormatError        = 500  //参数格式不满足通讯协议
	StatusParamValueError         = 501  //参数值不满足协议
	StatusDataAlreadyExist      = 502       //创建账本时账本数据不已存在
	StatusDataNotExist      = 503        //查询账本时，账本数据不存在
	StatusPermissionError   = 504       //无权限
	StatusBlockchainError   = 505      //链码内部错误
	StatusCreatKeyFailed   = 506       //创建索引失败
	StatusTheNewKeyRepeatedError   = 507   //新建的索引重复
	StatusAutoritySetNotExist = 508         //权限集不存在
	StatusDataFormatError      = 509        //账本数据格式错误
	StatusDataPathNotSet      = 510        //数据地址未设置
	StatusUnknownError      = 520
)


//err
const(
	MsgSuccess  =  "success"
	MsgIncorrectParamNum = "Incorrect number of parameter."
	MsgMissParam = "Miss required parameter."
	MsgParamFormatError = "Invalid parameter format."
	MsgParamValueError = "Para data value Invalid."
	MsgInvalidFunction = "Invalid function."
	MsgDataAlreadyExist = "Data already exist."
	MsgDataNotExist = "Ledger data not exist."
	MsgDataFormatError = "Ledger data format invalid"
	MsgAutoritySetNotExist = "Autority Set Not Exist"
	MsgNoPermission = "No permission."
	MsgBlockchainError = "Blockchain internal error."
	MsgCreatKeyFailed = "Creat composite key Failed."
	MsgJsonFormatError = "Json string format error."
	MsgDataPathNotSet = "Data path not set."
	MsgUnknownError = "Unknown error."
)


//复合键规则
func CreateKey(stub shim.ChaincodeStubInterface, objectName string, keywordsList []string, keyType string) (string, error) {
	if keyType == Key_type_simple {
		key := objectName
		for _,element := range keywordsList{
			key += "." + element
		}
		return key, nil
	} else {
		key, err := stub.CreateCompositeKey(objectName, keywordsList)
		if err != nil {
			return "", err
		} else {
			return key, nil
		}
	}
}

//key是否存证
func IsKeyExist(stub shim.ChaincodeStubInterface, key string) (exist bool, err error) {
	if result, err := stub.GetState(key); err != nil {
		//获取异常
		return false, err
	} else if len(result) > 0 {
		//已存在
		return true, nil
	}
	return false, nil
}

//查询
func GetState(stub shim.ChaincodeStubInterface, key string, ret interface{}) (value interface{}, err error) {
	//查询
	var result []byte
	if result, err = stub.GetState(key); err != nil {
		return
	}

	//json
	value = reflect.ValueOf(ret)
	if err = json.Unmarshal(result,&value); err != nil {
		return
	}

	return value, nil
}

//插入
func PutState(stub shim.ChaincodeStubInterface, key string, value interface{}, bookName string) (err error) {
	//序列化
	var jsonByte []byte
	if jsonByte, err = json.Marshal(&value); err != nil {
		return err
	}

	//插入账本名
	var js *simplejson.Json
	if js ,err = simplejson.NewJson(jsonByte); err != nil {
		return
	}

	js.SetPath([]string{"bookName"}, bookName)

	if jsonByte ,err = js.Encode(); err != nil {
		return
	}

	//存入数据库
	if err := stub.PutState(key, jsonByte); err != nil {
		return err
	}

	return nil
}

//删除
func DeleteState(stub shim.ChaincodeStubInterface, key string) (err error) {
	//数据删除
	if err := stub.DelState(key); err != nil {
		return err
	}
	return nil
}
func ResSuccess() peer.Response {
	return peer.Response {
		Status:  StatusOK,
		Message: MsgSuccess,
	}
}

func ResSuccessPayload( v interface{} ) peer.Response {
	payload, err := json.Marshal(v)
	if err != nil {
		return ResBlockchainError()
	}
	return peer.Response {
		Status:  StatusOK,
		Message: MsgSuccess,
		Payload: payload,
	}
}

func ResError(code int32, msg string) peer.Response {
	return peer.Response{
		Status:  code,
		Message: msg,
	}
}

func ResParamError() peer.Response {
	return peer.Response{
		Status:  StatusParamFormatError ,
		Message: MsgParamFormatError,
	}
}

func ResBlockchainError() peer.Response {
	return peer.Response{
		Status:  StatusBlockchainError,
		Message: MsgBlockchainError,
	}
}

func ResPermissionError() peer.Response {
	return peer.Response{
		Status:  StatusPermissionError,
		Message: MsgNoPermission,
	}
}

const (
	Key_type_simple = "simple"
	Key_type_composite ="composite"
)
