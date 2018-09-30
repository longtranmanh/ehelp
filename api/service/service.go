package service

import (
	"ehelp/middleware"
	"ehelp/o/service"
	"ehelp/o/tool"
	"ehelp/x/rest"
	"g/x/web"
	"github.com/gin-gonic/gin"
)

type ServiceServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewServiceServer(parent *gin.RouterGroup, name string) {
	var s = &ServiceServer{
		RouterGroup: parent.Group(name),
	}
	s.Use(middleware.MustBeSuperAdmin())
	s.POST("/create", s.handleCreate)
	s.GET("/list", s.handleList)
	s.GET("/list_tool", s.handleListServiceTool)
	s.GET("/delete/:id", s.handleDelete)
	s.POST("/tool/create", s.handleCreateTool)
	s.GET("/tool/list", s.handleListTool)

}

func (s *ServiceServer) handleListServiceTool(ctx *gin.Context) {
	services, err := service.GetToolServices("")
	rest.AssertNil(err)
	s.SendData(ctx, services)
}

func (s *ServiceServer) handleCreate(ctx *gin.Context) {
	var srv *service.Service
	web.AssertNil(ctx.ShouldBindJSON(&srv))
	web.AssertNil(srv.Create())
	s.SendData(ctx, srv)
}
func (s *ServiceServer) handleList(ctx *gin.Context) {
	services, err := service.GetServices()
	rest.AssertNil(err)
	s.SendData(ctx, services)
}

func (s *ServiceServer) handleDelete(ctx *gin.Context) {
	// var id = ctx.Param("id")
	// web.AssertNil(service.DeleteByID(id))
	// s.Success(ctx)
}

func (s *ServiceServer) handleCreateTool(ctx *gin.Context) {
	var tool *tool.Tool
	web.AssertNil(ctx.ShouldBindJSON(&tool))
	web.AssertNil(tool.Create())
	s.SendData(ctx, tool)
}
func (s *ServiceServer) handleListTool(ctx *gin.Context) {
	tools, err := tool.GetTools()
	rest.AssertNil(err)
	s.SendData(ctx, tools)
}
