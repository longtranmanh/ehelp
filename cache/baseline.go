package cache

import (
	"ehelp/common"
	"ehelp/o/order"
	"ehelp/x/rest"
)

func MoveOrderToOpen() {
	var idOder, err = order.GetAllOrderToExpired()
	rest.AssertNil(rest.WrapBadRequest(err, ""))
	var errs, _ = order.UpdateStatusOrders(idOder, string(common.ORDER_STATUS_OPEN))
	rest.AssertNil(errs)
}
