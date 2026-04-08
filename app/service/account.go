package service

import (
	"context"
	"errors"
	"time"

	"gopay/app/dao"
	"gopay/app/dm"
	"gopay/app/model"
	ec "gopay/errcode"

	"github.com/go-pay/xlog"
	"github.com/go-pay/xtime"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("gopay-platform-secret-key-2026")

// 内置账号，数据库不可用时使用
var builtinAccounts = map[string]*dm.Account{
	"admin": {ID: 1, Uname: "admin", Pwd: "admin", RealName: "超级管理员", Phone: "13800000001", Role: "admin", Status: 1},
}

// getAccount 获取账户，优先数据库，fallback 内置账号
func (s *Service) getAccount(ctx context.Context, uname string) (*dm.Account, error) {
	account, err := s.dao.GetAccountByUname(ctx, uname)
	if err != nil && !errors.Is(err, dao.ErrNoDatabase) {
		return nil, err
	}
	if account != nil {
		return account, nil
	}
	if acc, ok := builtinAccounts[uname]; ok {
		return acc, nil
	}
	return nil, nil
}

// Login 用户登录
func (s *Service) Login(ctx context.Context, req *model.LoginReq) (*model.LoginRsp, error) {
	account, err := s.getAccount(ctx, req.Username)
	if err != nil {
		xlog.Errorf("Login getAccount(%s), err:%v", req.Username, err)
		return nil, ec.ServerErr
	}
	if account == nil {
		return nil, ec.LoginFailed
	}
	if account.Status == 0 {
		return nil, ec.LoginFailed
	}

	// 验证密码：优先 bcrypt，兼容明文
	if errBcrypt := bcrypt.CompareHashAndPassword([]byte(account.Pwd), []byte(req.Password)); errBcrypt != nil {
		if account.Pwd != req.Password {
			return nil, ec.LoginFailed
		}
	}

	// 更新最后登录时间
	now := time.Now()
	_ = s.dao.UpdateAccountLastLogin(ctx, account.ID, now)

	// 生成 JWT token
	claims := jwt.MapClaims{
		"uid":   account.ID,
		"uname": account.Uname,
		"iat":   now.Unix(),
		"exp":   now.Add(24 * time.Hour).Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
	if err != nil {
		xlog.Errorf("Login jwt.SignedString, err:%v", err)
		return nil, ec.ServerErr
	}

	return &model.LoginRsp{
		Token: token,
		UserInfo: &model.UserInfo{
			ID:        account.ID,
			Username:  account.Uname,
			RealName:  account.RealName,
			Phone:     account.Phone,
			Email:     account.Email,
			Role:      account.Role,
			LastLogin: now.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

// GetUserInfo 获取用户信息（从 token 解析）
func (s *Service) GetUserInfo(ctx context.Context, tokenStr string) (*model.UserInfo, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ec.TokenInvalid
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, ec.TokenInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ec.TokenInvalid
	}

	uname, _ := claims["uname"].(string)
	account, err := s.getAccount(ctx, uname)
	if err != nil {
		xlog.Errorf("GetUserInfo getAccount(%s), err:%v", uname, err)
		return nil, ec.ServerErr
	}
	if account == nil {
		return nil, ec.TokenInvalid
	}

	lastLogin := ""
	if account.LastLogin != nil {
		lastLogin = account.LastLogin.Time().Format(xtime.TimeLayout)
	}

	return &model.UserInfo{
		ID:        account.ID,
		Username:  account.Uname,
		RealName:  account.RealName,
		Phone:     account.Phone,
		Email:     account.Email,
		Role:      account.Role,
		LastLogin: lastLogin,
	}, nil
}

// ChangePwd 修改密码
func (s *Service) ChangePwd(ctx context.Context, uid int64, req *model.ChangePwdReq) error {
	if req.NewPassword != req.ConfirmPassword {
		return ec.RequestErr
	}
	account, err := s.dao.GetAccountByID(ctx, uid)
	if err != nil {
		xlog.Errorf("ChangePwd GetAccountByID(%d), err:%v", uid, err)
		return ec.ServerErr
	}
	if account == nil {
		return ec.NotFound
	}
	// 验证旧密码：优先 bcrypt，兼容明文
	if errBcrypt := bcrypt.CompareHashAndPassword([]byte(account.Pwd), []byte(req.OldPassword)); errBcrypt != nil {
		if account.Pwd != req.OldPassword {
			return ec.LoginFailed
		}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		xlog.Errorf("ChangePwd bcrypt.GenerateFromPassword, err:%v", err)
		return ec.ServerErr
	}
	return s.dao.UpdateAccountPwd(ctx, uid, string(hash))
}

// UpdateProfile 更新个人资料
func (s *Service) UpdateProfile(ctx context.Context, uid int64, req *model.ProfileReq) error {
	return s.dao.UpdateAccountProfile(ctx, uid, req.RealName, req.Phone, req.Email)
}

// CreateOperationLog 创建操作日志
func (s *Service) CreateOperationLog(log *dm.OperationLog) error {
	return s.dao.CreateOperationLog(log)
}
