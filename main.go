package main

import (
	"fmt"
	"project/core"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func main() {
	//启动链码
	if err := shim.Start(core.NewClient()); err != nil {
		fmt.Printf("Starting ChainCode Failed:%s", err.Error())
	}else {
		fmt.Println("Starting ChainCode Success")
	}
}
