package router

import (
	"gopay/app/model"
	"gopay/errcode"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/web"
	"github.com/go-pay/xlog"
)

// ---- 用户管理 ----

func systemUserList(c *gin.Context) {
	req := new(model.UserListReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("systemUserList ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	rsp, err := svc.SystemUserList(c, req)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, rsp, nil)
}

func systemUserAdd(c *gin.Context) {
	req := new(model.UserAddReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("systemUserAdd ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	if err := svc.SystemUserAdd(c, req); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}

func systemUserUpdate(c *gin.Context) {
	req := new(model.UserUpdateReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("systemUserUpdate ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	if err := svc.SystemUserUpdate(c, req); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}

func systemUserToggleStatus(c *gin.Context) {
	req := new(model.IDReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("systemUserToggleStatus ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	if err := svc.SystemUserToggleStatus(c, req.ID); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}

func systemUserResetPwd(c *gin.Context) {
	req := new(model.IDReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("systemUserResetPwd ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	newPwd, err := svc.SystemUserResetPwd(c, req.ID)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, gin.H{"password": newPwd}, nil)
}

// ---- 角色管理 ----

func systemRoleList(c *gin.Context) {
	req := new(model.RoleListReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("systemRoleList ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	rsp, err := svc.SystemRoleList(c, req)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, rsp, nil)
}

func systemRoleAdd(c *gin.Context) {
	req := new(model.RoleAddReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("systemRoleAdd ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	if err := svc.SystemRoleAdd(c, req); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}

func systemRoleUpdate(c *gin.Context) {
	req := new(model.RoleUpdateReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("systemRoleUpdate ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	if err := svc.SystemRoleUpdate(c, req); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}

func systemRoleToggleStatus(c *gin.Context) {
	req := new(model.IDReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("systemRoleToggleStatus ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	if err := svc.SystemRoleToggleStatus(c, req.ID); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}

func systemRolePermsUpdate(c *gin.Context) {
	req := new(model.RolePermsUpdateReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("systemRolePermsUpdate ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	if err := svc.SystemRolePermsUpdate(c, req); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}

func systemRolePermsList(c *gin.Context) {
	req := new(model.RolePermsListReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("systemRolePermsList ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	perms, err := svc.SystemRolePermsList(c, req.RoleID)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, perms, nil)
}

// ---- 操作日志 ----

func systemLogList(c *gin.Context) {
	req := new(model.LogListReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("systemLogList ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	rsp, err := svc.SystemLogList(c, req)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, rsp, nil)
}

func systemLogDetail(c *gin.Context) {
	req := new(model.IDReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("systemLogDetail ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	log, err := svc.SystemLogDetail(c, req.ID)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, log, nil)
}

func systemLogExport(c *gin.Context) {
	req := new(model.LogListReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("systemLogExport ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	list, err := svc.SystemLogExport(c, req)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, list, nil)
}
