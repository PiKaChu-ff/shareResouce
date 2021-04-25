package project

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/msp"
	"github.com/hyperledger/fabric/protos/peer"
	"project/api/common"
	"project/core/authority"
	"project/core/ledger"
	"project/models/base"
	"project/models/model"
	"project/utils"
	"time"
)

const (
	//账本
	GateBook      = "GateBook"
	MemberBook    = "MemberBook" //成员账本名称
	CopyRightBook = "CopyRightBook"
	EvdBook       = "EvdBook"
	EvdAccessBook = "EvdAccessBook"
	RoleBook       = "RoleBook"
	DataExBook    ="DataExBook"
)

const(
	MethothGetService  =  "getServiceList"
	MethothCreatGate = "creatGate"
	MethothGetGate = "getGate"
	MethothCreatRole = "creatRole"
	MethothQueryRole = "queryRole"
	MethothModRole = "modRole"
	MethothAssignRole = "assignRole"
	MethothUserRegister = "userRegister"
	MethothQueryRegister = "queryRegister"
	MethothModRegister = "modRegister"
	MethothSaveEvd = "saveEvd"
	MethothQueryEvdList = "queryEvdList"
	MethothQueryEvd = "queryEvd"
	MethothModEvd = "modEvd"
)

const (
	InstitutionName = ""
	InstitutionCode = ""
)

/*协议请求结构体*/
//查询返回
type serviceListResp struct {
	Msg              int //返回状态
	Data             map[string]string
}

//查询返回
type queryResp struct {
	Msg              int //返回状态
	Data             []string
}

//修改返回
type modResp struct {
	Msg              int //返回状态
	Date_s           string //登记日期（年/月/日/时）
}

type gateReq struct {
	GateName_s       string  //网关名称
	Admin_s          string  //管理员名称
	AdminCert_s      string  //管理员证书
	Organization_s   string  //所属单位
	Remark_s         string //简介
}

//企业成员注册
type userRegisterReq struct {
	ApplyId_s         string        //成员ID
	Class_s          string        //用户单位名称（单位或组织）
	Pubkey_s        string         //用户公钥
	UserCert_s      string         //用户证书
	Role_s           string        //单位社会唯一信用码
	SecurityLevel_i  int           //安全密级
	Remark_s         string        //单位介绍
}

//创建角色
type creatRoleReq struct {
	RoleName_s      string         //角色名称
	ApplyId_s       string         //申请成员ID
	OprList_a       []string      //权限操作列表
	Remark_s        string        //角色描述
}

//版权登记
type regisCopyrightReq struct {
	DataTable_s    string //数据标签
	ApplyId_s      string //申请成员ID
	UserId_s       string //版权方ID
	SummaryHash_s  string //数据摘要hash
	Type_s         string //Hash摘要算法
	Subject_s      string //学科
	AppFields_s    string //应用领域
	Suitable_s     string //适用对象
	GenerateTime_s string //数据生成日期（年/月/日/时）
	Remark_s       string //数据特点介绍
}

type regisCopyrightResp struct {
	Msg           string //返回状态
	CertNo_s      string //证书编号
	ValidPeriod_s string //有效期
	Institution_s string //区块链认证机构
	Code_s        string //认证机构社会唯一信用码
	Date_s        string //认证日期（年/月/日/时）
}

//数据存证
type saveEvdReq struct {
	DataTable_s       string   //数据标签
	Handle_s          string   //存证句柄
	DataType_s        string   //存证数据类型
	UserId_s          string   //数据拥有方ID
	Content_s         string   //数据内容
	SecurityLevel_i   int      //密级
	VisibleList_a     []string     //可见列表
	WriteList_a       []string //可操作用户ID列表
	ReadList_a        []string //可访问用户ID列表（普通存证有效）
	StorType_s        string   //存储类型。”normal”:默认普通存储；”privacy”：隐私存储；
	Indentify_s	      string   //隐私区域标识(隐私存证有效)
	DataFeature_s     string   //数据特征
	AnalysiRule_s     string   //数据定义
	BlockHash_s       string   //存证数据HASH
	Remark_s          string   //存证数据描述
}

type saveEvdResp struct {
	Msg         string //返回状态
	Handle_s    string //索引
}

type queryEvdListResp struct {
	DataTable_s       string   //数据标签
	Handle_s          string   //存证句柄
	DataType_s        string   //存证数据类型
	UserId_s          string   //数据拥有方ID
	Class_s	          string   //权限方部门
	Organization_s    string   //权限方组织
	StorType_s        string   //存储类型。”normal”:默认普通存储；”privacy”：隐私存储；
	Indentify_s	      string   //隐私区域标识(隐私存证有效)
	DataFeature_s     string   //数据特征
	AnalysiRule_s     string   //数据定义
	Date_s            string   //存证时间
	BlockHash_s       string   //存证数据HASH
	Remark_s          string   //存证数据描述
}

/*************
//数据交换请求
type applyDataExchangeReq struct {
	UserId_s      string //申请方ID
	Signature_s   string //验签字符串
	ResCode_s     string //资源编号
	ValidPeriod_s string //交易截至日期
	ValidTimes_s  string    //限制有效读取次数
	TargetId_s    string //交易目标方用户ID
	Status_s	  string //确认状态
	Remark_s      string //任务描述
}

type applyDataExchangeResp struct {
	Msg           int //返回状态
	TaskNo_s      string //任务编号 ）
	ValidPeriod_s string //交易截至日期
	ValidTimes_s  string    //限制有效读取次数
	ApplyDate_s   string //任务提交时间
	Status_s      string    //任务执行状态
}

type ExApply struct {
	UserId_s     string //申请方ID
	Signature_s  string //验签字符串
	TaskNo_s     string //任务编号
}

type ExResp struct {
	Msg          int //消息
	Path_s       string //路径
	SecretKey_s  string //加密后的数据
}
*//////////////

var ServiceMap map[string]string

func SetServiceList(list map[string]string){
	ServiceMap = make(map[string]string)
	ServiceMap = list
}

func  GetServiceList(stub shim.ChaincodeStubInterface, strings []string)  peer.Response {

	resq := serviceListResp{base.StatusOK, ServiceMap}

	return base.ResSuccessPayload(resq)
}

func  RoleAuthCheck(stub shim.ChaincodeStubInterface, strings []string,code string)  *peer.Response{
	//日志
	fmt.Println("开始角色检查1......")
	defer fmt.Println("角色检查结束1......")

	//反序列化go
	var req map[string]interface{}
	if err := json.Unmarshal([]byte(strings[0]), &req); err != nil {
		fmt.Println("反序列化转换失败")
		res := base.ResParamError()
		return &res
	}

	user, ok := req["ApplyId_s"].(string)
	if !ok {
		res := base.ResParamError()
		return &res
	}
	userAuth,err := authority.GetUserAuth(stub,user)
	if err != nil {
		res :=base.ResBlockchainError()
		return &res
	}
	if userAuth == nil {
		res := base.ResPermissionError()
		return &res
	}

	ret,msg :=userAuth.RoleAuthentication(stub,code)
	if ret != base.StatusOK {
		res := base.ResError(int32(ret),msg)
		return &res
	}
	return nil
}


func getEvd(content common.ComWorkSet, user string ,cert string,indenty string) (ret int, msg string,data string) {

	if indenty == "" {
		args := []string{EvdBook, cert}
		//获取账本
		ret, msg, data = common.GetOneZb(content, args)
	} else {
		ret = base.StatusPermissionError
		msg = base.MsgNoPermission
	}

	return
}

func getGate(content common.ComWorkSet ,key string) (ret int, msg string,data string) {

	args := []string{GateBook, key}

	ret, msg, data = common.GetOneZb(content, args)


	return
}

func getTaskApply(content common.ComWorkSet, user string,task string) (ret int, msg string,data string) {

	args := []string{DataExBook , task}

	ret, msg, data = common.GetOneZb(content, args)
	return
}

func getMspId(stub shim.ChaincodeStubInterface) string{
	creatorByte, err := stub.GetCreator()
	if err != nil {
		return ""
	}
	si := &msp.SerializedIdentity{}
	err = proto.Unmarshal(creatorByte, si)
	if err != nil {

		return ""
	}
	key := si.GetMspid()

	return key
}

func getCert(stub shim.ChaincodeStubInterface) string {
	creatorByte, err := stub.GetCreator()
	if err != nil {
		return ""
	}
	si := &msp.SerializedIdentity{}
	err = proto.Unmarshal(creatorByte, si)
	if err != nil {
		return ""
	}

	return string(si.GetIdBytes())
}


//创建网关
func CreatGate(stub shim.ChaincodeStubInterface, strings []string) peer.Response {
	//日志
	fmt.Println("开始创建网关......")
	defer fmt.Println("创建网关结束......")

	//反序列化go
	var req gateReq
	if err := json.Unmarshal([]byte(strings[0]), &req); err != nil {
		fmt.Println("反序列化转换失败\n")
		return base.ResParamError()
	}

	if req.GateName_s == ""  || req.Admin_s == "" || req.AdminCert_s == ""{
		fmt.Println("输入参数错误\n")
		return base.ResParamError()
	}

	key := getMspId(stub)
	cert := getCert(stub)

	if key == "" || cert == "" {
		return base.ResBlockchainError()
	}
	args := []string{GateBook, cert}

	book := model.GateBook{key,cert,req.GateName_s,req.Admin_s,req.AdminCert_s,
		    req.Organization_s,req.Remark_s,1}
	content := common.ComWorkSet{stub, nil, ledger.G_Books}
	//数据完整性校验
	ret, msg := common.Creat(content, args, book)
	if ret != base.StatusOK {
		fmt.Println("创建失败 %s\n", msg)
		return base.ResError(int32(ret), msg)
	}

	timeStr := time.Now().Format("2006-01-02 15:04:05")

	//创建角色
	RoleName := "RoleManager"
	RoleOpr := []string{ServiceMap[MethothCreatRole],ServiceMap[MethothQueryRole],ServiceMap[MethothModRole],ServiceMap[MethothAssignRole]}
	bookRole := model.RoleBook{RoleName,RoleOpr,req.Remark_s,1}

	args = []string{RoleBook, RoleName}
	//数据完整性校验
	ret, msg = common.Creat(content, args, bookRole)
	if ret != base.StatusOK {
		fmt.Println("创建角色失败 %s\n", msg)
		return base.ResError(int32(ret), msg)
	}

	//创建网关管理员账户
	bookAdmin := model.MemberBook{req.Admin_s,  0, req.Organization_s,"",
		cert,"",req.AdminCert_s,RoleName,"", timeStr, 1}

	args = []string{MemberBook, req.Admin_s}

	//数据完整性校验
	ret, msg = common.Creat(content, args, bookAdmin)
	if ret != base.StatusOK {
		fmt.Println("创建失败 %s\n", msg)
		return base.ResError(int32(ret), msg)
	}

	return base.ResSuccessPayload(ret)
}

//获取网关
func GetGate(stub shim.ChaincodeStubInterface, strings []string) peer.Response {
	//日志
	fmt.Println("开始获取网关......")
	defer fmt.Println("获取网关结束......")

	var resq queryResp
	content := common.ComWorkSet{stub, nil, ledger.G_Books}
	if strings[0] != "current" {
		args := []string{GateBook}
		//数据完整性校验
		ret, msg, data := common.GetZbs(content, args)
		if ret != base.StatusOK {
			fmt.Println("获取失败1 %s\n", msg)
			return base.ResError(int32(ret), msg)
		}
		resq = queryResp{ret, data}
	} else {
		args := []string{GateBook, getCert(stub)}

		//数据完整性校验
		ret, msg, data := common.GetOneZb(content, args)
		if ret != base.StatusOK {
			fmt.Println("获取失败2 %s\n", msg)
			return base.ResError(int32(ret), msg)
		}

		reqData := append([]string{}, data)
		resq = queryResp{ret, reqData}
	}

	return base.ResSuccessPayload(resq)
}

//修改网关
func ModGate(stub shim.ChaincodeStubInterface, strings []string) peer.Response {
	//日志
	fmt.Println("开始修改网关信息......")
	defer fmt.Println("修改网关信息结束......")

	//反序列化go
	var req map[string]interface{}
	if err := json.Unmarshal([]byte(strings[0]), &req); err != nil {
		fmt.Println("反序列化转换失败")
		return base.ResParamError()
	}

	validStrs := []string{ "GateName_s",  "Remark_s", "Status_i"}
	for k, _ := range req {
		if utils.InArry(k, validStrs) {
			continue
		} else {
			return base.ResParamError()
		}
	}

	cert := getCert(stub)

	if cert == "" {
		return base.ResBlockchainError()
	}
	args := []string{GateBook, cert}
	content := common.ComWorkSet{stub, nil, ledger.G_Books}

	//本接口默认只有权限方修改
	ret, msg := common.Set(content, args, req)
	if ret != base.StatusOK {
		return base.ResError(int32(ret), msg)
	}

	timeStr := time.Now().Format("2006-01-02 15:04:05")
	resq := modResp{ret, timeStr}

	return base.ResSuccessPayload(resq)
}

//成员注册
func UserRegister(stub shim.ChaincodeStubInterface, strings []string) peer.Response {
	//日志
	fmt.Println("开始创建创建成员......")
	defer fmt.Println("创建成员结束......")

	//反序列化go
	var req userRegisterReq
	if err := json.Unmarshal([]byte(strings[0]), &req); err != nil {
		fmt.Println("反序列化转换失败\n")
		return base.ResParamError()
	}
	if req.ApplyId_s == ""  {
		fmt.Println("输入参数错误\n")
		return base.ResParamError()
	}

	cert := getCert(stub)
	content := common.ComWorkSet{stub, nil, ledger.G_Books}
	ret, msg ,data := getGate(content ,cert)
	if ret != base.StatusOK {
		fmt.Println("无权限的网关接入 %s\n", msg)
		return base.ResError(int32(ret), msg)
	}

	var gateBook model.GateBook
	if err := json.Unmarshal([]byte(data), &gateBook); err != nil {
		fmt.Println("反序列化转换失败\n")
		return base.ResParamError()
	}

	timeStr := time.Now().Format("2006-01-02 15:04:05")
	book := model.MemberBook{req.ApplyId_s,  req.SecurityLevel_i, gateBook.Organization_s,req.Class_s,
		cert,req.Pubkey_s,req.UserCert_s,req.Role_s,req.Remark_s, timeStr, 1}

	args := []string{MemberBook, req.ApplyId_s}

	//数据完整性校验
	ret, msg = common.Creat(content, args, book)
	if ret != base.StatusOK {
		fmt.Println("创建失败 %s\n", msg)
		return base.ResError(int32(ret), msg)
	}

	resq := modResp{ret, timeStr}

	return base.ResSuccessPayload(resq)
}

//查询成员注册
func QueryRegister(stub shim.ChaincodeStubInterface, strings []string) peer.Response {
	//日志
	fmt.Println("开始查询成员1......")
	defer fmt.Println("查询成员结束1......")

	//反序列化go
	var req map[string]interface{}
	if err := json.Unmarshal([]byte(strings[0]), &req); err != nil {
		fmt.Println("反序列化转换失败")
		return base.ResParamError()
	}

	//检查参数是否超出范围
	validStrs := []string{"ApplyId_s", "UserId_s"}
	for k, _ := range req {
		if utils.InArry(k, validStrs) {
			continue
		} else {
			return base.ResParamError()
		}
	}

	user, ok := req["ApplyId_s"].(string)
	if !ok {
		return base.ResParamError()
	}
	userAuth,err := authority.GetUserAuth(stub,user)
	if err != nil {
		return base.ResBlockchainError()
	}
	if userAuth == nil {
		return base.ResPermissionError()
	}
	content := common.ComWorkSet{stub, userAuth, ledger.G_Books}


	var owner_Id = ""
	if id, Ok := req["UserId_s"].(string); Ok {
		owner_Id = id
	}

	var resq queryResp
	if owner_Id == "" {
		args := []string{MemberBook}
		//数据完整性校验

		ret, msg, data := common.GetZbs(content, args)
		if ret != base.StatusOK {
			fmt.Println("获取失败1 %s\n", msg)
			return base.ResError(int32(ret), msg)
		}
		resq = queryResp{ret, data}
	} else {
		args := []string{MemberBook, owner_Id}

		//数据完整性校验
		var data string
		ret, msg, data := common.GetOneZb(content, args)
		if ret != base.StatusOK {
			fmt.Println("获取失败2 %s\n", msg)
			return base.ResError(int32(ret), msg)
		}

		reqData := append([]string{}, data)
		resq = queryResp{ret, reqData}
	}

	return base.ResSuccessPayload(resq)
}

//修改注册成员信息
func ModRegister(stub shim.ChaincodeStubInterface, strings []string) peer.Response {
	//日志
	fmt.Println("开始修改注册成员信息......")
	defer fmt.Println("修改注册成员结束......")

	//反序列化go
	var req map[string]interface{}
	if err := json.Unmarshal([]byte(strings[0]), &req); err != nil {
		fmt.Println("反序列化转换失败")
		return base.ResParamError()
	}

	validStrs := []string{ "ApplyId_s", "Class_s", "Pubkey_s", "UserCert_s", "SecurityLevel_i", "Remark_s"}
	for k, _ := range req {
		if utils.InArry(k, validStrs) {
			continue
		} else {
			return base.ResParamError()
		}
	}

	user, ok := req["ApplyId_s"].(string)
	if !ok {
		return base.ResParamError()
	}
	userAuth,err := authority.GetUserAuth(stub,user)
	if err != nil {
		return base.ResBlockchainError()
	}
	if userAuth == nil {
		return base.ResPermissionError()
	}

	delete(req, "ApplyId_s")

	content := common.ComWorkSet{stub, userAuth, ledger.G_Books}
	args := []string{MemberBook, user}

	//本接口默认只有权限方修改
	ret, msg := common.Set(content, args, req)
	if ret != base.StatusOK {
		return base.ResError(int32(ret), msg)
	}

	timeStr := time.Now().Format("2006-01-02 15:04:05")
	resq := modResp{ret, timeStr}

	return base.ResSuccessPayload(resq)
}

//成员角色
func CreatRole(stub shim.ChaincodeStubInterface, strings []string) peer.Response {
	//日志
	fmt.Println("开始创建成员角色......")
	defer fmt.Println("创建成员角色结束......")

	//反序列化go
	var req creatRoleReq
	if err := json.Unmarshal([]byte(strings[0]), &req); err != nil {
		fmt.Println("反序列化转换失败\n")
		return base.ResParamError()
	}
	if req.RoleName_s == "" || len(req.OprList_a) == 0 {
		fmt.Println("输入参数错误\n")
		return base.ResParamError()
	}

	content := common.ComWorkSet{stub, nil, ledger.G_Books}
	book := model.RoleBook{req.RoleName_s,req.OprList_a,req.Remark_s,1}

	args := []string{RoleBook, req.RoleName_s}
	//数据完整性校验
	ret, msg := common.Creat(content, args, book)
	if ret != base.StatusOK {
		fmt.Println("创建失败 %s\n", msg)
		return base.ResError(int32(ret), msg)
	}

	timeStr := time.Now().Format("2006-01-02 15:04:05")
	resq := modResp{ret, timeStr}

	return base.ResSuccessPayload(resq)
}

//查询角色
func QueryRole(stub shim.ChaincodeStubInterface, strings []string) peer.Response {
	//日志
	fmt.Println("开始查询角色......")
	defer fmt.Println("查询角色结束......")

	//反序列化go
	var req map[string]interface{}
	if err := json.Unmarshal([]byte(strings[0]), &req); err != nil {
		fmt.Println("反序列化转换失败")
		return base.ResParamError()
	}

	//检查参数是否超出范围
	validStrs := []string{"ApplyId_s", "RoleName_s"}
	for k, _ := range req {
		if utils.InArry(k, validStrs) {
			continue
		} else {
			return base.ResParamError()
		}
	}

	user, ok := req["ApplyId_s"].(string)
	if !ok {
		return base.ResParamError()
	}
	userAuth,err := authority.GetUserAuth(stub,user)
	if err != nil {
		return base.ResBlockchainError()
	}
	if userAuth == nil {
		return base.ResPermissionError()
	}

	content := common.ComWorkSet{stub, userAuth, ledger.G_Books}

	roleName, ok:= req["RoleName_s"].(string)
	if !ok {
		return base.ResParamError()
	}

	var resq queryResp
	if roleName == "" {
		args := []string{RoleBook}

		var data []string
		ret, msg, data := common.GetZbs(content, args)
		if ret != base.StatusOK {
			fmt.Println("获取失败1 %s\n", msg)
			return base.ResError(int32(ret), msg)
		}
		resq = queryResp{ret, data}
	} else {
		args := []string{RoleBook, roleName}

		//数据完整性校验
		var data string
		ret, msg, data := common.GetOneZb(content, args)
		if ret != base.StatusOK {
			fmt.Println("获取失败2 %s\n", msg)
			return base.ResError(int32(ret), msg)
		}

		reqData := append([]string{}, data)
		resq = queryResp{ret, reqData}
	}

	return base.ResSuccessPayload(resq)
}

//修改角色
func ModRole(stub shim.ChaincodeStubInterface, strings []string) peer.Response {
	//日志
	fmt.Println("开始修改角色信息......")
	defer fmt.Println("修改角色信息结束......")

	//反序列化go
	var req map[string]interface{}
	if err := json.Unmarshal([]byte(strings[0]), &req); err != nil {
		fmt.Println("反序列化转换失败")
		return base.ResParamError()
	}

	validStrs := []string{ "ApplyId_s", "RoleName_s", "OprList_a", "Remark_s","Enable_i"}
	for k, _ := range req {
		if utils.InArry(k, validStrs) {
			continue
		} else {
			return base.ResParamError()
		}
	}

	user, ok := req["ApplyId_s"].(string)
	if !ok {
		return base.ResParamError()
	}
	userAuth,err := authority.GetUserAuth(stub,user)
	if err != nil {
		return base.ResBlockchainError()
	}
	if userAuth == nil {
		return base.ResPermissionError()
	}

	content := common.ComWorkSet{stub, userAuth, ledger.G_Books}

	var roleName = ""
	if id, Ok := req["RoleName_s"].(string); Ok {
		roleName = id
	}

	if roleName == "" {
		return base.ResParamError()
	}
	delete(req, "ApplyId_s")
	delete(req, "RoleName_s")


	args := []string{RoleBook, roleName}
	//本接口默认只有权限方修改
	ret, msg := common.Set(content, args, req)
	if ret != base.StatusOK {
		return base.ResError(int32(ret), msg)
	}

	timeStr := time.Now().Format("2006-01-02 15:04:05")
	resq := modResp{ret, timeStr}

	return base.ResSuccessPayload(resq)
}

//修改角色
func assignRole(stub shim.ChaincodeStubInterface, strings []string) peer.Response {
	//日志
	fmt.Println("开始修改注册成员信息......")
	defer fmt.Println("修改注册成员结束......")

	//反序列化go
	var req map[string]interface{}
	if err := json.Unmarshal([]byte(strings[0]), &req); err != nil {
		fmt.Println("反序列化转换失败")
		return base.ResParamError()
	}

	validStrs := []string{ "ApplyId_s", "Role_s", "UserId_s"}
	for k, _ := range req {
		if utils.InArry(k, validStrs) {
			continue
		} else {
			return base.ResParamError()
		}
	}

	user, ok := req["UserId_s"].(string)
	if !ok {
		return base.ResParamError()
	}

	delete(req, "ApplyId_s")
	delete(req, "UserId_s")

	content := common.ComWorkSet{stub, nil, ledger.G_Books}
	args := []string{MemberBook, user}

	//本接口默认只有权限方修改
	ret, msg := common.Set(content, args, req)
	if ret != base.StatusOK {
		return base.ResError(int32(ret), msg)
	}

	timeStr := time.Now().Format("2006-01-02 15:04:05")
	resq := modResp{ret, timeStr}

	return base.ResSuccessPayload(resq)
}

/********************
//版权登记
func RegisCopyright(stub shim.ChaincodeStubInterface, strings []string) peer.Response {
	//日志
	fmt.Println("开始版权登记......")
	defer fmt.Println("版权登记结束......")

	//反序列化go
	var req regisCopyrightReq
	if err := json.Unmarshal([]byte(strings[0]), &req); err != nil {
		fmt.Println("反序列化转换失败")
		return base.ResParamError()
	}
	if req.DataTable_s == "" || req.UserId_s == "" || req.SummaryHash_s == "" {
		return base.ResParamError()
	}

	content := common.ContentPointer{&stub, authority.G_AuthUser, ledger.G_Books}

	timeStr := time.Now().Format("2006-01-02 15:04:05")
	randStr := utils.GetRandomString()

	var resq regisCopyrightResp
	for {
		certNo := string(utils.Base58Encode([]byte(req.UserId_s + "." + randStr)))
		args := []string{CopyRightBook, certNo}
		period := "long-term"
		book := model.CopyRightBook{certNo, req.DataTable_s, req.SummaryHash_s, req.Type_s, req.UserId_s,
		 req.GenerateTime_s, req.Remark_s, period,
			InstitutionName, InstitutionCode, timeStr}
		//数据完整性校验
		ret, msg := common.Creat(content.Books, *content.Stub, args, book)
		if ret == base.StatusOK {
			resq = regisCopyrightResp{msg, certNo, period, InstitutionName, InstitutionCode, timeStr}
			break
		}
		if ret != base.StatusTheNewKeyRepeatedError {
			fmt.Println("创建失败 %s\n", msg)
			return base.ResError(int32(ret), msg)
		}
	}

	return base.ResSuccessPayload(resq)
}

//版权查询
func QueryCopyright(stub shim.ChaincodeStubInterface, strings []string) peer.Response {
	//日志
	fmt.Println("开始查询版权......")
	defer fmt.Println("查询版权结束......")

	//反序列化go
	var req map[string]string
	if err := json.Unmarshal([]byte(strings[0]), &req); err != nil {
		fmt.Println("反序列化转换失败")
		return base.ResParamError()
	}

	validStrs := []string{"CertNo_s", "DataTable_s", "UserId_s", "Type_s", "OwnerId_s", "Signature_s"}
	for k, _ := range req {
		if utils.InArry(k, validStrs) {
			continue
		} else {
			return base.ResParamError()
		}
	}

	content := common.ContentPointer{&stub, authority.G_AuthUser, ledger.G_Books}
	fmt.Println("\nreq", req)
	var user string
	if id, Ok := req["UserId_s"]; Ok {
		ret, msg,_ := getUserRegister(content, id)
		if ret != base.StatusOK {
			return base.ResError(int32(ret), msg)
		}
		user = id
	} else {
		return base.ResParamError()
	}
	delete(req, "UserId_s")

	if _, Ok := req["Signature_s"]; Ok {
		delete(req, "Signature_s")
	} else {
		fmt.Println("输入参数错误3\n")
		return base.ResParamError()
	}
	//把req修改为检索条件
	if id, Ok := req["OwnerId_s"]; Ok {
		req["UserId_s"] = id
		delete(req, "OwnerId_s")
	}

	var certNo_s = ""
	if _, ok := req["CertNo_s"]; ok {
		certNo_s = req["CertNo_s"]
	}
	fmt.Println("\ncertNo:", certNo_s)
	var resq queryResp
	if certNo_s != "" {
		args := []string{CopyRightBook, certNo_s}
		//获取账本
		ret, msg, data := common.GetOneZb(content, args, user)
		if ret != base.StatusOK {
			return base.ResError(int32(ret), msg)
		}
		reqData := append([]string{}, data)
		resq = queryResp{ret, reqData}

	} else {
		if len(req) == 0 {
			args := []string{CopyRightBook}
			//数据完整性校验
			ret, msg, data := common.GetZbs(content, args, user)
			if ret != base.StatusOK {
				fmt.Println("获取失败 %s\n", msg)
				return base.ResError(int32(ret), msg)
			}
			resq = queryResp{ret, data}
		} else {

			ret, msg, data := common.Query(content, CopyRightBook, req, user)
			if ret != base.StatusOK {
				fmt.Println("查询失败 %s\n", msg)
				return base.ResError(int32(ret), msg)
			}
			resq = queryResp{ret, data}
		}

	}

	return base.ResSuccessPayload(resq)
}
**********************/

//数据存证
func SaveEvd(stub shim.ChaincodeStubInterface, strings []string) peer.Response {
	//日志
	fmt.Println("开始存证......")
	defer fmt.Println("存证结束......")

	//反序列化go
	var req saveEvdReq
	if err := json.Unmarshal([]byte(strings[0]), &req); err != nil {
		fmt.Println("反序列化转换失败")
		return base.ResParamError()
	}
	if req.DataTable_s == "" || req.UserId_s == "" || req.Content_s == "" || req.StorType_s == ""  || req.BlockHash_s == ""{
		fmt.Println("输入参数错误\n")
		return base.ResParamError()
	}

	recordList := []string{}
	readlist := []string{}

	if req.StorType_s == "privacy" {
		if req.Indentify_s == ""{
			return base.ResParamError()
		}
	} else {
		readlist = req.ReadList_a
	}
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	handle := req.Handle_s
	if handle == "" {
		handle = string(utils.Base58Encode([]byte(req.UserId_s + "." + req.DataTable_s)))
	}

	args := []string{EvdBook, handle}
	content := common.ComWorkSet{stub, nil, ledger.G_Books}

	encHandle := ledger.ObscureKey(handle)
	book := model.EvdBook{encHandle , req.DataTable_s, req.UserId_s, req.Content_s, req.DataType_s,
		req.StorType_s,req.SecurityLevel_i,req.VisibleList_a,req.WriteList_a, readlist, req.DataFeature_s,
		req.AnalysiRule_s,recordList,req.Indentify_s,timeStr, req.BlockHash_s,req.Remark_s}


	if req.StorType_s != "privacy" {
		ret, _ := common.Creat(content, args, book)

		if ret != base.StatusOK {
			return base.ResParamError()
		}

	} else {
		//暂时不支持隐私区域存证
		return base.ResParamError()
	}

	resq := saveEvdResp{base.MsgSuccess, encHandle}
	return base.ResSuccessPayload(resq)
}

//存证列表查询
func QueryEvdList(stub shim.ChaincodeStubInterface, strings []string) peer.Response {
	//日志
	fmt.Println("开始查询存证......")
	defer fmt.Println("查询存证结束......")

	//反序列化go
	var req map[string]interface{}
	if err := json.Unmarshal([]byte(strings[0]), &req); err != nil {
		fmt.Println("反序列化转换失败")
		return base.ResParamError()
	}

	validStrs := []string{"ApplyId_s"}
	for k, _ := range req {
		if utils.InArry(k, validStrs) {
			continue
		} else {
			fmt.Println("输入参数错误1\n")
			return base.ResParamError()
		}
	}

	user, ok := req["ApplyId_s"].(string)
	if !ok {
		return base.ResParamError()
	}
	userAuth,err := authority.GetUserAuth(stub,user)
	if err != nil {
		return base.ResBlockchainError()
	}
	if userAuth == nil {
		return base.ResPermissionError()
	}

	content := common.ComWorkSet{stub, userAuth, ledger.G_Books}

	args := []string{EvdBook}
	ret, msg, data := common.GetZbs(content, args)
	if ret != base.StatusOK {
		fmt.Println("获取信息失败2 %s\n", msg)
		return base.ResError(int32(ret), msg)
	}

	var valData []string

	for _, t := range data {
		var bookContent model.EvdBook
		if err := json.Unmarshal([]byte(t), &bookContent); err != nil {
			fmt.Println("反序列化转换失败")
			return base.ResParamError()
		}
		ownerAuth,err := authority.GetUserAuth(stub,bookContent.UserId_s)
		if err != nil {
			return base.ResBlockchainError()
		}
		if userAuth == nil {
			continue
		}

		respItem := queryEvdListResp{bookContent.DataTable_s,bookContent.Handle_s,bookContent.DataType_s,bookContent.UserId_s,
			ownerAuth.Class,ownerAuth.Organization,bookContent.StorType_s,bookContent.Indentify_s,bookContent.DataFeature_s,
			bookContent.AnalysiRule_s,bookContent.Date_s,bookContent.BlockHash_s,bookContent.Remark_s}
		if tmp, err := json.Marshal(respItem);err != nil {
			valData = append(valData,string(tmp))
		}
	}
	resq := queryResp{ret, valData}

	return base.ResSuccessPayload(resq)
}

//存证查询
func QueryEvd(stub shim.ChaincodeStubInterface, strings []string) peer.Response {
	//日志
	fmt.Println("开始查询存证......")
	defer fmt.Println("查询存证结束......")

	//反序列化go
	var req map[string]interface{}
	if err := json.Unmarshal([]byte(strings[0]), &req); err != nil {
		fmt.Println("反序列化转换失败")
		return base.ResParamError()
	}

	validStrs := []string{"ApplyId_s", "Purpose_s", "Handle_s"}
	for k, _ := range req {
		if utils.InArry(k, validStrs) {
			continue
		} else {
			fmt.Println("输入参数错误1\n")
			return base.ResParamError()
		}
	}

	user, ok := req["ApplyId_s"].(string)
	if !ok {
		return base.ResParamError()
	}
	userAuth,err := authority.GetUserAuth(stub,user)
	if err != nil {
		return base.ResBlockchainError()
	}
	if userAuth == nil {
		return base.ResPermissionError()
	}

	content := common.ComWorkSet{stub, userAuth, ledger.G_Books}


	var handle_s = ""
	if id, ok := req["Handle_s"].(string); ok {
		handle_s = id
	}
	if handle_s == "" {
		fmt.Println("输入参数错误3\n")
		return base.ResParamError()
	}

	decHandle := ledger.VisualKey(handle_s)
	args := []string{EvdBook, decHandle}

	var resq queryResp
	ret, msg, data := common.GetOneZb(content, args)
	if ret != base.StatusOK {
		fmt.Println("获取信息失败1 %s\n", msg)
		return base.ResError(int32(ret), msg)
	}
	var dataModel model.EvdBook
	if err := json.Unmarshal([]byte(data), &dataModel); err != nil {
		fmt.Println("反序列化转换失败")
		return base.ResParamError()
	}

	if userAuth.SecurityAuthentication(dataModel.SecurityLevel_i) {
		return base.ResPermissionError()
	}

	//添加数据使用记录
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	book := model.EvdAccessBook{userAuth.User,"",userAuth.Class,timeStr,req["Purpose_s"].(string)}
	key := string(utils.Base58Encode([]byte(userAuth.User + "." + req["Purpose_s"].(string) + "." + timeStr)))
	argsTmp := []string{EvdAccessBook, key}
	ret,_ = common.Creat(content, argsTmp, book)

	if ret == base.StatusOK {
		var modData map[string]interface{}
		record := append(dataModel.RecordKey_a,key)
		modData["RecordKey_a"] = record
		ret, _ = common.Set(content, args, modData)

	}

	if ret != base.StatusOK {
		fmt.Println("获取信息失败3 %s\n", msg)
		return base.ResError(int32(ret), msg)
	}


	reqData := append([]string{}, data)
	resq = queryResp{ret, reqData}

	return base.ResSuccessPayload(resq)
}

//修改存证属性
func ModEvdAttr(stub shim.ChaincodeStubInterface, strings []string) peer.Response {
	//日志
	fmt.Println("开始修改存证属性......")
	defer fmt.Println("修改存证属性结束......")

	//反序列化go
	var req map[string]interface{}
	if err := json.Unmarshal([]byte(strings[0]), &req); err != nil {
		fmt.Println("反序列化转换失败")
		return base.ResParamError()
	}

	validStrs := []string{"ApplyId_s", "Handle_s", "Oprlist_a", "Readlist_a","SecurityLevel_i"}
	for k, _ := range req {
		if utils.InArry(k, validStrs) {
			continue
		} else {
			fmt.Println("输入参数错误\n")
			return base.ResParamError()
		}
	}

	user, ok := req["ApplyId_s"].(string)
	if !ok {
		return base.ResParamError()
	}
	userAuth,err := authority.GetUserAuth(stub,user)
	if err != nil {
		return base.ResBlockchainError()
	}
	if userAuth == nil {
		return base.ResPermissionError()
	}

	content := common.ComWorkSet{stub,userAuth , ledger.G_Books}

	var handle = ""
	if id, Ok := req["Handle_s"].(string); Ok {
		handle = id
	}
	if handle == "" {
		return base.ResParamError()
	}
	delete(req, "ApplyId_s")
	delete(req, "Handle_s")



	args := []string{EvdBook, handle}
	strs := append(args, "StorType_s")
	contentRead := common.ComWorkSet{stub,nil , ledger.G_Books}
	ret, msg, data := common.GetField(contentRead, strs)
	if data == "privacy" {
		fmt.Println("私有存证不允许修改\n")
		return base.ResPermissionError()
	}
	//数据完整性校验
	ret, msg = common.Set(content, args, req)
	if ret != base.StatusOK {
		fmt.Println("设置信息错误 %s\n", msg)
		return base.ResError(int32(ret), msg)
	}


	timeStr := time.Now().Format("2006-01-02 15:04:05")
	resq := modResp{ret, timeStr}

	return base.ResSuccessPayload(resq)
}



//申请资源交换
/**************************************************
func ApplyDataEx(stub shim.ChaincodeStubInterface, strings []string) peer.Response {
	//日志
	fmt.Println("开始申请资源交换......")
	defer fmt.Println("结束资源结束......")

	//反序列化go
	var req applyDataExchangeReq
	if err := json.Unmarshal([]byte(strings[0]), &req); err != nil {
		fmt.Println("反序列化转换失败")
		return base.ResParamError()
	}
	if req.UserId_s == "" || req.ResCode_s == "" || req.TargetId_s == "" {
		fmt.Println("输入参数错误\n")
		return base.ResParamError()
	}

	content := common.ContentPointer{&stub, authority.G_AuthUser, ledger.G_Books}
	ret, msg,_ := getUserRegister(content, req.UserId_s)
	if ret != base.StatusOK {
		fmt.Println("获取用户注册错误\n")
		return base.ResError(int32(ret), msg)
	}

	var resData string
	ret, msg, resData = getRes(content, req.UserId_s, req.ResCode_s)
	if ret != base.StatusOK {
		return base.ResError(int32(ret), msg)
	}
	var res model.ResBook
	if err := json.Unmarshal([]byte(resData), &res); err != nil {
		fmt.Println("反序列化转换失败")
		return base.ResParamError()
	}
	if res.UserId_s != req.UserId_s {
		return base.ResPermissionError()
	}
	taskNo := string(utils.Base58Encode([]byte(req.ResCode_s + "." + req.TargetId_s)))
	args := []string{DataExBook, taskNo}

	timeStr := time.Now().Format("2006-01-02 15:04:05")
	status := "0"
	if req.Status_s != "" {
		status = "1"
	}
	book := model.ExchangeBook{taskNo, req.ResCode_s, req.UserId_s,  req.ValidTimes_s,req.ValidTimes_s,req.TargetId_s,req.Remark_s,
		timeStr,"",status}

	ret, msg = common.Creat(content.Books, *content.Stub, args, book)
	if ret != base.StatusOK {
		fmt.Println("创建错误\n")
		return base.ResError(int32(ret), msg)
	}
	resq := applyDataExchangeResp{ret, taskNo,req.ValidPeriod_s,req.ValidTimes_s,timeStr,status}

	return base.ResSuccessPayload(resq)
}

//确认请求
func ConfirmEx(stub shim.ChaincodeStubInterface, strings []string) peer.Response {
	//日志
	fmt.Println("开始确认请求......")
	defer fmt.Println("确认请求结束......")

	//反序列化go
	var req ExApply
	if err := json.Unmarshal([]byte(strings[0]), &req); err != nil {
		fmt.Println("反序列化转换失败")
		return base.ResParamError()
	}

	content := common.ContentPointer{&stub, authority.G_AuthUser, ledger.G_Books}
	var data string
	ret, msg, data := getResExApply(content, req.UserId_s, req.TaskNo_s)
	if ret != base.StatusOK {
		return base.ResError(int32(ret), msg)
	}
	var apply model.ExchangeBook
	if err := json.Unmarshal([]byte(data), &apply); err != nil {
		fmt.Println("反序列化转换失败")
		return base.ResParamError()
	}
	if apply.UserId_s != req.UserId_s || apply.Status_s != "0"{
		return base.ResPermissionError()
	}

	args := []string{DataExBook, req.TaskNo_s}
	reqData := make(map[string]interface{})
	reqData["Status_s"] = "1"

	//数据完整性校验
	ret, msg = common.Set(content, args, reqData, req.UserId_s)
	if ret != base.StatusOK {
		fmt.Println("设置信息错误\n", msg)
		return base.ResError(int32(ret), msg)
	}

	timeStr := time.Now().Format("2006-01-02 15:04:05")
	resq := modResp{ret, timeStr}

	return base.ResSuccessPayload(resq)
}

//获取数据
func GetTaskData(stub shim.ChaincodeStubInterface, strings []string) peer.Response {
	//日志
	fmt.Println("开始获取请求......")
	defer fmt.Println("获取数据结束......")

	//反序列化go
	var req ExApply
	if err := json.Unmarshal([]byte(strings[0]), &req); err != nil {
		fmt.Println("反序列化转换失败")
		return base.ResParamError()
	}

	content := common.ContentPointer{&stub, authority.G_AuthUser, ledger.G_Books}
	var data string
	ret, msg, data := getResExApply(content, req.UserId_s, req.TaskNo_s)
	if ret != base.StatusOK {
		return base.ResError(int32(ret), msg)
	}
	var apply model.ExchangeBook
	if err := json.Unmarshal([]byte(data), &apply); err != nil {
		fmt.Println("反序列化转换失败")
		return base.ResParamError()
	}
	if apply.TargetId_s != req.UserId_s || apply.Status_s != "1"{
		return base.ResPermissionError()
	}
	var resData string
	ret, msg, resData = getRes(content, req.UserId_s, apply.ResCode_s)
	if ret != base.StatusOK {
		return base.ResError(int32(ret), msg)
	}
	var res model.ResBook
	if err := json.Unmarshal([]byte(resData), &res); err != nil {
		fmt.Println("反序列化转换失败")
		return base.ResParamError()
	}
	if res.Path_s == "" {
		ret = base.StatusDataPathNotSet
		msg = base.MsgDataPathNotSet
		return base.ResError(int32(ret), msg)
	}

	secretKey := ""
	if res.KeyHandle_s != "" {
		ret, msg, data = getEvd(content, apply.UserId_s,res.KeyHandle_s,res.Indentify_s)
		if ret != base.StatusOK {
			fmt.Println("获取存证错误\n")
			return base.ResError(int32(ret), msg)
		}
		var evd model.EvdBook
		if err := json.Unmarshal([]byte(data), &evd); err != nil {
			fmt.Println("反序列化转换失败")
			return base.ResParamError()
		}
		ret, msg,data = getUserRegister(content, req.UserId_s)
		if ret != base.StatusOK {
			fmt.Println("获取用户信息错误\n")
			return base.ResError(int32(ret), msg)
		}
		var member model.MemberBook
		if err := json.Unmarshal([]byte(data), &member); err != nil {
			fmt.Println("反序列化转换失败")
			return base.ResParamError()
		}
		encData,err := utils.RsaEncrypt([]byte(member.PublicKey_s),[]byte(evd.Contend_s))
		if err != nil {
			return base.ResBlockchainError()
		}
		secretKey = string(encData)
	}

	resq := ExResp{ret, res.Path_s,secretKey}

	return base.ResSuccessPayload(resq)
}
***************************************************************************/