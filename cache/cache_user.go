package cache

import (
	"ehelp/o/user/customer"
	"ehelp/o/user/employee"
	"fmt"
)

var CacheEmps = map[string]*employee.Employee{}
var CacheCusts = map[string]*customer.Customer{}

func SetCacheEmp() {
	var emps, _ = employee.GetAllUser()
	fmt.Printf("Số Nv:", len(emps))
	if emps != nil && len(emps) > 0 {
		for _, item := range emps {
			CacheEmps[item.ID] = item
		}
	}
}

func SetCacheCus() {
	var emps, _ = customer.GetAllCus()
	fmt.Printf("Số khách:", len(emps))
	if emps != nil && len(emps) > 0 {
		for _, item := range emps {
			CacheCusts[item.ID] = item
		}
	}
}

func GetEmpID(empID string) (*employee.Employee, error) {
	if val, ok := CacheEmps[empID]; ok {
		return val, nil
	} else {
		val, err := employee.GetByID(empID)
		if val != nil {
			CacheEmps[empID] = val
		}
		return val, err
	}
}

func GetCusID(cusID string) (*customer.Customer, error) {
	if val, ok := CacheCusts[cusID]; ok {
		fmt.Printf("KHACH == ", val)
		return val, nil
	} else {
		val, err := customer.GetByID(cusID)
		if val != nil {
			CacheCusts[cusID] = val
		}
		fmt.Printf("KHACH == ", val)
		return val, err
	}
}
