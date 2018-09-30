package user

import (
	"gopkg.in/mgo.v2/bson"
)

func (u *Admin) CreateAdmin() error {
	return AdminTable.CreateUnique(bson.M{
		"uname": u.UName,
	}, u)
}
