package user

import (
	"ehelp/middleware"
	"ehelp/o/auth"
	"ehelp/o/user"
	"ehelp/x/rest"
	"github.com/gin-gonic/gin"
)

type UserServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewUserServer(parent *gin.RouterGroup, name string) {
	var s = UserServer{
		RouterGroup: parent.Group(name),
	}
	s.GET("/super_admin", s.handleExistSuperAdmin)
	s.POST("/register", s.handleRegister)
	s.POST("/create", s.handleCreate).Use(middleware.MustBeSuperAdmin())
}

func (s *UserServer) handleRegister(ctx *gin.Context) {
	var u *user.Staff
	ctx.BindJSON(&u)
	rest.AssertNil(u.Create())
	auth, err := auth.Create(u.ID, string(u.Role))
	rest.AssertNil(err)
	s.SendData(ctx, map[string]interface{}{
		"access_token": auth.ID,
	})
}
func (s *UserServer) handleCreate(ctx *gin.Context) {
	s.createUser(ctx)
}
func (s *UserServer) createUser(ctx *gin.Context) {
	var u *user.Staff
	ctx.BindJSON(&u)
	rest.AssertNil(u.Create())
	s.SendData(ctx, u)
}
func (s *UserServer) handleExistSuperAdmin(ctx *gin.Context) {
	s.SendData(ctx, user.GetSuperUser())
}
