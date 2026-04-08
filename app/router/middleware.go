package router

import (
	"strings"
	"time"

	"gopay/app/dm"
	ec "gopay/errcode"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/web"
	"github.com/go-pay/xlog"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("gopay-platform-secret-key-2026")

const (
	CtxKeyUsername = "ctx_username"
	CtxKeyUID      = "ctx_uid"
)

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			web.JSON(c, nil, ec.TokenInvalid)
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			web.JSON(c, nil, ec.TokenInvalid)
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ec.TokenInvalid
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			web.JSON(c, nil, ec.TokenInvalid)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			web.JSON(c, nil, ec.TokenInvalid)
			c.Abort()
			return
		}

		uname, _ := claims["uname"].(string)
		uid, _ := claims["uid"].(float64)
		c.Set(CtxKeyUsername, uname)
		c.Set(CtxKeyUID, int64(uid))
		c.Next()
	}
}

// OperationLogger 操作日志中间件 (异步写入)
func OperationLogger(module, action, description string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Milliseconds()

		username, _ := c.Get(CtxKeyUsername)
		uname, _ := username.(string)
		if uname == "" {
			uname = "anonymous"
		}

		success := int8(1)
		if c.Writer.Status() >= 400 {
			success = 0
		}

		log := &dm.OperationLog{
			Operator:    uname,
			Module:      module,
			Action:      action,
			Description: description,
			IP:          c.ClientIP(),
			UserAgent:   c.GetHeader("User-Agent"),
			Success:     success,
			Duration:    int(duration),
		}

		go func() {
			if svc != nil {
				if err := svc.CreateOperationLog(log); err != nil {
					xlog.Errorf("OperationLogger CreateOperationLog err:%v", err)
				}
			}
		}()
	}
}
