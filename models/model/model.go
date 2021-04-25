package model

//账本与请求模型
const (
//账本名称
)

//账本基础共用模型
type BaseModel struct {

}

//返回值基础模型
type ResList struct {
	List  interface{} `json:"list"`
	Total int         `json:"total"`
}
const(
	//成员状态
	MemberStatusBan = 0			//禁用
	MemberStatusUsing = 1		//使用中

	//工作状态
	TaskInit = 0               //任务初始状态
	TaskDoing = 1              //任务执行中
	TaskSuccess = 2            //任务执行成功
	TaskFailed = 3             //任务执行失败
)

//企业成员注册账本
type GateBook struct {
	MspId_s         string         //网关标识
	Cert_s          string         //网关节点证书
	GateName_s		string         //网关名称
	Admin_s         string         //管理员名称
	AdminCert_s     string         //管理员证书
	Organization_s  string         //所属单位
	Remark_s        string         //简介
	Status_i        int            //状态0:禁用，其他：使用中
}

//企业成员注册账本
type MemberBook struct {
	UserId_s        string         //注册成员ID
	SecurityLevel_i int            //密级
	Organization    string         //组织
	Class_s         string         //部门
	GateCert_s		string         //网关标识，网关证书
	Pubkey_s        string         //用户公钥
	UserCert_s      string         //用户证书
	Role_s          string         //成员角色
	Remark_s        string         //成员介绍
	Date_s          string         //登记日期（年/月/日/时）
	Status_i        int            //状态0:禁用，其他：使用中
}

//角色账本
type RoleBook struct {
	RoleName_s      string        //角色名称
	OprList_a       []string      //权限列表
	Remark_s	    string        //角色介绍
	Enable_i        int           //0作废，1启用
}

//版权账本
type CopyRightBook struct {
	CertNo_s	    string        //证书编号
	DataTable_s	    string        //数据标签
	SummaryHash_s	string        //摘要hash
	Type_s	        string        //摘要算法类型
	UserId_s	    string        //数据拥有方ID
	GenerateTime_s	string        //数据生成日期（年/月/日/时）
	Remark_s	    string        //数据特点介绍
	Subject_s       string        //学科
	AppFields_s     string        //应用领域
	Suitable_s      string        //适用对象
	RecordKey_s     string        //访问记录地址
	Institution_s	string        //区块链认证机构
	Code_s	        string        //认证机构社会唯一信用码
	Date_s	        string        //认证日期（年/月/日/时）
}

//数据存证账本
type EvdBook struct {
	Handle_s        string	     //存证句柄
	DataTable_s	    string	     //数据标签
	UserId_s	    string	     //数据拥有方ID
	Content_s	    string	     //数据内容
	DataType_s      string	     //数据类型。数据存证；文件存证
	StorType_s      string       //存储类型。”normal”:默认普通存储；”privacy”：隐私存储；
	SecurityLevel_i int          //密级
	VisibleList_a   []string     //可见列表
	WriteList_a       []string    //编辑列表
	ReadList_a      []string	 //读取列表
	DataFeature_s   string       //数据特征
	AnalysiRule_s   string       //数据规则
	RecordKey_a     []string      //访问记录地址
	Indentify_s	    string	     //隐私区域标识(隐私存证有效)
	Date_s          string       //存证时间
	BlockHash_s     string       //存证数据hash
	Remark_s	    string	     //存证数据描述
}

//存证访问记录账本
type EvdAccessBook struct {
	Accesser_s     string       //数据访问者
	Organization_s string       //访问单位
	Class_s        string       //访问部门
	AccessTime_s   string       //访问时间
	Purpose_s      string       //用途
}
