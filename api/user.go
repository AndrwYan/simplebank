package api

import (
	"database/sql"
	db "github.com/AndrewLoveMei/simplebank/db/sqlc"
	"github.com/AndrewLoveMei/simplebank/db/util"
	"github.com/lib/pq"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//创建请求体
type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}
type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

//之所以创建函数防止user结构体中的password的密码被泄露
func newUserResponse(user db.User) userResponse {
	userResp := userResponse{
		Email:             user.Email,
		Username:          user.Username,
		FullName:          user.FullName,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
	return userResp
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	//将http请求中的数据绑定到相应的结构体之中
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	password, err := util.HashPassword(req.password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	//创建user字面量值
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: password,
		Email:          req.Email,
		FullName:       req.FullName,
	}
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok { //因为外键约束而增加的条件
			switch pqErr.Code.Name() {
			case "unique_violation": //外键约束和唯一性约束
				ctx.JSON(http.StatusForbidden, errorResponse(err))
			}
			log.Println(pqErr.Code.Name())
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	//调用该函数过滤密码
	rsp := newUserResponse(user)
	//重新封装结构体返回response
	ctx.JSON(http.StatusOK, rsp)
}

//封装请求的结构体
type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	password string `json:"password" binding:"required,min=6"`
}

//封装返回的结构体
type loginUserResponse struct {
	AccessToken string       `json:"username" binding:"required,alphanum"`
	User        userResponse `json:"user"`
}

//编写用于登录的代码
func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	//去数据库中查找username对应的user
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	err = util.CheckedPassword(req.password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}

	accessToken, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}
