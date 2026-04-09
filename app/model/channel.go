package model

import "github.com/go-pay/xtime"

// ---- 支付通道 ----

type ChannelListReq struct {
	PageReq
	Name   string `json:"name"`
	Code   string `json:"code"`
	Type   string `json:"type"`   // ""=全部
	Status int8   `json:"status"` // -1=全部
}

type ChannelAddReq struct {
	Name       string   `json:"name" binding:"required"`
	Code       string   `json:"code" binding:"required"`
	Type       string   `json:"type" binding:"required"`
	MerchantID int64    `json:"merchantId" binding:"required"`
	PayMethods []string `json:"payMethods" binding:"required"`
	FeeRate    float64  `json:"feeRate"`
	Remark     string   `json:"remark"`
}

type ChannelUpdateReq struct {
	ID         int64    `json:"id" binding:"required"`
	Name       string   `json:"name" binding:"required"`
	Type       string   `json:"type" binding:"required"`
	MerchantID int64    `json:"merchantId" binding:"required"`
	PayMethods []string `json:"payMethods" binding:"required"`
	FeeRate    float64  `json:"feeRate"`
	Remark     string   `json:"remark"`
}

type ChannelConfigReq struct {
	ChannelID  int64  `json:"channelId" binding:"required"`
	AppID      string `json:"appId"`
	MchID      string `json:"mchId"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
	ApiKey     string `json:"apiKey"`
	SerialNo   string `json:"serialNo"`
	NotifyUrl  string `json:"notifyUrl"`
	SignType   string `json:"signType"`
	Sandbox    bool   `json:"sandbox"`
}

// ChannelResp 通道响应（payMethods 转为数组）
type ChannelResp struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name"`
	Code         string     `json:"code"`
	Type         string     `json:"type"`
	MerchantID   int64      `json:"merchantId"`
	MerchantName string     `json:"merchantName"`
	PayMethods   []string   `json:"payMethods"`
	FeeRate      float64    `json:"feeRate"`
	Status       int8       `json:"status"`
	Remark       string     `json:"remark"`
	Ctime        xtime.Time `json:"ctime"`
}

// ChannelDetailResp 通道详情（含配置）
type ChannelDetailResp struct {
	*ChannelResp
	Config *ChannelConfigResp `json:"config"`
}

// ChannelConfigResp 通道配置响应（敏感字段脱敏）
type ChannelConfigResp struct {
	AppID      string `json:"appId"`
	MchID      string `json:"mchId"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
	ApiKey     string `json:"apiKey"`
	SerialNo   string `json:"serialNo"`
	NotifyUrl  string `json:"notifyUrl"`
	SignType   string `json:"signType"`
	Sandbox    int8   `json:"sandbox"`
}
