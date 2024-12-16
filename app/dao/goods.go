package dao

import (
	"context"
	"gopay/app/dm"
)

func (d *Dao) GoodsById(c context.Context, goodsId int) (sg *dm.ShopGoods, err error) {
	sg = new(dm.ShopGoods)
	err = d.BizDB.Select([]string{"id", "sku_id", "goods_name", "goods_desc", "unit_price"}).
		Table(sg.TableName()).
		Where("id = ?", goodsId).
		Take(sg).Error
	if err != nil {
		return nil, err
	}
	return sg, nil
}
