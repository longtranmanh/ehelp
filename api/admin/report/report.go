package report

import (
	"ehelp/o/order_hst"
	"ehelp/x/rest"
	"github.com/gin-gonic/gin"
)

type ReportServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewReportServer(parent *gin.RouterGroup, name string) {
	var s = &ReportServer{
		RouterGroup: parent.Group(name),
	}
	s.GET("/order-history", s.handleOrderHistory)
	s.GET("/general", s.handleGeneralReport)
}
func (s *ReportServer) handleOrderHistory(ctx *gin.Context) {
	var res, err = order_hst.GetOrderREport()
	rest.AssertNil(err)
	s.SendData(ctx, res)
}
func (s *ReportServer) handleGeneralReport(ctx *gin.Context) {
	var types = ctx.Query("types")
	var res, err = order_hst.GetGeneralReportByTime(types)
	rest.AssertNil(err)
	s.SendData(ctx, res)
}
