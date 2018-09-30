package user

import (
	"ehelp/x/db/mongodb"
	"ehelp/x/rest"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type Admin struct {
	mongodb.BaseModel `bson:",inline"`
	UName             string   `bson:"uname" json:"uname"`
	Password          Password `bson:"password" json:"password"`
	Name              string   `bson:"name" json:"name"`
	Role              Role     `bson:"role" json:"role"`
}

type Role string

const (
	SUPERADMIN = Role("super-admin")
	ADMIN      = Role("admin")
)

const (
	ERR_USER_NOT_FOUND = rest.BadRequest("Sai tên đăng nhập hoặc mật khẩu")
)

var AdminTable = mongodb.NewTable("admin", "ADM", 5)

func GetByUNamePwd(uname, pwd, role string) (*Admin, error) {
	var user *Admin
	var query = bson.M{"uname": strings.ToLower(uname), "role": role}
	if role == string(ADMIN) {
		query = bson.M{
			"uname": strings.ToLower(uname),
			"$or": []bson.M{
				bson.M{"role": ADMIN},
				bson.M{"role": SUPERADMIN},
			},
		}
	}
	err := AdminTable.FindOne(query, &user)
	if err != nil {
		if rest.IsNotFound(err) {
			return nil, ERR_USER_NOT_FOUND
		}
		return nil, err
	}
	if err := user.Password.ComparePassword(pwd); err != nil {
		return nil, ERR_USER_NOT_FOUND
	}
	return user, nil
}

func GetSuperUser() *Admin {
	var user *Admin
	err := AdminTable.FindOne(bson.M{"role": SUPERADMIN}, &user)
	if err != nil {
		return nil
	}
	return user
}

func GetAdmins() ([]*Admin, error) {
	var users []*Admin
	err := AdminTable.FindWhere(bson.M{"role": ADMIN}, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (u *Admin) Create() error {
	var queryUnique = bson.M{"uname": u.UName, "role": ADMIN}
	hashed, _ := u.Password.GererateHashedPassword()
	u.Password = hashed
	return AdminTable.CreateUnique(queryUnique, u)
}

func (u *Admin) Update() error {
	hashed, _ := u.Password.GererateHashedPassword()
	u.Password = hashed
	return AdminTable.UpdateId(u.ID, bson.M{
		"$set": bson.M{
			"name":     u.Name,
			"uname":    u.UName,
			"password": u.Password,
		},
	})
}

func DeleteByID(id string) error {
	return AdminTable.DeleteByID(id)
}
