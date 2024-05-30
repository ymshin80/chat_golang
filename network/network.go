package network

import (
	"chat_server_golang/service"
	"log"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

type Server struct{
	engine *gin.Engine

	service *service.Service
	
	port string
	ip string
}

// type Network struct{
// 	engine *gin.Engine
// 	service *service.Service
// 	repository *repository.Repository
	
// 	port string
// 	ip string
// }

func NewServer(service *service.Service,  port string ) *Server {
	s := &Server{engine: gin.New(), service: service, port: port}

	///////////////////////middleware 설정///////////////////////////////
	//default -- log4J
	s.engine.Use(gin.Logger())
	s.engine.Use(gin.Recovery())
	//cross site 설정
	s.engine.Use(cors.New(cors.Config{
		AllowWebSockets: true,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET","POST","PUT", "DELETE", "PATCH"},
		AllowHeaders: []string{"ORIGIN","Content-Type", "Content-Length", "Access-Control-Allow-Headers", "Access-Control-Allow-Origin", "Authorization", "X-Requested-With", "expires"},
		ExposeHeaders: []string{"ORIGIN","Content-Type", "Content-Length", "Access-Control-Allow-Headers", "Access-Control-Allow-Origin", "Authorization", "X-Requested-With", "expires"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
	}))

	registerServer(s)
	//default
	//r := NewRoom()
	//goroutine  백그라운드에서 동작.. 
	//go r.RunInit()

	//n.engine.GET("/room", r.SocketServe)
	return s
}

func (n *Server) StartServer() error {
	log.Print("starting server")

	return n.engine.Run(n.port)
}