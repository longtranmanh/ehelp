package cache

// import (
// 	"ehelp/o/order"
// 	"ehelp/x/rest"
// )

// func SetCacheOrder() {
// 	var all, err = order.GetAllOrderOpenAccept()
// 	if err != nil {
// 		panic(err)
// 	}
// 	for _, s := range all {
// 		globalMap.add(s)
// 	}
// }

// func SetOrderCache(ord *order.Order) {
// 	var _, ok = globalMap.cache[ord.ID]
// 	if !ok {
// 		globalMap.add(ord)
// 	}
// }

// func mustGetEntry(id string) *entry {
// 	var e, ok = globalMap.cache[id]
// 	if ok {
// 		return e
// 	}

// 	var s, err = order.GetOrderById(id)
// 	if err != nil {
// 		rest.AssertNil(err)
// 	}

// 	Refresh(s)
// 	e = globalMap.cache[id]
// 	if e == nil {
// 		// the entry was not added to cache
// 		e = &entry{ord: s}
// 	}
// 	return e
// }

// func MustGetByID(id string) *order.Order {
// 	var e = mustGetEntry(id)
// 	return e.ord
// }
