package model

type GetGoodsListRsp struct {
	List []*GoodsInfo `json:"list"`
}

type GoodsInfo struct {
	GoodsId   int    `json:"goods_id"`
	GoodsName string `json:"goods_name"`
	GoodsDesc string `json:"goods_desc"`
	UnitPrice int    `json:"unit_price"`
}
