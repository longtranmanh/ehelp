package employee

import (
	"ehelp/x/rest"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

func (emp *Employee) CrateEmployee() *Employee {
	rest.AssertNil(emp.create())
	rest.AssertNil(EmployeeTable.Create(emp))
	return emp
}

func (e *Employee) Update() error {
	return EmployeeTable.UpdateId(e.ID, e)
}

func DeleteEmployee(id string) error {
	return EmployeeTable.DeleteByID(id)
}

type EmpUpdate struct {
	AddressWork string `json:"address_work"`
	Service     string `json:"service"`
	TypePayment int    `json:"type_payment"`
	Image       string `json:"image"`
	LinkAvatar  string `json:"link_avatar"`
}

func (e *Employee) UpdateInfo(updEmp *EmpUpdate) (err error) {
	var empWork = EmployeeWork{
		AddressWork: updEmp.AddressWork,
		ServiceIds:  []string{updEmp.Service},
	}

	var mapUp = map[string]interface{}{
		"type_payment": updEmp.TypePayment,
		"image":        updEmp.Image,
		"emp_work":     empWork,
		"link_avatar":  updEmp.LinkAvatar,
	}
	err = EmployeeTable.UpdateId(e.ID, bson.M{
		"$set": mapUp,
	})
	if err == nil {
		e.EmployeeWork.AddressWork = updEmp.AddressWork
		e.EmployeeWork.ServiceIds = empWork.ServiceIds
	}
	return
}

func ActiveEmployee(id string) error {
	return EmployeeTable.UpdateId(id, bson.M{"$set": bson.M{"is_active": true}})
}
func DeactiveEmployee(id string) error {
	return EmployeeTable.UpdateId(id, bson.M{"$set": bson.M{"is_active": false}})
}
func UpdateEmployee(phone string) {
	var err, res = CheckEmployeeByPhone(phone)
	rest.AssertNil(err)
	res.update()
	rest.AssertNil(EmployeeTable.UpdateId(res.ID, bson.M{
		"$set": res,
	}))
}

func UpdateEmployeeByCusNewAndHour(empID string, countCusInOrder int, allHourItemOrder float32) {
	var queryUpdate = bson.M{}
	if countCusInOrder > 0 {
		queryUpdate["$inc"] = bson.M{
			"all_customer":  1,
			"all_hour_work": allHourItemOrder,
		}
	} else {
		queryUpdate["$inc"] = bson.M{
			"all_hour_work": allHourItemOrder,
		}
	}
	rest.AssertNil(EmployeeTable.UpdateId(empID, queryUpdate))
}

func UpdateRate(empID string, rate int) {
	var rateCol = "rate" + strconv.Itoa(rate)
	var queryUpdate = bson.M{}
	if rate > 0 {
		queryUpdate["$inc"] = bson.M{
			rateCol: 1,
		}
	}
	rest.AssertNil(EmployeeTable.UpdateId(empID, queryUpdate))
}

func (emp *Employee) DeleteEmployee() {
	emp.delete()
	rest.AssertNil(EmployeeTable.Update(emp.ID, bson.M{
		"$set": emp,
	}))
}

func (cus *Employee) UpdateEmployeeFb(fbId string, fbToken string) error {
	cus.update()
	var err = EmployeeTable.UpdateId(cus.ID, bson.M{
		"$set": map[string]interface{}{
			"fb_id":    fbId,
			"fb_token": fbToken,
		},
	})
	if err != nil {
		return errors.New("UpdateEmployeeFb error")
	}

	cus.FbToken = fbToken
	return nil
}

func (cus *Employee) UpdateEmployeeGmail(gmId string, gmToken string) error {
	cus.update()
	var err = EmployeeTable.UpdateId(cus.ID, bson.M{
		"$set": map[string]interface{}{
			"gm_id":    gmId,
			"gm_token": gmToken,
		},
	})
	if err != nil {
		return errors.New("UpdateEmployeeGmail error")
	}

	cus.GmToken = gmToken
	return nil
}

func DeleteUserByID(userID string) error {
	var err = EmployeeTable.RemoveId(userID)
	return err
}

func UploadCertificateBase64(userID string, certificate string) error {
	return EmployeeTable.UpdateId(userID, bson.M{"$set": bson.M{"certificate": certificate}})
}
func UploadAvatarBase64(userID string, avatar string, linkAvatar string) error {
	var data = bson.M{
		"avatar":      avatar,
		"image":       avatar,
		"link_avatar": linkAvatar,
	}
	return EmployeeTable.UpdateId(userID, bson.M{"$set": data})
}
