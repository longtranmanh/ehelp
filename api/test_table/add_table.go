package test_table

import (
	"ehelp/o/service"
	"ehelp/x/rest"
	"github.com/gin-gonic/gin"
)

type TableServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewTableServer(parent *gin.RouterGroup, name string) {
	var s = TableServer{
		RouterGroup: parent.Group(name),
	}
	s.POST("/add_service", s.handlerAddService)
}

func (s *TableServer) handlerAddService(ctx *gin.Context) {
	var body = service.Service{}
	ctx.BindJSON(&body)
	body.Create()
	s.SendData(ctx, map[string]interface{}{})
}
