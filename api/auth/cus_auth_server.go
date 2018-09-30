package auth

import (
	oAuth "ehelp/o/user/auth"
	"ehelp/x/rest"
	"github.com/gin-gonic/gin"
)

type AuthServerMux struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewAuthCustomerServer(parent *gin.RouterGroup, name string) {
	var s = AuthServerMux{
		RouterGroup: parent.Group(name),
	}
	s.POST("/login", s.handleLoginCus)
	s.POST("/logout", s.handleLogoutCus)
	s.POST("/login_facebook", s.handleLoginFacebookCus)
	s.POST("/loginfb_update", s.handleLoginFbUpdateCus)
	s.POST("/login_gmail", s.handleLoginGmailCus)
	s.POST("/logingm_update", s.handleLoginGmUpdateCus)
	s.POST("/register", s.handleRegisterCus)
}

func (s *AuthServerMux) handleLoginCus(ctx *gin.Context) {
	var body = oAuth.LoginUser{}
	ctx.BindJSON(&body)
	var u, p = oAuth.LoginCustomer(&body)
	s.SendData(ctx, map[string]interface{}{
		"customer":     u,
		"access_token": p,
	})

}

func (s *AuthServerMux) handleLoginFacebookCus(ctx *gin.Context) {
	var body = oAuth.LoginFB{}
	ctx.BindJSON(&body)
	var u, p = oAuth.LoginCustomerFaceBook(&body)
	s.SendData(ctx, map[string]interface{}{
		"customer":     u,
		"access_token": p,
	})
}

func (s *AuthServerMux) handleLoginFbUpdateCus(ctx *gin.Context) {
	var body = oAuth.LoginFB{}
	ctx.BindJSON(&body)
	var u, p = oAuth.CreateCusFacebook(&body)
	s.SendData(ctx, map[string]interface{}{
		"customer":     u,
		"access_token": p,
	})
}

func (s *AuthServerMux) handleLoginGmUpdateCus(ctx *gin.Context) {
	var body = oAuth.LoginGmail{}
	ctx.BindJSON(&body)
	var u, p = oAuth.CreateCusGmail(&body)
	s.SendData(ctx, map[string]interface{}{
		"customer":     u,
		"access_token": p,
	})
}

func (s *AuthServerMux) handleLoginGmailCus(ctx *gin.Context) {
	var body = oAuth.LoginGmail{}
	ctx.BindJSON(&body)
	var u, p = oAuth.LoginCustomerGmail(&body)
	s.SendData(ctx, map[string]interface{}{
		"customer":     u,
		"access_token": p,
	})
}

func (s *AuthServerMux) handleRegisterCus(ctx *gin.Context) {
	var body = oAuth.RegisterUser{}
	ctx.BindJSON(&body)
	var u, p = oAuth.RegisterCustomer(&body)
	s.SendData(ctx, map[string]interface{}{
		"customer":     u,
		"access_token": p,
	})

}

func (s *AuthServerMux) handleLogoutCus(ctx *gin.Context) {
	var body = struct {
		//UserID string `json:"user_id"`
		Token string `json:"token"`
	}{}
	ctx.BindJSON(&body)
	rest.AssertNil(oAuth.Logout(body.Token))
	s.SendData(ctx, nil)
}
