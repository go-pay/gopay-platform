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
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}
