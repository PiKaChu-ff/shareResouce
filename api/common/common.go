package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"project/core/authority"
	"project/core/ledger"
	"project/models/base"
	"project/utils"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type ComWorkSet struct {
	Stub     shim.ChaincodeStubInterface
	AuthUser *authority.UserAuth
	Books    *ledger.Ledgers
}


func Creat(set ComWorkSet, args []string, v interface{}) (ret int, msg string) {

	var data string
	ret = base.StatusOK
	msg = base.MsgSuccess
	position := ":api common-Creat"
	max := len(args) - 1

	if max < 1 || args[0] == "" || args[1] == "" {
		ret = base.StatusParamFormatError
		msg = base.MsgIncorrectParamNum + position
		return
	}

	A := args[0]
	regulars := make(map[string]reflect.Type)
	ret, msg = set.Books.GetFieldsTypeList(A, &regulars)
	if ret != base.StatusOK {
		return
	}

	ret, msg, data = utils.ContentsCheck(v, regulars)
	if ret != base.StatusOK {
		return
	}

	key, err := base.CreateKey(set.Stub, args[0], args[1:],ledger.Key_type_current)
	if err != nil {
		ret = base.StatusCreatKeyFailed
		msg = base.MsgCreatKeyFailed + position
		return
	}
	Avalbytes, err1 := set.Stub.GetState(key)
	if err1 != nil {
		ret = base.StatusBlockchainError
		msg = base.MsgBlockchainError + position
		return
	}
	if Avalbytes != nil {
		ret = base.StatusDataAlreadyExist
		msg = base.MsgDataAlreadyExist + position
		return
	}

	err = set.Stub.PutState(key, []byte(data))
	if err != nil {
		ret = base.StatusBlockchainError
		msg = base.MsgBlockchainError + position
		return
	}
	fmt.Println("creat putstate key = %s,value = %s\n", key, data)
	return
}

// Transaction makes payment of X units from A to B
func Set(set ComWorkSet, args []string, v map[string]interface{}) (ret int, msg string) {

	ret = base.StatusOK
	msg = base.MsgSuccess
	position := ":api common-Set"


	if len(args) < 2 || args[0] == "" || args[1] == "" {
		ret = base.StatusParamFormatError
		msg = base.MsgIncorrectParamNum
		return
	}

	A := args[0]
	B, _ := base.CreateKey(set.Stub, args[0], args[1:], ledger.Key_type_current)

	if set.AuthUser != nil {
		ret,msg = set.AuthUser.WriteAuthentication(set.Stub,B)
		if ret != base.StatusOK {
			return
		}
	}

	regulars := make(map[string]reflect.Type)
	ret, msg = set.Books.GetFieldsTypeList(A, &regulars)
	if ret != base.StatusOK {
		return
	}

	bytes, err := set.Stub.GetState(B)
	if err != nil {
		ret = base.StatusBlockchainError
		msg = base.MsgBlockchainError + position
		return
	}
	if bytes == nil {
		ret = base.StatusDataNotExist
		msg = base.MsgDataNotExist + position
		return
	}

	for k, t := range v {
		value := ""
		if _, ok := t.(string); !ok {
			tmp, _ := json.Marshal(t)
			value = string(tmp)
		} else if value = t.(string); value == "" {
			ret = base.StatusParamValueError
			msg = base.MsgParamValueError
			return
		}
		paras := []string{k, value}
		var iRet int
		bytes, iRet, _ = utils.KeyValueRefresh(paras, bytes, regulars)
		if iRet != base.StatusOK {

			return
		}
	}

	err = set.Stub.PutState(B, bytes)
	if err != nil {
		ret = base.StatusBlockchainError
		msg = base.MsgBlockchainError + position
		return
	}

	return
}

// query callback representing the query of a chaincode
func GetHistory(set ComWorkSet, args []string) (ret int, msg string) {
	ret = base.StatusOK
	msg = base.MsgSuccess
	position := ":api common-GetHistory"


	if len(args) < 2 || args[0] == "" || args[1] == "" {
		ret = base.StatusParamFormatError
		msg = base.MsgIncorrectParamNum + position
		return
	}

	A := args[0]
	regulars := make(map[string]reflect.Type)
	ret, msg = set.Books.GetFieldsTypeList(A, &regulars)
	if ret != base.StatusOK {
		return
	}

	B, _ := base.CreateKey(set.Stub, args[0], args[1:], ledger.Key_type_current)
	if set.AuthUser != nil {
		ret,msg = set.AuthUser.AccessAuthentication(set.Stub,B)
		if ret != base.StatusOK {
			return
		}
	}

	resultsIterator, err := set.Stub.GetHistoryForKey(B)
	if err != nil {
		ret = base.StatusBlockchainError
		msg = base.MsgBlockchainError + position
		return
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	bArrayMemberAlreadyWritten := false
	buffer.WriteString("[")

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			ret = base.StatusBlockchainError
			msg = base.MsgBlockchainError + position
			return
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	msg = buffer.String()

	return
}

// query callback representing the query of a chaincode
func GetField(set ComWorkSet, args []string) (ret int, msg string, data string) {

	ret = base.StatusOK
	msg = base.MsgSuccess
	position := ":api common-GetField"


	var Avalbytes string
	if len(args) < 3 || args[0] == "" || args[1] == "" || args[2] == "" {
		ret = base.StatusParamFormatError
		msg = base.MsgIncorrectParamNum + position
		return
	}

	A := args[0]
	regulars := make(map[string]reflect.Type)
	ret, msg = set.Books.GetFieldsTypeList(A, &regulars)
	if ret != base.StatusOK {
		return
	}

	B, _ := base.CreateKey(set.Stub, args[0], args[1:(len(args)-1)], ledger.Key_type_current)
	if set.AuthUser != nil {
		ret,msg = set.AuthUser.AccessAuthentication(set.Stub,B)
		if ret != base.StatusOK {
			return
		}
	}

	bytes, err := set.Stub.GetState(B)

	if err != nil {
		ret = base.StatusBlockchainError
		msg = base.MsgBlockchainError + position
		return
	}
	if bytes == nil {
		ret = base.StatusDataNotExist
		msg = base.MsgDataNotExist + position
		return
	}

	C := args[len(args)-1]
	if C == "" {
		ret = base.StatusParamValueError
		msg = base.MsgParamValueError + position
		return
	}

	var iRet int
	Avalbytes, iRet, msg = utils.KeyValueGet(C, bytes, regulars)

	if iRet != base.StatusOK {
		return
	}
	data = Avalbytes

	return
}

func GetOneZb(set ComWorkSet, args []string) (ret int, msg string, data string) {

	ret = base.StatusOK
	msg = base.MsgSuccess
	position := ":api common-GetOneZb"


	if len(args) < 2 || args[0] == "" || args[1] == "" {
		ret = base.StatusParamFormatError
		msg = base.MsgIncorrectParamNum + position
		return
	}

	//检验是否为注册过的账本
	A := args[0]
	regulars := make(map[string]reflect.Type)
	ret, msg = set.Books.GetFieldsTypeList(A, &regulars)
	if ret != base.StatusOK {
		return
	}

	B, _ := base.CreateKey(set.Stub, args[0], args[1:len(args)], ledger.Key_type_current)
	if set.AuthUser != nil {
		ret,msg = set.AuthUser.AccessAuthentication(set.Stub,B)
		if ret != base.StatusOK {
			return
		}
	}

	bytes, err := set.Stub.GetState(B)
	if err != nil {
		ret = base.StatusBlockchainError
		msg = base.MsgBlockchainError + position
		return
	}
	if bytes == nil {
		ret = base.StatusDataNotExist
		msg = base.MsgDataNotExist + position
		return
	}
	var r interface{}
	err = json.Unmarshal(bytes, &r)
	fmt.Println("\nresult : ", r)

	data = string(bytes)
	fmt.Println("getzb getstate key = %s,value = %s\n", B, data)
	return
}

func GetZbs(set ComWorkSet, args []string) (ret int, msg string, data []string) {
	ret = base.StatusOK
	msg = base.MsgSuccess
	data = []string{}
	position := ":api common-GetZbs"

	if len(args) < 1 {
		ret = base.StatusParamFormatError
		msg = base.MsgIncorrectParamNum + position

		return
	}

	A := args[0]
	isExist := set.Books.IsLedgerExist(A)
	if !isExist {
		ret = base.StatusParamFormatError
		msg = base.MsgParamFormatError + position
		return
	}

	paras := []string{}
	if len(args) > 1 {
		paras = args[1:]
	}
	resultIterator, err := set.Stub.GetStateByPartialCompositeKey(args[0], paras)
	if err != nil {
		ret = base.StatusBlockchainError
		msg = base.MsgBlockchainError + position

		return
	}
	defer resultIterator.Close()

	for resultIterator.HasNext() {
		item, err_ := resultIterator.Next()
		fmt.Println("\nvalue : ", item.Key, item.Value)
		if err_ != nil {
			ret = base.StatusBlockchainError
			msg = base.MsgBlockchainError + position
			return
		}

		if set.AuthUser != nil {
			ret,msg = set.AuthUser.VisibleAuthentication(set.Stub,item.Key)
			if ret != base.StatusOK {
				return
			}
		}

		data = append(data, string(item.Value))
	}

	fmt.Println("getzbs getstate value = %s\n", data)
	return
}

func Add(set ComWorkSet, args []string, v interface{}) (ret int, msg string) {

	ret = base.StatusOK
	msg = base.MsgSuccess
	position := ":api common-Add"

	if len(args) < 3 || args[0] == "" || args[1] == "" || args[2] == "" {
		ret = base.StatusParamFormatError
		msg = base.MsgIncorrectParamNum + position
		return
	}
	regulars := make(map[string]reflect.Type)
	isExist := set.Books.IsLedgerExist(args[0])
	if !isExist {
		ret = base.StatusParamFormatError
		msg = base.MsgParamFormatError + position
		return
	}

	key, _ := base.CreateKey(set.Stub, args[0], args[(len(args)-1):], ledger.Key_type_current)
	ret, msg = set.Books.GetFieldsTypeList(key, &regulars)
	if ret != base.StatusOK {
		return
	}

	ret, msg, _ = utils.ContentsCheck(v, regulars)
	if ret != base.StatusOK {
		return
	}

	B, _ := base.CreateKey(set.Stub, args[0], args[1:(len(args)-1)], ledger.Key_type_current)

	if set.AuthUser != nil {
		ret,msg = set.AuthUser.WriteAuthentication(set.Stub,B)
		if ret != base.StatusOK {
			return
		}
	}

	bytes, err := set.Stub.GetState(B)
	if bytes == nil {
		ret = base.StatusBlockchainError
		msg = base.MsgBlockchainError + position
		return
	}
	var r interface{}
	err = json.Unmarshal(bytes, &r)
	if err != nil {
		ret = base.StatusBlockchainError
		msg = base.MsgBlockchainError + position
		return
	}

	gobook, ok := r.(map[string]interface{})
	if ok {
		keyVal, success := gobook[args[2]].([]interface{})
		if !success {
			ret = base.StatusParamValueError
			msg = base.MsgParamValueError + position
			return
		}
		ResultPara, err_ := json.Marshal(v)
		if err_ != nil {
			ret = base.StatusBlockchainError
			msg = base.MsgBlockchainError + position
			return
		}
		keyvarNew := append(keyVal, ResultPara)
		gobook[args[2]] = keyvarNew

		newBytes, err_ := json.Marshal(gobook)
		if err_ != nil {
			ret = base.StatusBlockchainError
			msg = base.MsgBlockchainError + position
			return
		}
		err = set.Stub.PutState(B, newBytes)
		if err != nil {
			ret = base.StatusBlockchainError
			msg = base.MsgBlockchainError + position
			return
		}
		fmt.Println(string(newBytes))
	}

	return
}

func Query(set ComWorkSet, keyfilter string, v map[string]string) (ret int, msg string, data []string) {
	ret = base.StatusOK
	msg = base.MsgSuccess
	data = []string{}
	position := ":api common-Query"

	mapLen := len(v)
	if mapLen < 1 {
		ret = base.StatusParamFormatError
		msg = base.MsgParamFormatError + position
		return
	}

	var arry []interface{}
	var bytes []byte

	queryString := `{"selector":{`
	for k, t := range v {
		queryString += fmt.Sprintf(`"%s":"%s"`, k, t)
		mapLen--
		if mapLen > 0 {
			queryString += `,`
		}
	}
	queryString += `}}`

	fmt.Println("\n query key", queryString)
	resultsIterator, err := set.Stub.GetQueryResult(queryString)
	if err != nil {
		ret = base.StatusBlockchainError
		msg = base.MsgBlockchainError + position
		return
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		item, err_ := resultsIterator.Next()
		fmt.Println("\narry1 : ", arry)
		fmt.Println("\nvalue : ", item.Key, string(item.Value))
		if err_ != nil {
			ret = base.StatusBlockchainError
			msg = base.MsgBlockchainError + position
			return
		}
		if strings.Contains(item.Key, keyfilter) {
			if set.AuthUser != nil {
				ret,msg = set.AuthUser.VisibleAuthentication(set.Stub,item.Key)
				if ret != base.StatusOK {
					return
				}
			}
			var r interface{}
			err = json.Unmarshal(item.Value, &r)
			if err == nil {
				fmt.Println("\nr : ", r)
				arry = append(arry, r)
			}
		}
	}

	fmt.Println("\n arry:", arry)
	for _, t := range arry {
		bytes, err = json.Marshal(t)
		if err != nil {
			ret = base.StatusBlockchainError
			msg = base.MsgBlockchainError + position
			return
		}
		if bytes == nil {
			ret = base.StatusDataNotExist
			msg = base.MsgDataNotExist + position
			return
		}
		data = append(data, string(bytes))
	}

	return
}
