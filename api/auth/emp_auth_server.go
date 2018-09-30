package auth

import (
	//"bytes"
	"ehelp/common"
	oAuth "ehelp/o/user/auth"
	"ehelp/o/user/employee"
	//"ehelp/x/fcm"
	"ehelp/x/rest"
	//"encoding/base64"
	//"errors"
	"github.com/gin-gonic/gin"
	//"image/jpeg"
	//"os"
	//"strings"
)

type EmpAuthServerMux struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewAuthEmployeeServer(parent *gin.RouterGroup, name string) {
	var s = EmpAuthServerMux{
		RouterGroup: parent.Group(name),
	}
	s.POST("/login", s.handleLoginEmp)
	s.POST("/logout", s.handleLogoutEmp)
	s.POST("/login_facebook", s.handleLoginFacebookEmp)
	s.POST("/loginfb_update", s.handleLoginFbUpdateEmp)
	s.POST("/login_gmail", s.handleLoginGmailEmp)
	s.POST("/logingm_update", s.handleLoginGmUpdateEmp)
	s.POST("/register", s.handleRegisterEmp)
	s.POST("/update_info", s.handleUpdateInfo)
}

func (s *EmpAuthServerMux) handleUpdateInfo(ctx *gin.Context) {
	var emp = oAuth.GetEmpFromToken(ctx.Request)
	var up *employee.EmpUpdate
	rest.AssertNil(ctx.ShouldBindJSON(&up))
	if up.TypePayment == 0 {
		up.TypePayment = int(common.TYPE_PAYMENT_MONEY)
	}
	// b64data := up.Image[strings.IndexByte(up.Image, ',')+1:]
	// unbased, _ := base64.StdEncoding.DecodeString(string(b64data))
	// var errs = errors.New("Lỗi upload ảnh! Thử lại!")
	// jpgI, err := jpeg.Decode(bytes.NewReader(unbased))
	// if err != nil {
	// 	rest.AssertNil(errs)
	// }
	// f, err := os.Create("./upload/avatar/" + emp.ID + ".jpg")
	// if err != nil {
	// 	rest.AssertNil(errs)
	// }
	// defer f.Close()
	// err = jpeg.Encode(f, jpgI, &jpeg.Options{Quality: 75})
	// if err != nil {
	// 	rest.AssertNil(errs)
	// }
	// up.LinkAvatar = fcm.LINK_AVATAR + emp.ID + ".jpg"
	rest.AssertNil(emp.UpdateInfo(up))
	s.Success(ctx)
}

func (s *EmpAuthServerMux) handleLoginEmp(ctx *gin.Context) {
	var body = oAuth.LoginUser{}
	ctx.BindJSON(&body)
	var u, p = oAuth.LoginEmployee(&body)
	s.SendData(ctx, map[string]interface{}{
		"employee":     u,
		"access_token": p,
	})
}
func (s *EmpAuthServerMux) handleLoginFbUpdateEmp(ctx *gin.Context) {
	var body = oAuth.LoginFB{}
	ctx.BindJSON(&body)
	var u, p = oAuth.CreateEmpFacebook(&body)
	s.SendData(ctx, map[string]interface{}{
		"employee":     u,
		"access_token": p,
	})
}
func (s *EmpAuthServerMux) handleLoginFacebookEmp(ctx *gin.Context) {
	var body = oAuth.LoginFB{}
	ctx.BindJSON(&body)
	var u, p = oAuth.LoginEmployeeFaceBook(&body)
	s.SendData(ctx, map[string]interface{}{
		"employee":     u,
		"access_token": p,
	})
}
func (s *EmpAuthServerMux) handleLoginGmUpdateEmp(ctx *gin.Context) {
	var body = oAuth.LoginGmail{}
	ctx.BindJSON(&body)
	var u, p = oAuth.CreateEmpGmail(&body)
	s.SendData(ctx, map[string]interface{}{
		"employee":     u,
		"access_token": p,
	})
}
func (s *EmpAuthServerMux) handleLoginGmailEmp(ctx *gin.Context) {
	var body = oAuth.LoginGmail{}
	ctx.BindJSON(&body)
	var u, p = oAuth.LoginCustomerGmail(&body)
	s.SendData(ctx, map[string]interface{}{
		"employee":     u,
		"access_token": p,
	})
}
func (s *EmpAuthServerMux) handleRegisterEmp(ctx *gin.Context) {
	var body = oAuth.RegisterUser{}
	ctx.BindJSON(&body)
	var u, p = oAuth.RegisterEmployee(&body)
	s.SendData(ctx, map[string]interface{}{
		"employee":     u,
		"access_token": p,
	})
}
func (s *EmpAuthServerMux) handleLogoutEmp(ctx *gin.Context) {
	var body = struct {
		//UserID string `json:"user_id"`
		Token string `json:"token"`
	}{}
	ctx.BindJSON(&body)
	rest.AssertNil(oAuth.Logout(body.Token))
	s.SendData(ctx, nil)
}
