package system

import (
	"ehelp/common"
	"ehelp/o/order"
	"ehelp/x/mrw/encode"
	"fmt"
)

type actionHandler func(action *OrderAction)
type actionHandlers map[common.OrderStatus]actionHandler

type OrderAction struct {
	Action  common.OrderStatus
	OrderID string
	EmpId   string
	Order   *order.OrderCheckChange `json:"order"`
	Error   OrderActionError        `json:"error"`
	Extra   encode.RawMessage       `json:"extra"`
	doneC   chan struct{}
	used    bool // must be trigger at most once
}

type OrderActionError struct {
	s string
}

func (a *OrderAction) Done() bool {
	fmt.Printf("VO ACTION DONE")
	a.doneC <- struct{}{}
	return a.GetError() == nil
}

func (e *OrderActionError) Error() string {
	return e.s
}

func (e *OrderActionError) GetError() error {
	if len(e.s) > 0 {
		return e
	}
	return nil
}

func (a *OrderAction) SetError(err error) {
	if err == nil {
		return
	}
	a.Error = OrderActionError{s: err.Error()}
}

func (a *OrderAction) GetError() error {
	return a.Error.GetError()
}

func (a *OrderAction) Wait() (*order.Order, error) {

	if a.doneC == nil {
		panic("no done channel")
	}
	<-a.doneC
	var err = a.GetError()
	if err != nil {
		fmt.Printf("LOI :", err)
		return nil, err
	}
	if a.Order == nil {
		fmt.Println("KHÔNG CÓ ORDER")
		return nil, nil
	}
	return a.Order.Order, nil
}

func NewOrderAction() *OrderAction {
	var a = &OrderAction{
		doneC: make(chan struct{}, 1),
	}
	return a
}
