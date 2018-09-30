package customer

import (
	"ehelp/common"
	"ehelp/x/rest"
	"errors"
	"gopkg.in/mgo.v2/bson"
)

func ActiveCustomer(id string) error {
	return CustomerTable.UpdateId(id, bson.M{"$set": bson.M{"is_active": true}})
}
func DeactiveCustomer(id string) error {
	return CustomerTable.UpdateId(id, bson.M{"$set": bson.M{"is_active": false}})
}
func (cus *Customer) CrateCustomer() *Customer {
	rest.AssertNil(cus.create())
	rest.AssertNil(CustomerTable.Create(cus))
	return cus
}

func (c *Customer) Update() error {
	return CustomerTable.UpdateId(c.ID, c)
}

func (cus *Customer) UpdateCustomer(value interface{}) {
	cus.update()
	rest.AssertNil(CustomerTable.UpdateId(cus.ID, bson.M{
		"$set": value,
	}))
}

func (cus *Customer) DeleteCustomer() {
	cus.delete()
	rest.AssertNil(CustomerTable.Insert(cus))
}

func (cus *Customer) UpdateCustomerFb(fbId string, fbToken string) error {
	cus.update()
	var err = CustomerTable.UpdateId(cus.ID, bson.M{
		"$set": map[string]interface{}{
			"fb_id":    fbId,
			"fb_token": fbToken,
		},
	})
	if err != nil && err.Error() != common.NOT_EXIST {
		return errors.New("UpdateCustomerFb error")
	}
	cus.FbToken = fbToken
	return nil
}

func (cus *Customer) UpdateCustomerGmail(gmId string, gmToken string) error {
	cus.update()
	var err = CustomerTable.UpdateId(cus.ID, bson.M{
		"$set": map[string]interface{}{
			"gm_id":    gmId,
			"gm_token": gmToken,
		},
	})
	if err != nil && err.Error() != common.NOT_EXIST {
		return errors.New("UpdateEmployeeGmail error")
	}

	cus.GmToken = gmToken
	return nil
}

func DeleteUserByID(userID string) error {
	var err = CustomerTable.RemoveId(userID)
	return err
}
