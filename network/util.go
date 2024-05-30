package network

import (
	"chat_server_golang/types"

	"github.com/gin-gonic/gin"
)

func response(c *gin.Context, s int, res interface{}, data ...string) {
	c.JSON(s, types.NewRes(s, res, data...))
}