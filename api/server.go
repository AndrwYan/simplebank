package api

import (
	db "github.com/AndrewLoveMei/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	//Context很重要，我们在处理程序的时候，所做的一切都会影响到这个上下文。
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.POST("/transfer", server.createTransfer)

	//add router to router
	server.router = router
	return server
}

// Start 此函数的作用就是在输入地址上运行http服务器开始监听API请求。
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"key": err.Error()}
}
