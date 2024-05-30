package network

import (
	"chat_server_golang/types"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type api struct {
	server *Server
}

func registerServer(server *Server) {
	a := &api{server: server}

	server.engine.GET("/room-list", a.roomList)
	
	server.engine.GET("/room", a.room)
	server.engine.GET("/enter-room", a.enterRoom)

	server.engine.POST("/room", a.makeRoom)

	 r := NewRoom(server.service)
	 go r.RunInit()

	server.engine.GET("/room-chat", r.ServerHTTP)

}


func (a *api) roomList(c *gin.Context){
	if res, err := a.server.service.RoomList(); err != nil {
		response(c, http.StatusInternalServerError, err.Error())
	} else {
		response(c, http.StatusOK, res)
	}
}

func (a *api) makeRoom(c *gin.Context){
	//데이터 가져오기
	var req types.BodyRoomReq
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		response(c, http.StatusUnprocessableEntity, err.Error())
	} else if err := a.server.service.MakeRoom(req.Name); err != nil {
		response(c, http.StatusInternalServerError, err.Error())
	} else {
		response(c, http.StatusOK, "Success Job")
	}
}


func (a *api) room(c *gin.Context){
	var req types.FormRoomReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response(c, http.StatusUnprocessableEntity, err.Error())
	} else {
		res, err := a.server.service.Room(req.Name)
		
		if  err != nil && !strings.Contains(err.Error(), "sql: no rows in result set") {
			response(c, http.StatusInternalServerError, err.Error())
		} else {
			response(c, http.StatusOK, res)
		}
	}
}

func (a *api) enterRoom(c *gin.Context) {
	var req types.FormRoomReq
	
	if err := c.ShouldBindQuery(&req); err != nil {
		response(c, http.StatusUnprocessableEntity, err.Error())
	} else if res, err := a.server.service.EnterRoom(req.Name); err != nil {
		response(c, http.StatusInternalServerError, err.Error())
	} else {
		response(c, http.StatusOK, res)
	}
}



