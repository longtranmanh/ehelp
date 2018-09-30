package public

import (
	"ehelp/common"
	"ehelp/middleware"
	"ehelp/o/area_city"
	"ehelp/o/order"
	"ehelp/o/service"
	"ehelp/o/tool"
	"ehelp/o/user/auth"
	"ehelp/x/fcm"
	"ehelp/x/rest"
	"github.com/gin-gonic/gin"
)

type PublicServerMux struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewPublicServer(parent *gin.RouterGroup, name string) {
	var s = PublicServerMux{
		RouterGroup: parent.Group(name),
	}
	s.GET("service/list", s.handleList)
	s.POST("service/avatar", s.handleUpAvatar)
	s.GET("service_tool/list", s.handleListTool)
	s.GET("test", s.handleTEST)
	s.GET("service/list_by_services", s.handleListByService)
	s.GET("area/list_area", s.handleListArea)
	s.POST("order/suggest_price", s.handleSuggestPrice).Use(middleware.AuthenticateAppRole(auth.RoleCustomer))
}

func (s *PublicServerMux) handleUpAvatar(ctx *gin.Context) {
	var file, err = ctx.FormFile("icon")
	rest.AssertNil(err)
	var idSer = ctx.Query("id")
	var link = fcm.LINK_AVATAR + idSer
	rest.AssertNil(ctx.SaveUploadedFile(file, "./upload/avatar/"+idSer))
	rest.AssertNil(service.UpdateIconLink(idSer, link))
	s.SendData(ctx, nil)
}

func (s *PublicServerMux) handleList(ctx *gin.Context) {
	auth.MustGetKey(ctx.Request)
	services, err := service.GetAllServiceAndTool()
	rest.AssertNil(err)
	s.SendData(ctx, services)
}

func (s *PublicServerMux) handleTEST(ctx *gin.Context) {
	services, err := order.GetAllOrderToExpired()
	rest.AssertNil(err)
	var errs, count = order.UpdateStatusOrders(services, string(common.ORDER_STATUS_OPEN))
	rest.AssertNil(errs)
	s.SendData(ctx, count)
}

func (s *PublicServerMux) handleListArea(ctx *gin.Context) {
	areas, err := area_city.GetAreaCitys()
	rest.AssertNil(err)
	s.SendData(ctx, areas)
}

func (s *PublicServerMux) handleListTool(ctx *gin.Context) {
	tools, err := tool.GetTools()
	rest.AssertNil(err)
	s.SendData(ctx, tools)
}

// "tbx_nzCa378tfJly",
// "tbx_VvGEZGvrTW2j"
func (s *PublicServerMux) handleListByService(ctx *gin.Context) {
	var body = struct {
		Services []string `json:"services"`
	}{}
	ctx.BindJSON(&body)
	tools, err := service.GetServiceAndTool(body.Services)
	rest.AssertNil(err)
	s.SendData(ctx, tools)
}

func (s *PublicServerMux) handleSuggestPrice(ctx *gin.Context) {
	auth.GetCusFromToken(ctx.Request)
	var body *common.MathPriceOrder
	ctx.BindJSON(&body)
	var allHour, priceAllHour, priceTool, priceEnd, err = body.MathPriceOrder()
	rest.AssertNil(err)
	s.SendData(ctx, map[string]interface{}{
		"all_hour_work":  allHour,
		"price_all_hour": priceAllHour,
		"price_end":      priceEnd,
		"price_tool":     priceTool,
	})
}
