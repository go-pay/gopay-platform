package dm

import "github.com/go-pay/xtime"

// SysRole 角色表
type SysRole struct {
	ID          int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Code        string     `gorm:"column:code" json:"code"`
	Name        string     `gorm:"column:name" json:"name"`
	Description string     `gorm:"column:description" json:"description"`
	BuiltIn     int8       `gorm:"column:built_in" json:"builtIn"`
	Status      int8       `gorm:"column:status" json:"status"`
	Ctime       xtime.Time `gorm:"column:ctime;autoCreateTime" json:"ctime"`
	Utime       xtime.Time `gorm:"column:utime;autoUpdateTime" json:"utime"`
}

func (SysRole) TableName() string { return "sys_role" }

// SysRolePerm 角色权限表
type SysRolePerm struct {
	ID     int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RoleID int64      `gorm:"column:role_id" json:"roleId"`
	Perm   string     `gorm:"column:perm" json:"perm"`
	Ctime  xtime.Time `gorm:"column:ctime;autoCreateTime" json:"ctime"`
}

func (SysRolePerm) TableName() string { return "sys_role_perm" }

// OperationLog 操作日志表
type OperationLog struct {
	ID          int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Operator    string     `gorm:"column:operator" json:"operator"`
	Module      string     `gorm:"column:module" json:"module"`
	Action      string     `gorm:"column:action" json:"action"`
	Description string     `gorm:"column:description" json:"description"`
	IP          string     `gorm:"column:ip" json:"ip"`
	UserAgent   string     `gorm:"column:user_agent" json:"userAgent"`
	Success     int8       `gorm:"column:success" json:"success"`
	Duration    int        `gorm:"column:duration" json:"duration"`
	RequestData string     `gorm:"column:request_data" json:"requestData"`
	Ctime       xtime.Time `gorm:"column:ctime;autoCreateTime" json:"ctime"`
}

func (OperationLog) TableName() string { return "operation_log" }
