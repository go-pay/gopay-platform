package dao

import (
	"context"
	"gopay/app/dm"
)

func (d *Dao) GoodsById(c context.Context, goodsId int) (sg *dm.ShopGoods, err error) {
	sg = new(dm.ShopGoods)
	err = d.BizDB.Select([]string{"id", "developer_id", "appid", "request_order_url", "report_callback_url", "report_type"}).
		Table(sg.TableName()).
		Where("developer_id = ?", goodsId).
		Take(sg).Error
	if err != nil {
		return nil, err
	}
	return sg, nil
}
