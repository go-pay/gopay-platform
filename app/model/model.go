package model

// LoginReq 登录请求
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginRsp 登录响应
type LoginRsp struct {
	Token    string    `json:"token"`
	UserInfo *UserInfo `json:"userInfo"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	RealName  string `json:"realName"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	LastLogin string `json:"lastLogin"`
}

// ChangePwdReq 修改密码请求
type ChangePwdReq struct {
	OldPassword     string `json:"oldPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

// ProfileReq 更新个人资料请求
type ProfileReq struct {
	RealName string `json:"realName"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}
