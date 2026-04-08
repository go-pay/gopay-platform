package errcode

import "github.com/go-pay/ecode"

var (
	New = ecode.New

	Success = New(0, "SUCCESS", "success")

	// service error
	RequestErr         = New(10400, "PARAM_ERROR", "parameter error")
	HeaderVerifyFailed = New(10401, "HEADER_ERROR", "header verify failed")
	LoginFailed        = New(10402, "LOGIN_FAILED", "用户名或密码错误")
	TokenInvalid       = New(10403, "TOKEN_INVALID", "登录已过期，请重新登录")
	ServerErr          = New(10500, "SERVER_ERROR", "server error")
	UnAvailableErr     = New(10501, "UNAVAILABLE", "unavailable")
)
