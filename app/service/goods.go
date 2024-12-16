package service

import (
	"context"
	"gopay/app/model"
)

// GetGoodsList 支付码获取(用户扫码支付)
func (s *Service) GetGoodsList(c context.Context) (rsp *model.GetGoodsListRsp, err error) {
	// DB 查询商品列表

	// 模拟数据返回
	rsp = &model.GetGoodsListRsp{
		List: []*model.GoodsInfo{
			{
				GoodsId:   1,
				GoodsName: "商品1",
				GoodsDesc: "我是商品1",
				UnitPrice: 1,
			},
			{
				GoodsId:   2,
				GoodsName: "商品2",
				GoodsDesc: "我是商品2",
				UnitPrice: 1,
			},
			{
				GoodsId:   3,
				GoodsName: "商品3",
				GoodsDesc: "我是商品3",
				UnitPrice: 1,
			},
		},
	}
	return rsp, nil
}
