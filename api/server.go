package api

import (
	"fmt"
	db "github.com/AndrewLoveMei/simplebank/db/sqlc"
	"github.com/AndrewLoveMei/simplebank/db/util"
	"github.com/AndrewLoveMei/simplebank/token"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      *db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer(config util.Config, store *db.Store) (*Server, error) {

	//初始化maker，并组合进server
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can not create token maker: %w", err)
	}
	server := &Server{
		tokenMaker: tokenMaker,
		store:      store,
		config:     config,
	}

	//向gin注册我们自己编写的验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency) //第一个参数是需要验证的参数,第二个参数是自定义的校验器
	}

	server.setupRouter()
	return server, nil
}

//注册api路由的函数
func (server *Server) setupRouter() {
	router := gin.Default()

	//Context很重要，我们在处理程序的时候，所做的一切都会影响到这个上下文。
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.POST("/transfer", server.createTransfer)
	router.POST("/users", server.createUser)
	router.POST("users/login", server.loginUser)

	//add router to router
	server.router = router
}

// Start 此函数的作用就是在输入地址上运行http服务器开始监听API请求。
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"key": err.Error()}
}
