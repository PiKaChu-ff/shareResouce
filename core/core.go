package core

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"project/core/authority"
	"project/core/ledger"
	"project/core/router"
	"project/models/base"
)

//ChainCode
type ChainCode struct {
	router   router.Router
	authUser authority.Authority
	books    ledger.Ledgers
}

func NewClient() (client *ChainCode) {
	client = new(ChainCode)
	client.router.RegisterRouter()
	client.books.Register()
	client.authUser.Register()
	return client
}

func (t *ChainCode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("ChainCode Init start......")
	defer fmt.Println("ChainCode Init end......")
	//初始化基础设置和环境参数，例如创建用户基础权限账本
	//...
	return base.ResSuccess()
}

func (t *ChainCode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	//获取参数
	function, args := stub.GetFunctionAndParameters()

	//参数异常
	if len(args) == 0 {
		base.ResParamError()
	}

	//寻找路由
	return t.router.Routing(stub, function, args)
}
