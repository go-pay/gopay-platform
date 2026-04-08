package model

// ---- 用户管理 ----

type UserListReq struct {
	PageReq
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Status   int8   `json:"status"` // -1=全部
}

type UserAddReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	RealName string `json:"realName"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type UserUpdateReq struct {
	ID       int64  `json:"id" binding:"required"`
	RealName string `json:"realName"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type IDReq struct {
	ID int64 `json:"id" binding:"required"`
}

// ---- 角色管理 ----

type RoleListReq struct {
	PageReq
}

type RoleAddReq struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type RoleUpdateReq struct {
	ID          int64  `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type RolePermsUpdateReq struct {
	RoleID int64    `json:"roleId" binding:"required"`
	Perms  []string `json:"perms"`
}

type RolePermsListReq struct {
	RoleID int64 `json:"roleId" binding:"required"`
}

// ---- 操作日志 ----

type LogListReq struct {
	PageReq
	Operator string `json:"operator"`
	Module   string `json:"module"`
	Action   string `json:"action"`
	Date     string `json:"date"`
}
