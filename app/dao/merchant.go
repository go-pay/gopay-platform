package dao

import (
	"context"
	"errors"

	"gopay/app/dm"
	"gopay/app/model"

	"gorm.io/gorm"
)

// ---- 商户 ----

// MerchantList 商户列表
func (d *Dao) MerchantList(ctx context.Context, req *model.MerchantListReq) (list []*dm.Merchant, total int64, err error) {
	if d.GopayDB == nil {
		return nil, 0, ErrNoDatabase
	}
	db := d.GopayDB.WithContext(ctx).Model(&dm.Merchant{})
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Contact != "" {
		db = db.Where("contact LIKE ?", "%"+req.Contact+"%")
	}
	if req.Status >= 0 {
		db = db.Where("status = ?", req.Status)
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Order("id DESC").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error
	return
}

// GetMerchantByID 根据 ID 查询商户
func (d *Dao) GetMerchantByID(ctx context.Context, id int64) (*dm.Merchant, error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	m := new(dm.Merchant)
	err := d.GopayDB.WithContext(ctx).Where("id = ?", id).First(m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return m, err
}

// CreateMerchant 创建商户
func (d *Dao) CreateMerchant(ctx context.Context, m *dm.Merchant) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Create(m).Error
}

// UpdateMerchant 更新商户
func (d *Dao) UpdateMerchant(ctx context.Context, id int64, updates map[string]interface{}) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Model(&dm.Merchant{}).Where("id = ?", id).Updates(updates).Error
}

// ToggleMerchantStatus 切换商户状态
func (d *Dao) ToggleMerchantStatus(ctx context.Context, id int64) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Exec("UPDATE merchant SET status = IF(status=1,0,1) WHERE id = ?", id).Error
}

// MerchantOptions 商户下拉选项（仅启用的商户）
func (d *Dao) MerchantOptions(ctx context.Context) (list []*dm.Merchant, err error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	err = d.GopayDB.WithContext(ctx).Select("id, name").Where("status = 1").Order("id ASC").Find(&list).Error
	return
}

// GetMerchantsByIDs 批量查询商户名称
func (d *Dao) GetMerchantsByIDs(ctx context.Context, ids []int64) (map[int64]string, error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	if len(ids) == 0 {
		return make(map[int64]string), nil
	}
	var merchants []*dm.Merchant
	if err := d.GopayDB.WithContext(ctx).Select("id, name").Where("id IN ?", ids).Find(&merchants).Error; err != nil {
		return nil, err
	}
	nameMap := make(map[int64]string, len(merchants))
	for _, m := range merchants {
		nameMap[m.ID] = m.Name
	}
	return nameMap, nil
}

// ---- 商户应用 ----

// MerchantAppList 商户应用列表
func (d *Dao) MerchantAppList(ctx context.Context, req *model.MerchantAppListReq) (list []*dm.MerchantApp, total int64, err error) {
	if d.GopayDB == nil {
		return nil, 0, ErrNoDatabase
	}
	db := d.GopayDB.WithContext(ctx).Model(&dm.MerchantApp{})
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Appid != "" {
		db = db.Where("appid LIKE ?", "%"+req.Appid+"%")
	}
	if req.PlatformType >= 0 {
		db = db.Where("platform_type = ?", req.PlatformType)
	}
	if req.MerchantID > 0 {
		db = db.Where("merchant_id = ?", req.MerchantID)
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Order("id DESC").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error
	return
}

// GetMerchantAppByID 根据 ID 查询商户应用
func (d *Dao) GetMerchantAppByID(ctx context.Context, id int64) (*dm.MerchantApp, error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	app := new(dm.MerchantApp)
	err := d.GopayDB.WithContext(ctx).Where("id = ?", id).First(app).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return app, err
}

// CreateMerchantApp 创建商户应用
func (d *Dao) CreateMerchantApp(ctx context.Context, app *dm.MerchantApp) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Create(app).Error
}

// UpdateMerchantApp 更新商户应用
func (d *Dao) UpdateMerchantApp(ctx context.Context, id int64, updates map[string]interface{}) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Model(&dm.MerchantApp{}).Where("id = ?", id).Updates(updates).Error
}
