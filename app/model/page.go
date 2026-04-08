package model

// PageReq 统一分页请求
type PageReq struct {
	Page     int `json:"page" binding:"required,min=1"`
	PageSize int `json:"pageSize" binding:"required,min=1,max=100"`
}

// PageResp 统一分页响应
type PageResp struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}
