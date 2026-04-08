package model

// ---- 商户管理 ----

type MerchantListReq struct {
	PageReq
	Name    string `json:"name"`
	Contact string `json:"contact"`
	Status  int8   `json:"status"` // -1=全部
}

type MerchantAddReq struct {
	Name    string `json:"name" binding:"required"`
	Contact string `json:"contact" binding:"required"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Remark  string `json:"remark"`
}

type MerchantUpdateReq struct {
	ID      int64  `json:"id" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Contact string `json:"contact" binding:"required"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Remark  string `json:"remark"`
}

// ---- 商户应用 ----

type MerchantAppListReq struct {
	PageReq
	Name         string `json:"name"`
	Appid        string `json:"appid"`
	PlatformType int8   `json:"platformType"` // -1=全部
	MerchantID   int64  `json:"merchantId"`
}

type MerchantAppAddReq struct {
	Name         string `json:"name" binding:"required"`
	Appid        string `json:"appid" binding:"required"`
	MerchantID   int64  `json:"merchantId" binding:"required"`
	PlatformType int8   `json:"platformType"`
	MerchantType int8   `json:"merchantType"`
	NotifyUrl    string `json:"notifyUrl"`
	ReturnUrl    string `json:"returnUrl"`
}

type MerchantAppUpdateReq struct {
	ID           int64  `json:"id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	MerchantID   int64  `json:"merchantId" binding:"required"`
	PlatformType int8   `json:"platformType"`
	MerchantType int8   `json:"merchantType"`
	NotifyUrl    string `json:"notifyUrl"`
	ReturnUrl    string `json:"returnUrl"`
}
