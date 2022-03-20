package main

import (
	"database/sql"
	"github.com/AndrewLoveMei/simplebank/api"
	db "github.com/AndrewLoveMei/simplebank/db/sqlc"
	"github.com/AndrewLoveMei/simplebank/db/util"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	//加载配置文件
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("无法加载配置文件！")
	}
	//建立数据库连接
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to the db!")
	}

	//建立服务
	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	//致命日志
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
