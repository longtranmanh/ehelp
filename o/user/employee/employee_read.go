package employee

import (
	"ehelp/common"
	"ehelp/o/user"
	"ehelp/x/mlog"
	"ehelp/x/rest"
	"errors"
	"gopkg.in/mgo.v2/bson"
)

var logEmp = mlog.NewTagLog("Employee")

func GetEmployeeByLogin(phone string, password string) (error, *Employee) {
	var cus *Employee
	err := EmployeeTable.FindOne(bson.M{"phone": phone}, &cus)
	if err != nil || cus == nil {
		err = rest.BadRequestNotFound(errors.New("tài khoản không tồn tại!"))
		logEmp.Errorf("GetEmployeeByLogin", err)
		return err, nil
	}
	var psd = user.Password(cus.Password)
	if err := psd.ComparePassword(password); err != nil {
		err = rest.BadRequestValid(errors.New("Password sai!"))
		logEmp.Errorf("GetEmployeeByLogin", err)
		return err, nil
	}
	// if cus != nil && cus.IsActive == false {
	// 	return rest.BadRequestValid(errors.New("Tài khoản chưa được kích hoạt!")), nil
	// }
	return nil, cus
}

func CheckEmployeeById(userId string) (error, *Employee) {
	var cus Employee
	return EmployeeTable.FindOne(bson.M{
		"id": userId,
	}, &cus), &cus
}

func CheckEmployeeByPhone(phone string) (error, *Employee) {
	var emp Employee
	return EmployeeTable.FindOne(bson.M{
		"phone": phone,
	}, &emp), &emp
}

func GetEmployeeByLoginFb(fbId string) (error, *Employee) {
	var cus *Employee
	err := EmployeeTable.FindOne(bson.M{"fb_id": fbId}, &cus)
	if err != nil {
		logEmp.Errorf("GetEmployeeByLoginFb", err)
		if err.Error() == common.NOT_EXIST {
			return rest.BadRequestNotFound(errors.New("tài khoản không tồn tại!")), nil
		}
		return err, nil
	}
	return nil, cus
}

func GetEmployeeByLoginGmail(gmId string) (error, *Employee) {
	var cus *Employee
	err := EmployeeTable.FindOne(bson.M{"gm_id": gmId}, &cus)
	if err != nil {
		logEmp.Errorf("GetEmployeeByLoginGmail", err)
		if err.Error() == common.NOT_EXIST {
			return rest.BadRequestNotFound(errors.New("tài khoản không tồn tại!")), nil
		}
		return err, nil
	}
	return nil, cus
}

func GetEmpByPhone(phone string) (int, error) {
	return EmployeeTable.CountWhere(bson.M{"phone": phone})
}

func GetEmployees() ([]*Employee, error) {
	var employees []*Employee
	err := EmployeeTable.FindWhere(bson.M{}, &employees)
	if err != nil && err.Error() == common.NOT_EXIST {
		err = nil
	} else {
		logEmp.Errorf("GetEmployees", err)
	}
	return employees, err
}

func GetByID(userId string) (*Employee, error) {
	var emp *Employee
	err := EmployeeTable.FindByID(userId, &emp)
	if err != nil && err.Error() == common.NOT_EXIST {
		err = nil
	} else {
		logEmp.Errorf("GetByID", err)
	}
	return emp, err
}

func (usr *Employee) GetMyAvgRate() float32 {
	var rate = usr.Rate1 + usr.Rate2 + usr.Rate3 + usr.Rate4 + usr.Rate5
	if rate == 0 {
		return 5
	} else {
		var temp = (float32)(usr.Rate1*1+usr.Rate2*2+usr.Rate3*3+usr.Rate4*4+usr.Rate5*5) / (float32)(rate)
		temp = common.Round2(temp, .5)
		if temp < 0 {
			return 5.00
		}
		return temp
	}
}

func GetListEmpVsOrderBidding(serviceOrder []string, addressOrder string) (empIds []string, err error) {
	err = EmployeeTable.Find(bson.M{
		"is_active":            true,
		"emp_work.service_ids": serviceOrder[0],
		//"$text":                bson.M{"$search": addressOrder}, // đây là index "emp_work.address_work"
	}).Distinct("_id", &empIds)
	if err != nil && err.Error() == common.NOT_EXIST {
		err = nil
	} else {
		logEmp.Errorf("GetListEmpVsOrderBidding", err)
	}
	return
}

func GetAllUser() (emps []*Employee, err error) {
	err = EmployeeTable.FindWhere(bson.M{"updated_at": bson.M{
		"$ne": 0,
	}}, &emps)
	if err != nil && err.Error() == common.NOT_EXIST {
		err = nil
	} else if err != nil {
		logEmp.Errorf("GetAllUser", err)
	}
	return
}
