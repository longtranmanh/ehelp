package auth

import (
	"ehelp/o/auth"
	"ehelp/o/push_token"
	"ehelp/o/user"
	oAuth "ehelp/o/user/auth"
	"ehelp/x/rest"
	"fmt"
	"github.com/gin-gonic/gin"
)

type AuthServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewAuthServer(parent *gin.RouterGroup, name string) {
	var s = AuthServer{
		RouterGroup: parent.Group(name),
	}
	s.POST("/signin", s.handleSignin)
	s.POST("/delete", s.handleDeleAccount)
}

func (s *AuthServer) handleSignin(ctx *gin.Context) {
	var loginInfo = struct {
		UName    string `json:"uname"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}{}
	ctx.BindJSON(&loginInfo)
	u, err := user.GetByUNamePwd(loginInfo.UName, loginInfo.Password, loginInfo.Role)
	rest.AssertNil(err)
	fmt.Println(err)
	auth, err := auth.Create(u.ID, string(u.Role))
	rest.AssertNil(err)
	s.SendData(ctx, map[string]interface{}{
		"access_token": auth.ID,
	})
}

func (s *AuthServer) handleDeleAccount(ctx *gin.Context) {
	var cus, emp = oAuth.GetUserFromToken(ctx.Request)
	if cus != nil {
		var userID = cus.ID
		rest.AssertNil(push_token.DeleteTokenByUser(userID))
		rest.AssertNil(oAuth.DeleteCustomerID(userID))
	} else {
		var userID = emp.ID
		rest.AssertNil(push_token.DeleteTokenByUser(userID))
		rest.AssertNil(oAuth.DeleteEmpID(userID))
	}
	s.Success(ctx)
}
