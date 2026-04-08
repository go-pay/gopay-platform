package dao

import (
	"context"
	"errors"

	"gopay/app/dm"
	"gopay/app/model"

	"gorm.io/gorm"
)

// CreateOperationLog 创建操作日志
func (d *Dao) CreateOperationLog(log *dm.OperationLog) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.Create(log).Error
}

// LogList 操作日志列表（分页）
func (d *Dao) LogList(ctx context.Context, req *model.LogListReq) (list []*dm.OperationLog, total int64, err error) {
	if d.GopayDB == nil {
		return nil, 0, ErrNoDatabase
	}
	db := d.GopayDB.WithContext(ctx).Model(&dm.OperationLog{})
	if req.Operator != "" {
		db = db.Where("operator LIKE ?", "%"+req.Operator+"%")
	}
	if req.Module != "" {
		db = db.Where("module = ?", req.Module)
	}
	if req.Action != "" {
		db = db.Where("action = ?", req.Action)
	}
	if req.Date != "" {
		db = db.Where("DATE(ctime) = ?", req.Date)
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Order("id DESC").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error
	return
}

// GetLogByID 获取操作日志详情
func (d *Dao) GetLogByID(ctx context.Context, id int64) (*dm.OperationLog, error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	log := new(dm.OperationLog)
	err := d.GopayDB.WithContext(ctx).Where("id = ?", id).First(log).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return log, err
}

// LogListAll 操作日志列表（不分页，用于导出）
func (d *Dao) LogListAll(ctx context.Context, req *model.LogListReq) (list []*dm.OperationLog, err error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	db := d.GopayDB.WithContext(ctx).Model(&dm.OperationLog{})
	if req.Operator != "" {
		db = db.Where("operator LIKE ?", "%"+req.Operator+"%")
	}
	if req.Module != "" {
		db = db.Where("module = ?", req.Module)
	}
	if req.Action != "" {
		db = db.Where("action = ?", req.Action)
	}
	if req.Date != "" {
		db = db.Where("DATE(ctime) = ?", req.Date)
	}
	err = db.Order("id DESC").Limit(10000).Find(&list).Error
	return
}
