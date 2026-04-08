package service

import (
	"context"

	"gopay/app/dm"
	"gopay/app/model"
	ec "gopay/errcode"

	"github.com/go-pay/xlog"
	"golang.org/x/crypto/bcrypt"
)

// ---- 用户管理 ----

// SystemUserList 用户列表
func (s *Service) SystemUserList(ctx context.Context, req *model.UserListReq) (*model.PageResp, error) {
	list, total, err := s.dao.AccountList(ctx, req)
	if err != nil {
		xlog.Errorf("SystemUserList AccountList, err:%v", err)
		return nil, ec.ServerErr
	}
	return &model.PageResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// SystemUserAdd 新增用户
func (s *Service) SystemUserAdd(ctx context.Context, req *model.UserAddReq) error {
	// 检查用户名是否已存在
	existing, err := s.dao.GetAccountByUname(ctx, req.Username)
	if err != nil {
		xlog.Errorf("SystemUserAdd GetAccountByUname(%s), err:%v", req.Username, err)
		return ec.ServerErr
	}
	if existing != nil {
		return ec.Conflict
	}
	// bcrypt 加密密码
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		xlog.Errorf("SystemUserAdd bcrypt.GenerateFromPassword, err:%v", err)
		return ec.ServerErr
	}
	account := &dm.Account{
		Uname:    req.Username,
		Pwd:      string(hash),
		RealName: req.RealName,
		Phone:    req.Phone,
		Email:    req.Email,
		Role:     req.Role,
		Status:   1,
	}
	if err = s.dao.CreateAccount(ctx, account); err != nil {
		xlog.Errorf("SystemUserAdd CreateAccount, err:%v", err)
		return ec.ServerErr
	}
	return nil
}

// SystemUserUpdate 编辑用户
func (s *Service) SystemUserUpdate(ctx context.Context, req *model.UserUpdateReq) error {
	account, err := s.dao.GetAccountByID(ctx, req.ID)
	if err != nil {
		xlog.Errorf("SystemUserUpdate GetAccountByID(%d), err:%v", req.ID, err)
		return ec.ServerErr
	}
	if account == nil {
		return ec.NotFound
	}
	updates := map[string]interface{}{
		"real_name": req.RealName,
		"phone":     req.Phone,
		"email":     req.Email,
		"role":      req.Role,
	}
	if err = s.dao.UpdateAccount(ctx, req.ID, updates); err != nil {
		xlog.Errorf("SystemUserUpdate UpdateAccount(%d), err:%v", req.ID, err)
		return ec.ServerErr
	}
	return nil
}

// SystemUserToggleStatus 切换用户状态
func (s *Service) SystemUserToggleStatus(ctx context.Context, id int64) error {
	account, err := s.dao.GetAccountByID(ctx, id)
	if err != nil {
		xlog.Errorf("SystemUserToggleStatus GetAccountByID(%d), err:%v", id, err)
		return ec.ServerErr
	}
	if account == nil {
		return ec.NotFound
	}
	if err = s.dao.ToggleAccountStatus(ctx, id); err != nil {
		xlog.Errorf("SystemUserToggleStatus ToggleAccountStatus(%d), err:%v", id, err)
		return ec.ServerErr
	}
	return nil
}

// SystemUserResetPwd 重置密码（返回新密码）
func (s *Service) SystemUserResetPwd(ctx context.Context, id int64) (string, error) {
	account, err := s.dao.GetAccountByID(ctx, id)
	if err != nil {
		xlog.Errorf("SystemUserResetPwd GetAccountByID(%d), err:%v", id, err)
		return "", ec.ServerErr
	}
	if account == nil {
		return "", ec.NotFound
	}
	newPwd := "123456"
	hash, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		xlog.Errorf("SystemUserResetPwd bcrypt.GenerateFromPassword, err:%v", err)
		return "", ec.ServerErr
	}
	if err = s.dao.UpdateAccountPwd(ctx, id, string(hash)); err != nil {
		xlog.Errorf("SystemUserResetPwd UpdateAccountPwd(%d), err:%v", id, err)
		return "", ec.ServerErr
	}
	return newPwd, nil
}

// ---- 角色管理 ----

// SystemRoleList 角色列表
func (s *Service) SystemRoleList(ctx context.Context, req *model.RoleListReq) (*model.PageResp, error) {
	list, total, err := s.dao.RoleList(ctx, req)
	if err != nil {
		xlog.Errorf("SystemRoleList RoleList, err:%v", err)
		return nil, ec.ServerErr
	}
	return &model.PageResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// SystemRoleAdd 新增角色
func (s *Service) SystemRoleAdd(ctx context.Context, req *model.RoleAddReq) error {
	existing, err := s.dao.GetRoleByCode(ctx, req.Code)
	if err != nil {
		xlog.Errorf("SystemRoleAdd GetRoleByCode(%s), err:%v", req.Code, err)
		return ec.ServerErr
	}
	if existing != nil {
		return ec.Conflict
	}
	role := &dm.SysRole{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Status:      1,
	}
	if err = s.dao.CreateRole(ctx, role); err != nil {
		xlog.Errorf("SystemRoleAdd CreateRole, err:%v", err)
		return ec.ServerErr
	}
	return nil
}

// SystemRoleUpdate 编辑角色
func (s *Service) SystemRoleUpdate(ctx context.Context, req *model.RoleUpdateReq) error {
	role, err := s.dao.GetRoleByID(ctx, req.ID)
	if err != nil {
		xlog.Errorf("SystemRoleUpdate GetRoleByID(%d), err:%v", req.ID, err)
		return ec.ServerErr
	}
	if role == nil {
		return ec.NotFound
	}
	if role.BuiltIn == 1 {
		return ec.Forbidden
	}
	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
	}
	if err = s.dao.UpdateRole(ctx, req.ID, updates); err != nil {
		xlog.Errorf("SystemRoleUpdate UpdateRole(%d), err:%v", req.ID, err)
		return ec.ServerErr
	}
	return nil
}

// SystemRoleToggleStatus 切换角色状态
func (s *Service) SystemRoleToggleStatus(ctx context.Context, id int64) error {
	role, err := s.dao.GetRoleByID(ctx, id)
	if err != nil {
		xlog.Errorf("SystemRoleToggleStatus GetRoleByID(%d), err:%v", id, err)
		return ec.ServerErr
	}
	if role == nil {
		return ec.NotFound
	}
	if role.BuiltIn == 1 {
		return ec.Forbidden
	}
	if err = s.dao.ToggleRoleStatus(ctx, id); err != nil {
		xlog.Errorf("SystemRoleToggleStatus ToggleRoleStatus(%d), err:%v", id, err)
		return ec.ServerErr
	}
	return nil
}

// SystemRolePermsUpdate 更新角色权限
func (s *Service) SystemRolePermsUpdate(ctx context.Context, req *model.RolePermsUpdateReq) error {
	role, err := s.dao.GetRoleByID(ctx, req.RoleID)
	if err != nil {
		xlog.Errorf("SystemRolePermsUpdate GetRoleByID(%d), err:%v", req.RoleID, err)
		return ec.ServerErr
	}
	if role == nil {
		return ec.NotFound
	}
	if err = s.dao.UpdateRolePerms(ctx, req.RoleID, req.Perms); err != nil {
		xlog.Errorf("SystemRolePermsUpdate UpdateRolePerms(%d), err:%v", req.RoleID, err)
		return ec.ServerErr
	}
	return nil
}

// SystemRolePermsList 获取角色权限列表
func (s *Service) SystemRolePermsList(ctx context.Context, roleID int64) ([]string, error) {
	perms, err := s.dao.GetRolePerms(ctx, roleID)
	if err != nil {
		xlog.Errorf("SystemRolePermsList GetRolePerms(%d), err:%v", roleID, err)
		return nil, ec.ServerErr
	}
	return perms, nil
}

// ---- 操作日志 ----

// SystemLogList 操作日志列表
func (s *Service) SystemLogList(ctx context.Context, req *model.LogListReq) (*model.PageResp, error) {
	list, total, err := s.dao.LogList(ctx, req)
	if err != nil {
		xlog.Errorf("SystemLogList LogList, err:%v", err)
		return nil, ec.ServerErr
	}
	return &model.PageResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// SystemLogDetail 操作日志详情
func (s *Service) SystemLogDetail(ctx context.Context, id int64) (*dm.OperationLog, error) {
	log, err := s.dao.GetLogByID(ctx, id)
	if err != nil {
		xlog.Errorf("SystemLogDetail GetLogByID(%d), err:%v", id, err)
		return nil, ec.ServerErr
	}
	if log == nil {
		return nil, ec.NotFound
	}
	return log, nil
}

// SystemLogExport 操作日志导出
func (s *Service) SystemLogExport(ctx context.Context, req *model.LogListReq) ([]*dm.OperationLog, error) {
	list, err := s.dao.LogListAll(ctx, req)
	if err != nil {
		xlog.Errorf("SystemLogExport LogListAll, err:%v", err)
		return nil, ec.ServerErr
	}
	return list, nil
}
