package user

import (
	"bytes"
	"ehelp/o/admin/user"
	"ehelp/o/user/customer"
	"ehelp/o/user/employee"
	"ehelp/x/fcm"
	"ehelp/x/rest"
	"ehelp/x/utils"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"image/jpeg"
	"os"
	"strings"
)

type UserServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewUserServer(parent *gin.RouterGroup, name string) {
	var s = &UserServer{
		RouterGroup: parent.Group(name),
	}
	s.POST("/create", s.handleCreate)
	s.POST("/update", s.handleUpdate)
	s.GET("/delete", s.handleDelete)
	s.GET("/list", s.handleList)
	s.GET("/customer/list", s.handleListCustomer)
	s.POST("/customer/update", s.handleUpdateCustomer)
	s.GET("/customer/active", s.handleActiveCustomer)
	s.GET("/customer/deactive", s.handleDeactiveCustomer)
	s.GET("/employee/list", s.handleListEmployee)
	s.POST("/employee/create", s.handleCreateEmployee)
	s.POST("/employee/update", s.handleUpdateEmployee)
	s.POST("/employee/avatar", s.handleUploadAvatar)
	s.POST("/employee/certificate", s.handleUploadCertificate)
	s.GET("/employee/active", s.handleActiveEmployee)
	s.GET("/employee/deactive", s.handleDeactiveEmployee)
	s.GET("/employee/delete", s.handleDeleteEmployee)
}

func (s *UserServer) handleCreateEmployee(ctx *gin.Context) {
	var e *employee.Employee
	rest.AssertNil(ctx.ShouldBindJSON(&e))
	e.CrateEmployee()
	go sendEmail(e.FullName, e.Email)
	s.SendData(ctx, e)
}

func sendEmail(name, email string) {
	var mail = utils.Mail{
		Subject: "Thư chào " + name,
		Body:    "Xin chào " + name,
		To:      email,
	}
	mail.Send()
}

type Upload struct {
	UserID string `json:"id"`
	Base64 string `json:"base64"`
}

func (s *UserServer) handleUploadCertificate(ctx *gin.Context) {
	// var file, err = ctx.FormFile("certificate")
	// rest.AssertNil(err)
	// rest.AssertNil(ctx.SaveUploadedFile(file, "./upload/certificate/"+ctx.Query("id")))
	var param *Upload
	ctx.BindJSON(&param)
	rest.AssertNil(employee.UploadCertificateBase64(param.UserID, param.Base64))
	s.Success(ctx)
}

func (s *UserServer) handleUploadAvatar(ctx *gin.Context) {
	// var file, err = ctx.FormFile("avatar")
	// rest.AssertNil(err)
	// rest.AssertNil(ctx.SaveUploadedFile(file, "./upload/avatar/"+ctx.Query("id")))
	var param *Upload
	ctx.BindJSON(&param)
	var up = employee.EmpUpdate{}
	up.Image = param.Base64
	b64data := up.Image[strings.IndexByte(up.Image, ',')+1:]
	fmt.Printf("ẢNH ", up.Image)
	unbased, _ := base64.StdEncoding.DecodeString(string(b64data))
	var errs = errors.New("Lỗi upload ảnh! Thử lại!")
	jpgI, err := jpeg.Decode(bytes.NewReader(unbased))
	if err != nil {
		rest.AssertNil(errs)
	}
	f, err := os.Create("./upload/avatar/" + param.UserID)
	if err != nil {
		rest.AssertNil(errs)
	}
	defer f.Close()
	err = jpeg.Encode(f, jpgI, &jpeg.Options{Quality: 75})
	if err != nil {
		rest.AssertNil(errs)
	}
	up.LinkAvatar = fcm.LINK_AVATAR + param.UserID
	rest.AssertNil(employee.UploadAvatarBase64(param.UserID, up.Image, up.LinkAvatar))
	s.Success(ctx)
}

func (s *UserServer) handleUpdateCustomer(ctx *gin.Context) {
	var c *customer.Customer
	rest.AssertNil(ctx.ShouldBindJSON(&c))
	rest.AssertNil(c.Update())
	s.SendData(ctx, c)
}

func (s *UserServer) handleUpdateEmployee(ctx *gin.Context) {
	var e *employee.Employee
	rest.AssertNil(ctx.ShouldBindJSON(&e))
	rest.AssertNil(e.Update())
	s.Success(ctx)
}

func (s *UserServer) handleDeleteEmployee(ctx *gin.Context) {
	var id = ctx.Query("id")
	rest.AssertNil(employee.DeleteEmployee(id))
	s.Success(ctx)
}
func (s *UserServer) handleDeactiveCustomer(ctx *gin.Context) {
	var id = ctx.Query("id")
	rest.AssertNil(customer.DeactiveCustomer(id))
	s.Success(ctx)
}

func (s *UserServer) handleActiveCustomer(ctx *gin.Context) {
	var id = ctx.Query("id")
	rest.AssertNil(customer.ActiveCustomer(id))
	s.Success(ctx)
}

func (s *UserServer) handleDeactiveEmployee(ctx *gin.Context) {
	var id = ctx.Query("id")
	rest.AssertNil(employee.DeactiveEmployee(id))
	s.Success(ctx)
}

func (s *UserServer) handleActiveEmployee(ctx *gin.Context) {
	var id = ctx.Query("id")
	rest.AssertNil(employee.ActiveEmployee(id))
	s.Success(ctx)
}

func (s *UserServer) handleCreate(ctx *gin.Context) {
	var admin *user.Admin
	rest.AssertNil(ctx.ShouldBindJSON(&admin))
	rest.AssertNil(admin.Create())
	s.SendData(ctx, admin)
}
func (s *UserServer) handleUpdate(ctx *gin.Context) {
	var admin *user.Admin
	rest.AssertNil(ctx.ShouldBindJSON(&admin))
	rest.AssertNil(admin.Update())
	s.SendData(ctx, admin)
}
func (s *UserServer) handleDelete(ctx *gin.Context) {
	var id = ctx.Query("id")
	rest.AssertNil(user.DeleteByID(id))
	s.Success(ctx)
}
func (s *UserServer) handleList(ctx *gin.Context) {
	admins, err := user.GetAdmins()
	rest.AssertNil(err)
	s.SendData(ctx, admins)
}
func (s *UserServer) handleListCustomer(ctx *gin.Context) {
	var customers, err = customer.GetCustomers()
	rest.AssertNil(err)
	s.SendData(ctx, customers)
}
func (s *UserServer) handleListEmployee(ctx *gin.Context) {
	var employees, err = employee.GetEmployees()
	rest.AssertNil(err)
	s.SendData(ctx, employees)
}
