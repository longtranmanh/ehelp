package system

import (
	"fmt"
	//"ehelp/common"
	"ehelp/x/utils"
)

func (c *CacheOrderWorker) OrderWorking(action *OrderAction) error {
	if action == nil {
		return nil
	}
	defer utils.Recover()
	defer action.Done()
	action.HandlerAction()
	if action.GetError() != nil {
		logSys.Infof(0, "validate", action)
	}
	// else if action.Action == common.ORDER_STATUS_BIDDING {
	// 	logSys.Infof(0, "validate", action.Order)
	// 	c.Orders[action.Order.ID] = action.Order
	// }
	fmt.Printf("ERROR", action.GetError())
	return action.GetError()
}
