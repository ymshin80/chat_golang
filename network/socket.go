package network

import (
	"chat_server_golang/service"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = &websocket.Upgrader{	
	ReadBufferSize: SocketBufferSize,
	WriteBufferSize: MessageBufferSize,
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Room struct {
	Forward chan *message //수신되는 메시지, 클라언트들에게 전송

	Join chan *client //socket이 연결되는 경우 채팅방에 입장한 사용자
	Leave chan *client	//socket이 끊어지는 경우 채팅방을 떠난 사용자

	Clients map[*client]bool //현재 방에 잇는 사용자 저장

	service *service.Service
}

type message struct {
	Name     string `json:"name"`
	Room     string `json:"room"`
	Message  string `json:"message"`
	SendDtm  time.Time  `json:"sendDtm"`
}

type client struct {
	Send   chan *message
	Room   *Room
	Name   string 	`json:"name"`
	Socket *websocket.Conn
}

const (
	SocketBufferSize  = 1024
	MessageBufferSize = 256
)


func NewRoom(service *service.Service) *Room {
	return &Room {
		Forward: make(chan *message),
		Join: make(chan *client),
		Leave: make(chan *client),
		Clients: make(map[*client]bool),
		service: service,
	}
}

func (r *Room) RunInit() {
	//Room 에 있는 모든 채널정보를 받는 역활
	for {
		select {
			case client := <- r.Join:
				r.Clients[client] = true

			case client := <- r.Leave:
				r.Clients[client] = false
				//map에서 빼줌.
				delete(r.Clients, client)
				//떠나니 .. 채널을 닫는다.
				close(client.Send)
			case msg := <- r.Forward:
				go r.service.InsertChatting(msg.Name, msg.Message, msg.Room)

				for client := range r.Clients {
					client.Send <- msg
				}
		}
	}
}

func (c *client) Read() {
	defer c.Socket.Close()
	for {
		var msg *message
		err :=c.Socket.ReadJSON(&msg)
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway){
				break
			} else {
				panic(err)
			}
		} else {
			log.Println("READ::",msg,"client::", c.Name)
			msg.SendDtm = time.Now()
			msg.Name = c.Name

			c.Room.Forward <- msg
		} 
	}
}
func (c *client) Write() {
	defer c.Socket.Close()
	for msg := range c.Send {
		err := c.Socket.WriteJSON(msg)
		log.Println("WRITE::",msg,"client::", c.Name)
		if err != nil {
 			panic(err)
		}
	}
}

func (r *Room) ServerHTTP(c *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	
	socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
 		log.Fatal("----ServerHTTP:::",err)
		return
	}
	
	userCookie, err := c.Request.Cookie("auth")
	if err != nil {
 		log.Fatal("----ServerHTTP:::",err)
		return
	}
	//입장 클라이언트 생성
	client := &client{
		Socket: socket,
		Send: make(chan *message, MessageBufferSize),
		Room: r,
		Name: userCookie.Value,
	}
	//채팅방에 join
	r.Join <- client

	//함수에서 벚어날때.. 실행 됨.
	//?? 클라이언트 화면에서 메신저를 이탈하고나 방에서 나가는 경우가 아닌데도. 
	//단순히 입장을 했으니 나중에 나가는 경우를 위해 나가는 사용자에도 등록만 해놓는 형태인지..
	defer func() { r.Leave <- client } ()

	log.Println("여기가 마지막...")


	go client.Write()



	client.Read()
}
//사용자가 메시전에 접속해서 방에 입장을 했을 때..
// func (r *Room) SocketServe(c *gin.Context) {
// 	socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
//  		panic(err)
// 	}
	
// 	userCookie, err := c.Request.Cookie("auth")
// 	if err != nil {
//  		panic(err)
// 	}
// 	//입장 클라이언트 생성
// 	client := &client{
// 		Socket: socket,
// 		Send: make(chan *message, types.MessageBufferSize),
// 		Room: r,
// 		Name: userCookie.Value,
// 	}
// 	//채팅방에 join
// 	r.Join <- client

// 	//함수에서 벚어날때.. 실행 됨.
// 	//?? 클라이언트 화면에서 메신저를 이탈하고나 방에서 나가는 경우가 아닌데도. 
// 	//단순히 입장을 했으니 나중에 나가는 경우를 위해 나가는 사용자에도 등록만 해놓는 형태인지..
// 	defer func() { r.Leave <- client } ()

// 	log.Println("여기가 마지막...")


// 	go client.Write()



// 	client.Read()
// }