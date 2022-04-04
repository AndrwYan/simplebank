package api

import (
	"database/sql"
	"errors"
	db "github.com/AndrwYan/simplebank/db/sqlc"

	"github.com/AndrwYan/simplebank/token"
	"github.com/lib/pq"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//创建请求体
type createAccountRequest struct {
	Currency string `json:"currency",binding:"required,oneof=USD EUR"` //前面的json文件代表这是绑定json
}

func (server *Server) createAccount(ctx *gin.Context) {

	var req createAccountRequest

	//将请求中的数据绑定到相应的结构体之中
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	//获取授权
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	//字面量
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok { //因为外键约束而增加的条件
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation": //外键约束和唯一性约束
				ctx.JSON(http.StatusForbidden, errorResponse(err))
			}
			log.Println(pqErr.Code.Name())
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

//通过id查询
type getAccountRequest struct {
	ID int64 `uri:"id",binding:"required,min=1"` //由于是从uri中获取参数所以是uri
}

//根据id查询account
func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest

	//绑定url中的参数
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	//代用
	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		//数据库连接的时候报错
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	//获取授权
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	//将查到的user和授权中的api查看
	if authPayload.Username != account.Owner {
		err := errors.New("account doesn't belong to the authentication user")

		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, account)
}

//通过id查询
type listAccountRequest struct {
	pageID   int32 `form:"page_id",binding:"required,min=1"`
	pageSize int32 `form:"page_size",binding:"required,min=5,max=10"`
}

//分页查询查询account
func (server *Server) listAccount(ctx *gin.Context) {
	var req listAccountRequest
	//绑定query中的参数，所以调用的是ShouldBuildingQuery()函数
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	//给分页设置参数
	arg := db.ListAccountsParams{
		Limit:  req.pageSize,
		Offset: (req.pageID - 1) * req.pageSize,
	}
	//分页
	account, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		//数据库连接的时候报错
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusAccepted, account)
}
