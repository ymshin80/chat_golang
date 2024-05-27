package network

import (
	"log"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

type Network struct{
	engine *gin.Engine
}

func NewServer() *Network {
	n := &Network{engine: gin.New()}

	///////////////////////middleware 설정///////////////////////////////
	//default -- log4J
	n.engine.Use(gin.Logger())
	n.engine.Use(gin.Recovery())
	//cross site 설정
	n.engine.Use(cors.New(cors.Config{
		AllowWebSockets: true,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET","POST","PUT"},
		AllowHeaders: []string{"*"},
		AllowCredentials: true,
	}))
	//default
	r := NewRoom()
	//goroutine  백그라운드에서 동작.. 
	go r.RunInit()

	n.engine.GET("/room", r.SocketServe)

	return n
}

func (n *Network) StartServer() error {
	log.Print("starting server")

	return n.engine.Run(":8080")
}