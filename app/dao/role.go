package dao

import (
	"context"
	"errors"

	"gopay/app/dm"
	"gopay/app/model"

	"gorm.io/gorm"
)

// RoleList 角色列表
func (d *Dao) RoleList(ctx context.Context, req *model.RoleListReq) (list []*dm.SysRole, total int64, err error) {
	if d.GopayDB == nil {
		return nil, 0, ErrNoDatabase
	}
	db := d.GopayDB.WithContext(ctx).Model(&dm.SysRole{})
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Order("id ASC").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error
	return
}

// GetRoleByID 根据 ID 查询角色
func (d *Dao) GetRoleByID(ctx context.Context, id int64) (*dm.SysRole, error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	role := new(dm.SysRole)
	err := d.GopayDB.WithContext(ctx).Where("id = ?", id).First(role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return role, err
}

// GetRoleByCode 根据 code 查询角色
func (d *Dao) GetRoleByCode(ctx context.Context, code string) (*dm.SysRole, error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	role := new(dm.SysRole)
	err := d.GopayDB.WithContext(ctx).Where("code = ?", code).First(role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return role, err
}

// CreateRole 创建角色
func (d *Dao) CreateRole(ctx context.Context, role *dm.SysRole) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Create(role).Error
}

// UpdateRole 更新角色
func (d *Dao) UpdateRole(ctx context.Context, id int64, updates map[string]interface{}) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Model(&dm.SysRole{}).Where("id = ?", id).Updates(updates).Error
}

// ToggleRoleStatus 切换角色状态
func (d *Dao) ToggleRoleStatus(ctx context.Context, id int64) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Exec("UPDATE sys_role SET status = IF(status=1,0,1) WHERE id = ?", id).Error
}

// GetRolePerms 获取角色权限列表
func (d *Dao) GetRolePerms(ctx context.Context, roleID int64) ([]string, error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	var perms []string
	err := d.GopayDB.WithContext(ctx).Model(&dm.SysRolePerm{}).Where("role_id = ?", roleID).Pluck("perm", &perms).Error
	return perms, err
}

// UpdateRolePerms 更新角色权限（全量替换）
func (d *Dao) UpdateRolePerms(ctx context.Context, roleID int64, perms []string) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&dm.SysRolePerm{}).Error; err != nil {
			return err
		}
		if len(perms) == 0 {
			return nil
		}
		records := make([]*dm.SysRolePerm, 0, len(perms))
		for _, p := range perms {
			records = append(records, &dm.SysRolePerm{RoleID: roleID, Perm: p})
		}
		return tx.Create(&records).Error
	})
}
