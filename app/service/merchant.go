package service

import (
	"context"
	"fmt"
	"time"

	"gopay/app/dm"
	"gopay/app/model"
	ec "gopay/errcode"

	"github.com/go-pay/xlog"
)

// ---- 商户管理 ----

// MerchantList 商户列表
func (s *Service) MerchantList(ctx context.Context, req *model.MerchantListReq) (*model.PageResp, error) {
	list, total, err := s.dao.MerchantList(ctx, req)
	if err != nil {
		xlog.Errorf("MerchantList, err:%v", err)
		return nil, ec.ServerErr
	}
	return &model.PageResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// MerchantAdd 新增商户
func (s *Service) MerchantAdd(ctx context.Context, req *model.MerchantAddReq) (int64, error) {
	m := &dm.Merchant{
		Name:    req.Name,
		Contact: req.Contact,
		Phone:   req.Phone,
		Email:   req.Email,
		Remark:  req.Remark,
		Status:  1,
	}
	if err := s.dao.CreateMerchant(ctx, m); err != nil {
		xlog.Errorf("MerchantAdd CreateMerchant, err:%v", err)
		return 0, ec.ServerErr
	}
	return m.ID, nil
}

// MerchantUpdate 编辑商户
func (s *Service) MerchantUpdate(ctx context.Context, req *model.MerchantUpdateReq) error {
	m, err := s.dao.GetMerchantByID(ctx, req.ID)
	if err != nil {
		xlog.Errorf("MerchantUpdate GetMerchantByID(%d), err:%v", req.ID, err)
		return ec.ServerErr
	}
	if m == nil {
		return ec.NotFound
	}
	updates := map[string]interface{}{
		"name":    req.Name,
		"contact": req.Contact,
		"phone":   req.Phone,
		"email":   req.Email,
		"remark":  req.Remark,
	}
	if err = s.dao.UpdateMerchant(ctx, req.ID, updates); err != nil {
		xlog.Errorf("MerchantUpdate UpdateMerchant(%d), err:%v", req.ID, err)
		return ec.ServerErr
	}
	return nil
}

// MerchantToggleStatus 切换商户状态
func (s *Service) MerchantToggleStatus(ctx context.Context, id int64) error {
	m, err := s.dao.GetMerchantByID(ctx, id)
	if err != nil {
		xlog.Errorf("MerchantToggleStatus GetMerchantByID(%d), err:%v", id, err)
		return ec.ServerErr
	}
	if m == nil {
		return ec.NotFound
	}
	if err = s.dao.ToggleMerchantStatus(ctx, id); err != nil {
		xlog.Errorf("MerchantToggleStatus ToggleMerchantStatus(%d), err:%v", id, err)
		return ec.ServerErr
	}
	return nil
}

// MerchantOptions 商户下拉选项
func (s *Service) MerchantOptions(ctx context.Context) (interface{}, error) {
	list, err := s.dao.MerchantOptions(ctx)
	if err != nil {
		xlog.Errorf("MerchantOptions, err:%v", err)
		return nil, ec.ServerErr
	}
	type option struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	opts := make([]*option, 0, len(list))
	for _, m := range list {
		opts = append(opts, &option{ID: m.ID, Name: m.Name})
	}
	return opts, nil
}

// ---- 商户应用 ----

// MerchantAppList 商户应用列表
func (s *Service) MerchantAppList(ctx context.Context, req *model.MerchantAppListReq) (*model.PageResp, error) {
	list, total, err := s.dao.MerchantAppList(ctx, req)
	if err != nil {
		xlog.Errorf("MerchantAppList, err:%v", err)
		return nil, ec.ServerErr
	}
	// 填充商户名称
	if len(list) > 0 {
		ids := make([]int64, 0, len(list))
		for _, app := range list {
			ids = append(ids, app.MerchantID)
		}
		nameMap, _ := s.dao.GetMerchantsByIDs(ctx, ids)
		for _, app := range list {
			app.MerchantName = nameMap[app.MerchantID]
		}
	}
	return &model.PageResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// MerchantAppAdd 新增商户应用
func (s *Service) MerchantAppAdd(ctx context.Context, req *model.MerchantAppAddReq) error {
	// 验证商户存在
	m, err := s.dao.GetMerchantByID(ctx, req.MerchantID)
	if err != nil {
		xlog.Errorf("MerchantAppAdd GetMerchantByID(%d), err:%v", req.MerchantID, err)
		return ec.ServerErr
	}
	if m == nil {
		return ec.NotFound
	}
	app := &dm.MerchantApp{
		Name:         req.Name,
		Appid:        req.Appid,
		MerchantID:   req.MerchantID,
		PlatformType: req.PlatformType,
		MerchantType: req.MerchantType,
		NotifyUrl:    req.NotifyUrl,
		ReturnUrl:    req.ReturnUrl,
		Status:       1,
	}
	if err = s.dao.CreateMerchantApp(ctx, app); err != nil {
		xlog.Errorf("MerchantAppAdd CreateMerchantApp, err:%v", err)
		return ec.ServerErr
	}
	return nil
}

// MerchantAppUpdate 编辑商户应用
func (s *Service) MerchantAppUpdate(ctx context.Context, req *model.MerchantAppUpdateReq) error {
	app, err := s.dao.GetMerchantAppByID(ctx, req.ID)
	if err != nil {
		xlog.Errorf("MerchantAppUpdate GetMerchantAppByID(%d), err:%v", req.ID, err)
		return ec.ServerErr
	}
	if app == nil {
		return ec.NotFound
	}
	updates := map[string]interface{}{
		"name":          req.Name,
		"merchant_id":   req.MerchantID,
		"platform_type": req.PlatformType,
		"merchant_type": req.MerchantType,
		"notify_url":    req.NotifyUrl,
		"return_url":    req.ReturnUrl,
	}
	if err = s.dao.UpdateMerchantApp(ctx, req.ID, updates); err != nil {
		xlog.Errorf("MerchantAppUpdate UpdateMerchantApp(%d), err:%v", req.ID, err)
		return ec.ServerErr
	}
	return nil
}

// ---- 进件管理 ----

// IncomingApplyList 进件申请列表
func (s *Service) IncomingApplyList(ctx context.Context, req *model.IncomingApplyListReq) (*model.PageResp, error) {
	list, total, err := s.dao.IncomingApplyList(ctx, req)
	if err != nil {
		xlog.Errorf("IncomingApplyList, err:%v", err)
		return nil, ec.ServerErr
	}
	s.fillIncomingMerchantNames(ctx, list)
	return &model.PageResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// IncomingApplyAdd 新建进件
func (s *Service) IncomingApplyAdd(ctx context.Context, req *model.IncomingApplyAddReq) error {
	// 验证商户存在
	m, err := s.dao.GetMerchantByID(ctx, req.MerchantID)
	if err != nil {
		xlog.Errorf("IncomingApplyAdd GetMerchantByID(%d), err:%v", req.MerchantID, err)
		return ec.ServerErr
	}
	if m == nil {
		return ec.NotFound
	}
	// 生成申请编号
	count, _ := s.dao.CountTodayApplies(ctx)
	applyNo := fmt.Sprintf("INC%s%03d", time.Now().Format("20060102"), count+1)

	status := int8(0) // 草稿
	if req.Submit {
		status = 1 // 审核中
	}
	apply := &dm.IncomingApply{
		ApplyNo:     applyNo,
		MerchantID:  req.MerchantID,
		ChannelType: req.ChannelType,
		MerchantNo:  req.MerchantNo,
		LicenseNo:   req.LicenseNo,
		LicenseImg:  req.LicenseImg,
		LegalPerson: req.LegalPerson,
		IdCardFront: req.IdCardFront,
		IdCardBack:  req.IdCardBack,
		Phone:       req.Phone,
		Remark:      req.Remark,
		Status:      status,
	}
	if err = s.dao.CreateIncomingApply(ctx, apply); err != nil {
		xlog.Errorf("IncomingApplyAdd CreateIncomingApply, err:%v", err)
		return ec.ServerErr
	}
	return nil
}

// IncomingApplySubmit 提交审核 (status 0→1)
func (s *Service) IncomingApplySubmit(ctx context.Context, id int64) error {
	apply, err := s.dao.GetIncomingApplyByID(ctx, id)
	if err != nil {
		xlog.Errorf("IncomingApplySubmit GetIncomingApplyByID(%d), err:%v", id, err)
		return ec.ServerErr
	}
	if apply == nil {
		return ec.NotFound
	}
	if apply.Status != 0 {
		return ec.RequestErr
	}
	return s.dao.UpdateIncomingApplyStatus(ctx, id, 1)
}

// IncomingApplyReview 审核进件
func (s *Service) IncomingApplyReview(ctx context.Context, req *model.IncomingApplyReviewReq, reviewer string) error {
	apply, err := s.dao.GetIncomingApplyByID(ctx, req.ID)
	if err != nil {
		xlog.Errorf("IncomingApplyReview GetIncomingApplyByID(%d), err:%v", req.ID, err)
		return ec.ServerErr
	}
	if apply == nil {
		return ec.NotFound
	}
	if apply.Status != 1 {
		return ec.RequestErr
	}
	var status int8
	switch req.Action {
	case "pass":
		status = 2
	case "reject":
		status = 3
	default:
		return ec.RequestErr
	}
	return s.dao.UpdateIncomingApplyReview(ctx, req.ID, status, reviewer, req.Remark)
}

// IncomingRecordList 进件记录列表
func (s *Service) IncomingRecordList(ctx context.Context, req *model.IncomingRecordListReq) (*model.PageResp, error) {
	list, total, err := s.dao.IncomingRecordList(ctx, req)
	if err != nil {
		xlog.Errorf("IncomingRecordList, err:%v", err)
		return nil, ec.ServerErr
	}
	s.fillIncomingMerchantNames(ctx, list)
	return &model.PageResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// IncomingRecordDetail 进件记录详情
func (s *Service) IncomingRecordDetail(ctx context.Context, id int64) (*dm.IncomingApply, error) {
	apply, err := s.dao.GetIncomingApplyByID(ctx, id)
	if err != nil {
		xlog.Errorf("IncomingRecordDetail GetIncomingApplyByID(%d), err:%v", id, err)
		return nil, ec.ServerErr
	}
	if apply == nil {
		return nil, ec.NotFound
	}
	// 填充商户名称
	nameMap, _ := s.dao.GetMerchantsByIDs(ctx, []int64{apply.MerchantID})
	apply.MerchantName = nameMap[apply.MerchantID]
	return apply, nil
}

// fillIncomingMerchantNames 批量填充进件列表的商户名称
func (s *Service) fillIncomingMerchantNames(ctx context.Context, list []*dm.IncomingApply) {
	if len(list) == 0 {
		return
	}
	ids := make([]int64, 0, len(list))
	for _, a := range list {
		ids = append(ids, a.MerchantID)
	}
	nameMap, _ := s.dao.GetMerchantsByIDs(ctx, ids)
	for _, a := range list {
		a.MerchantName = nameMap[a.MerchantID]
	}
}
