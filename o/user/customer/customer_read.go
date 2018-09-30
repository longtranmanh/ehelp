package customer

import (
	"ehelp/common"
	"ehelp/o/user"
	"ehelp/x/rest"
	"errors"
	"gopkg.in/mgo.v2/bson"
)

func CheckCustomer(userId string) (error, *Customer) {
	var cus Customer
	err := CustomerTable.FindOne(bson.M{"id": userId}, &cus)
	if err != nil {
		return err, nil
	}
	return nil, &cus
}

func GetCustomerByLogin(phone string, password string) (error, *Customer) {
	var cus *Customer
	err := CustomerTable.FindOne(bson.M{"phone": phone}, &cus)
	if err != nil {
		return rest.BadRequestNotFound(errors.New("tài khoản không tồn tại!")), nil
	}
	var psd = user.Password(cus.Password)
	if err := psd.ComparePassword(password); err != nil {
		return rest.BadRequestValid(errors.New("Password sai!")), nil
	}
	return nil, cus
}

func GetCusByPhone(phone string, email string) (int, error) {
	return CustomerTable.CountWhere(bson.M{
		"$or": []bson.M{
			bson.M{"phone": phone},
			bson.M{"email": email},
		},
	})
}

func GetByPhone(phone string) (int, error) {
	return CustomerTable.CountWhere(bson.M{
		"phone": phone,
	})
}

func GetCustomerByLoginFb(fbId string) (error, *Customer) {
	var cus *Customer
	err := CustomerTable.FindOne(bson.M{"fb_id": fbId}, &cus)
	if err != nil {
		return err, nil
	}
	return nil, cus
}

func GetCustomerByLoginGmail(gmId string) (error, *Customer) {
	var cus *Customer
	err := CustomerTable.FindOne(bson.M{"gm_id": gmId}, &cus)
	if err != nil {
		if err.Error() == common.NOT_EXIST {
			return rest.BadRequestNotFound(errors.New("tài khoản không tồn tại!")), nil
		}
		return err, nil
	}
	return nil, cus
}

func GetCustomers() ([]*Customer, error) {
	var customers []*Customer
	err := CustomerTable.FindWhere(bson.M{}, &customers)
	return customers, err
}

func GetByID(userID string) (*Customer, error) {
	var cus *Customer
	return cus, CustomerTable.FindByID(userID, &cus)
}

func GetAllCus() ([]*Customer, error) {
	var cus []*Customer
	return cus, CustomerTable.FindWhere(bson.M{"updated_at": bson.M{
		"$ne": 0,
	}}, &cus)
}
