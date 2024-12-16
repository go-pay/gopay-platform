package ecode

import "github.com/go-pay/ecode"

const (
	StatusCodeNotOK = "STATUS_CODE_NOT_200"
	RspCodeNotOK    = "RESPONSE_CODE_NOT_SUCCESS"
)

var (
	New = ecode.New

	// service error
	RequestErr         = New(10400, "PARAM_ERROR", "parameter error")
	HeaderVerifyFailed = New(10401, "HEADER_ERROR", "header verify failed")
	ServerErr          = New(10500, "SERVER_ERROR", "server error")
	UnAvailableErr     = New(10501, "UNAVAILABLE", "unavailable")

	WxNativePayErr = func(msg string) error {
		return New(10001, "WX_NATIVE_PAY_ERR", "wechat native pay failed: "+msg)
	}
)
