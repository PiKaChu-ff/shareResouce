package router

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"project/api/project"
)

//路由
type Router struct {
	handlers map[string]Handler
}


//添加handle
func (r *Router) addHandle(funcName string, handle HandlerFunc,code string, mw ...HandlerMiddleware) {
	r.handlers[funcName] = Handler{handle, code,mw}
}

//注册路由
func (r *Router) RegisterRouter() {
	r.handlers = make(map[string]Handler)

	fmt.Println("register function ......")
	defer fmt.Println("register function end......")

	//路由注册
	r.addHandle(project.MethothGetService, project.GetServiceList,"", nil)

	r.addHandle(project.MethothCreatGate, project.CreatGate,"",nil)
	r.addHandle(project.MethothGetGate, project.GetGate,"",nil)

	r.addHandle(project.MethothCreatRole, project.CreatRole,"3", project.RoleAuthCheck)
	r.addHandle(project.MethothQueryRole, project.QueryRole,"4", project.RoleAuthCheck)
	r.addHandle(project.MethothModRole, project.ModRole, "5",project.RoleAuthCheck)
	r.addHandle(project.MethothAssignRole, project.ModRole, "6",project.RoleAuthCheck)

	r.addHandle(project.MethothUserRegister, project.UserRegister,"7", project.RoleAuthCheck)
	r.addHandle(project.MethothQueryRegister, project.QueryRegister,"8", project.RoleAuthCheck)
	r.addHandle(project.MethothModRegister, project.ModRegister,"9" ,project.RoleAuthCheck)
	//r.addHandle("regisCopyright", project.RegisCopyright, nil)
	//r.addHandle("queryCopyright", project.QueryCopyright, nil)
	r.addHandle(project.MethothSaveEvd, project.SaveEvd,"10", project.RoleAuthCheck)
	r.addHandle(project.MethothQueryEvdList, project.QueryEvdList,"11", project.RoleAuthCheck)
	r.addHandle(project.MethothQueryEvd, project.QueryEvd, "12",project.RoleAuthCheck)
	r.addHandle(project.MethothModEvd, project.ModEvdAttr, "13",project.RoleAuthCheck)

	//r.addHandle("applyDataExchange", project.ApplyDataEx, nil)
	//r.addHandle("confirmEx", project.ConfirmEx, nil)
	//r.addHandle("getTaskData", project.GetTaskData, nil)
	project.SetServiceList(r.GetOprList())
}

func (r *Router) GetOprList() (list map[string]string) {

	for k,handle := range r.handlers {
		code := handle.GetMethodCode()
		if code != "" {
			list[k] = code
		}

	}
	return
}

func (r *Router) Routing(stub shim.ChaincodeStubInterface, function string, args []string) peer.Response {
	if handle, ok := r.handlers[function]; ok {
		return handle.HandleInvoke(stub, args)
	}
	return shim.Error("Invalid function name")
}

