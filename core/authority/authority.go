package authority

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"project/core/ledger"
	"project/models/base"
	"project/models/model"
	"project/utils"
)

type UserAuth struct {
	User         string
	Class        string
	Organization string
	Role         string
	Security     int
}

//获取权限对象UserAuth
func GetUserAuth(stub shim.ChaincodeStubInterface, user string) (auth *UserAuth,err error) {

	auth = nil
	args := []string{ledger.MemberBook, user}

	//获取账本
	var key string
	var bytes []byte

	key, err = base.CreateKey(stub, args[0], args[1:], ledger.Key_type_current)
	if err != nil {
		return
	}
	bytes, err = stub.GetState(key)
	if err != nil || bytes == nil{
		return
	}
	var applicant model.MemberBook
	if err := json.Unmarshal(bytes, &applicant); err != nil {
		return
	}
	auth = new (UserAuth)
	auth.User = user
	auth.Class = applicant.Class_s
	auth.Organization = applicant.Organization
	auth.Security = applicant.SecurityLevel_i
	auth.Role = applicant.Role_s

	return
}

//角色鉴权
func (r *UserAuth) RoleAuthentication(stub shim.ChaincodeStubInterface, code string) (ret int, msg string) {
	ret = base.StatusOK
	msg = base.MsgSuccess

	args := []string{ledger.RoleBook, r.Role}
	var bytes []byte

	key, err := base.CreateKey(stub, args[0], args[1:], ledger.Key_type_current)
	if err != nil {
		ret = base.StatusBlockchainError
		msg = base.MsgBlockchainError
		return
	}
	bytes, err = stub.GetState(key)
	if err != nil || bytes == nil{
		ret = base.StatusBlockchainError
		msg = base.MsgBlockchainError
		return
	}


	var RoleM model.RoleBook
	if err := json.Unmarshal([]byte(bytes), &RoleM); err != nil {
		ret = base.StatusBlockchainError
		msg = base.MsgBlockchainError
		return
	}

	if !utils.InArry(code, RoleM.OprList_a) {
		ret = base.StatusPermissionError
		msg = base.MsgNoPermission
		return
	}

	return
}

//鉴权
func (r *UserAuth) checkAuthentication(stub shim.ChaincodeStubInterface, key string,field string) (ret int, msg string) {
	ret = base.StatusOK
	msg = base.MsgSuccess
	isOk := false

	bytes, err := stub.GetState(key)
	if err != nil{
		ret = base.StatusBlockchainError
		msg = base.MsgBlockchainError
		return
	}
	if bytes == nil {
		ret = base.StatusParamValueError
		msg = base.MsgParamValueError
		return
	}
	//反序列化go
	var book map[string]interface{}
	if err := json.Unmarshal([]byte(bytes), &book); err != nil {
		ret = base.StatusBlockchainError
		msg = base.MsgBlockchainError
		return
	}

	var list []string
	if id, Ok := book[field].([]string); Ok {
		list = id
	} else {
		isOk = true
	}

	if r.User != "" && utils.InArry(r.User, list)  {
		isOk = true
	}

	if r.Class != "" && utils.InArry(r.Class, list)  {
		isOk = true
	}
	if !isOk {
		ret = base.StatusPermissionError
		msg = base.MsgNoPermission
	}

	return
}

//安全基本鉴权
func (r *UserAuth) SecurityAuthentication(index int) (isOk bool) {

	return index >= r.Security
}

//目录可见鉴权
func (r *UserAuth) VisibleAuthentication(stub shim.ChaincodeStubInterface, key string) (ret int, msg string) {

	return r.checkAuthentication(stub,key,"VisibleList_a")
}

//账本访问鉴权
func (r *UserAuth) AccessAuthentication(stub shim.ChaincodeStubInterface, key string) (ret int, msg string) {
	return r.checkAuthentication(stub,key,"ReadList_a")
}

//账本编辑鉴权
func (r *UserAuth) WriteAuthentication(stub shim.ChaincodeStubInterface, key string) (ret int, msg string) {
	return r.checkAuthentication(stub,key,"WriteList_a")
}
