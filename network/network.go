package network

import (
	"chat_server_golang/service"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

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

	//시작하면 control server에 알려줌. 



	//default
	//r := NewRoom()
	//goroutine  백그라운드에서 동작.. 
	//go r.RunInit()

	//n.engine.GET("/room", r.SocketServe)
	return s
}
func (s *Server) setServerInfo() {
	//ip, 

	if addrs, err := net.InterfaceAddrs(); err != nil {
		panic(err)
	} else {
		var ip net.IP

		for _, addr := range addrs{
			if ipnet, ok := addr.(*net.IPNet);  ok {
				if 	!ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
					ip = ipnet.IP
					break
				}

			}
		}
		if ip == nil {
			panic("no ip address found")
		} else {
			if err = s.service.ServerSet(ip.String()+s.port, true); err != nil {
				panic(err)
			} else {
				s.ip = ip.String()
			}
			//서버 시작 이벤트 발생
			s.service.PublishServerStatusEvent(ip.String()+s.port, true)
		}
	}


}
func (n *Server) StartServer() error {
	n.setServerInfo()
	
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGINT)

	//서버가 down 될때 이벤트 발생
	go func() {
		<- channel
		//이벤트가 발생되면 input
		
		if err := n.service.ServerSet(n.ip+n.port, false); err != nil {
			log.Println("Failed To Set ServerInfo When Close", "err", err)
		} 

		n.service.PublishServerStatusEvent(n.ip+n.port, false)
		os.Exit(1)

	}()


	return n.engine.Run(n.port)
}