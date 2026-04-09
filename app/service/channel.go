package service

import (
	"context"
	"strings"

	"gopay/app/dm"
	"gopay/app/model"
	ec "gopay/errcode"

	"github.com/go-pay/xlog"
)

// channelToResp 将 dm 转为响应 DTO（处理 payMethods 逗号分隔→数组）
func channelToResp(ch *dm.PaymentChannel) *model.ChannelResp {
	var methods []string
	if ch.PayMethods != "" {
		methods = strings.Split(ch.PayMethods, ",")
	}
	return &model.ChannelResp{
		ID:         ch.ID,
		Name:       ch.Name,
		Code:       ch.Code,
		Type:       ch.Type,
		MerchantID: ch.MerchantID,
		PayMethods: methods,
		FeeRate:    ch.FeeRate,
		Status:     ch.Status,
		Remark:     ch.Remark,
		Ctime:      ch.Ctime,
	}
}

// maskSensitive 脱敏处理
func maskSensitive(s string) string {
	if s == "" {
		return ""
	}
	return "******"
}

// ChannelList 支付通道列表
func (s *Service) ChannelList(ctx context.Context, req *model.ChannelListReq) (*model.PageResp, error) {
	list, total, err := s.dao.ChannelList(ctx, req)
	if err != nil {
		xlog.Errorf("ChannelList, err:%v", err)
		return nil, ec.ServerErr
	}
	// 转换响应 + 填充商户名称
	respList := make([]*model.ChannelResp, 0, len(list))
	if len(list) > 0 {
		ids := make([]int64, 0, len(list))
		for _, ch := range list {
			ids = append(ids, ch.MerchantID)
		}
		nameMap, _ := s.dao.GetMerchantsByIDs(ctx, ids)
		for _, ch := range list {
			r := channelToResp(ch)
			r.MerchantName = nameMap[ch.MerchantID]
			respList = append(respList, r)
		}
	}
	return &model.PageResp{
		List:     respList,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// ChannelAdd 新增通道
func (s *Service) ChannelAdd(ctx context.Context, req *model.ChannelAddReq) error {
	// code 唯一性检查
	existing, err := s.dao.GetChannelByCode(ctx, req.Code)
	if err != nil {
		xlog.Errorf("ChannelAdd GetChannelByCode(%s), err:%v", req.Code, err)
		return ec.ServerErr
	}
	if existing != nil {
		return ec.Conflict
	}
	ch := &dm.PaymentChannel{
		Name:       req.Name,
		Code:       req.Code,
		Type:       req.Type,
		MerchantID: req.MerchantID,
		PayMethods: strings.Join(req.PayMethods, ","),
		FeeRate:    req.FeeRate,
		Remark:     req.Remark,
		Status:     1,
	}
	if err = s.dao.CreateChannel(ctx, ch); err != nil {
		xlog.Errorf("ChannelAdd CreateChannel, err:%v", err)
		return ec.ServerErr
	}
	return nil
}

// ChannelUpdate 编辑通道
func (s *Service) ChannelUpdate(ctx context.Context, req *model.ChannelUpdateReq) error {
	ch, err := s.dao.GetChannelByID(ctx, req.ID)
	if err != nil {
		xlog.Errorf("ChannelUpdate GetChannelByID(%d), err:%v", req.ID, err)
		return ec.ServerErr
	}
	if ch == nil {
		return ec.NotFound
	}
	updates := map[string]interface{}{
		"name":        req.Name,
		"type":        req.Type,
		"merchant_id": req.MerchantID,
		"pay_methods": strings.Join(req.PayMethods, ","),
		"fee_rate":    req.FeeRate,
		"remark":      req.Remark,
	}
	if err = s.dao.UpdateChannel(ctx, req.ID, updates); err != nil {
		xlog.Errorf("ChannelUpdate UpdateChannel(%d), err:%v", req.ID, err)
		return ec.ServerErr
	}
	return nil
}

// ChannelToggleStatus 切换通道状态
func (s *Service) ChannelToggleStatus(ctx context.Context, id int64) error {
	ch, err := s.dao.GetChannelByID(ctx, id)
	if err != nil {
		xlog.Errorf("ChannelToggleStatus GetChannelByID(%d), err:%v", id, err)
		return ec.ServerErr
	}
	if ch == nil {
		return ec.NotFound
	}
	if err = s.dao.ToggleChannelStatus(ctx, id); err != nil {
		xlog.Errorf("ChannelToggleStatus, err:%v", err)
		return ec.ServerErr
	}
	return nil
}

// ChannelDetail 通道详情（含配置，敏感字段脱敏）
func (s *Service) ChannelDetail(ctx context.Context, id int64) (*model.ChannelDetailResp, error) {
	ch, err := s.dao.GetChannelByID(ctx, id)
	if err != nil {
		xlog.Errorf("ChannelDetail GetChannelByID(%d), err:%v", id, err)
		return nil, ec.ServerErr
	}
	if ch == nil {
		return nil, ec.NotFound
	}
	resp := channelToResp(ch)
	// 填充商户名称
	nameMap, _ := s.dao.GetMerchantsByIDs(ctx, []int64{ch.MerchantID})
	resp.MerchantName = nameMap[ch.MerchantID]

	detail := &model.ChannelDetailResp{ChannelResp: resp}
	// 获取配置
	cfg, err := s.dao.GetChannelConfig(ctx, id)
	if err != nil {
		xlog.Errorf("ChannelDetail GetChannelConfig(%d), err:%v", id, err)
	}
	if cfg != nil {
		detail.Config = &model.ChannelConfigResp{
			AppID:      cfg.AppID,
			MchID:      cfg.MchID,
			PrivateKey: maskSensitive(cfg.PrivateKey),
			PublicKey:  maskSensitive(cfg.PublicKey),
			ApiKey:     maskSensitive(cfg.ApiKey),
			SerialNo:   cfg.SerialNo,
			NotifyUrl:  cfg.NotifyUrl,
			SignType:   cfg.SignType,
			Sandbox:    cfg.Sandbox,
		}
	}
	return detail, nil
}

// ChannelConfig 保存通道参数配置
func (s *Service) ChannelConfig(ctx context.Context, req *model.ChannelConfigReq) error {
	ch, err := s.dao.GetChannelByID(ctx, req.ChannelID)
	if err != nil {
		xlog.Errorf("ChannelConfig GetChannelByID(%d), err:%v", req.ChannelID, err)
		return ec.ServerErr
	}
	if ch == nil {
		return ec.NotFound
	}
	sandbox := int8(0)
	if req.Sandbox {
		sandbox = 1
	}
	cfg := &dm.PaymentChannelConfig{
		ChannelID:  req.ChannelID,
		AppID:      req.AppID,
		MchID:      req.MchID,
		PrivateKey: req.PrivateKey,
		PublicKey:  req.PublicKey,
		ApiKey:     req.ApiKey,
		SerialNo:   req.SerialNo,
		NotifyUrl:  req.NotifyUrl,
		SignType:   req.SignType,
		Sandbox:    sandbox,
	}
	if err = s.dao.SaveChannelConfig(ctx, cfg); err != nil {
		xlog.Errorf("ChannelConfig SaveChannelConfig, err:%v", err)
		return ec.ServerErr
	}
	return nil
}
