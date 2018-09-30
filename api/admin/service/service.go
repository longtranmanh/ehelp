package service

import (
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
	s.POST("/create", s.handleCreate)
	s.POST("/update", s.handleUpdate)
	s.GET("/list", s.handleList)
	s.GET("/list-tool", s.handleListServiceTool)
	s.GET("/delete/:id", s.handleDelete)
	s.POST("/tool/create", s.handleCreateTool)
	s.POST("/tool/update", s.handleUpdateTool)
	s.GET("/tool/list", s.handleListTool)
	s.GET("/tool/delete", s.handleDeleteTool)
	s.GET("/delete", s.handleDeleteService)

}

func (s *ServiceServer) handleUpdate(ctx *gin.Context) {
	var srv *service.Service
	web.AssertNil(ctx.ShouldBindJSON(&srv))
	web.AssertNil(srv.Update())
	s.SendData(ctx, srv)
}

func (s *ServiceServer) handleUpdateTool(ctx *gin.Context) {
	var t *tool.Tool
	web.AssertNil(ctx.ShouldBindJSON(&t))
	web.AssertNil(t.Update())
	s.Success(ctx)
}

func (s *ServiceServer) handleDeleteService(ctx *gin.Context) {
	var id = ctx.Query("id")
	rest.AssertNil(service.DeleteServiceByID(id))
	s.Success(ctx)
}

func (s *ServiceServer) handleDeleteTool(ctx *gin.Context) {
	var id = ctx.Query("id")
	rest.AssertNil(tool.DeleteToolByID(id))
	s.Success(ctx)
}
func (s *ServiceServer) handleListServiceTool(ctx *gin.Context) {
	var types = ctx.Query("type")
	services, err := service.GetToolServices(types)
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
	var id = ctx.Param("id")
	web.AssertNil(service.DeleteServiceByID(id))
	s.Success(ctx)
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
