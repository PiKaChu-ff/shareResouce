package router

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

//HandlerFunc func
type HandlerFunc func(stub shim.ChaincodeStubInterface, strings []string) (p peer.Response)

//中间件
type HandlerMiddleware func(stub shim.ChaincodeStubInterface, strings []string,code string) (p *peer.Response)
type HandleMwChain []HandlerMiddleware

//handler
type Handler struct {
	Function HandlerFunc   //执行的方法
	MethodCode     string     //操作码
	Handlers HandleMwChain //中间件
}

func (h *Handler) HandleInvoke(stub shim.ChaincodeStubInterface, strings []string) (p peer.Response) {

	//中间件 若中间件已返回,则不执行func
	if p := h.handleMiddleware(stub, strings,h.MethodCode); p != nil {
		return *p
	}
	return h.Function(stub, strings)
}

func (h *Handler) handleMiddleware(stub shim.ChaincodeStubInterface, strings []string,code string) (p *peer.Response) {
	//顺序执行中间件
	for _, mw := range h.Handlers {
		if p := mw(stub, strings, code); p != nil {
			return p
		}
	}
	return nil
}

func (h *Handler) GetMethodCode() (code string) {

	code = h.MethodCode
	return
}
