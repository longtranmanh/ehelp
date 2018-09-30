package common

const (
	ERROR_DONT_ROLE = "Quyền không hợp lệ !"
	NOT_EXIST       = "not found" // không có data
)

var LINK_AVATAR string

type OrderStatus string
type ItemOrderStatus string

const (
	ORDER_STATUS_OPEN          = OrderStatus("missed")
	ORDER_STATUS_BIDDING       = OrderStatus("bidding")
	ORDER_STATUS_ACCEPTED      = OrderStatus("accepted")
	ORDER_STATUS_WORKING       = OrderStatus("working")
	ORDER_STATUS_FINISHED      = OrderStatus("finished")
	ORDER_STATUS_CANCELED      = OrderStatus("canceled")
	ITEM_ORDER_STATUS_NEW      = ItemOrderStatus("new")
	ITEM_ORDER_STATUS_WORKING  = ItemOrderStatus("working")
	ITEM_ORDER_STATUS_FINISHED = ItemOrderStatus("finished")
)

type TypeWork int

const (
	TYPE_ONE_WEEK   = TypeWork(1) // kiểu 1 tuần 1 lần
	TYPE_TWO_WEEK   = TypeWork(2) // kiểu 2 tuần 1 lần
	TYPE_THREE_WEEK = TypeWork(3) // kiểu 3 tháng 1 lần
	TYPE_FOUR_WEEK  = TypeWork(4) // kiểu 4 tuần 1 lần
)

type TypePayment int

const (
	TYPE_PAYMENT_CARD  = TypePayment(1) // kiểu thanh toán theo thẻ
	TYPE_PAYMENT_MONEY = TypePayment(2) //Thanh toán theo trả lương
)
