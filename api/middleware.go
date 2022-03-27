package api

import (
	"errors"
	"fmt"
	"github.com/AndrewLoveMei/simplebank/token"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey        = "authorization"
	authorizationHeaderTypeBearer = "bearer"
	authorizationPayloadKey       = "authorization_payload"
)

//创建权限认证的中间件
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//从http请求中获得请求头
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provider")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		//校验请求头
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		}
		//校验token类型是否是bearer
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationHeaderTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		//验证token，取出有效负载
		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		//将payload存储在上下文中
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
