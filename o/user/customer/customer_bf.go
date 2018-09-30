package customer

import (
	"ehelp/common"
	"ehelp/o/user"
	"ehelp/x/rest"
	"ehelp/x/rest/math"
	"errors"
	validator "gopkg.in/go-playground/validator.v9"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

var validate = validator.New()

func NewInviteCode() string {
	return strings.TrimPrefix(math.RandStringUpper("", 6), "_")
}

func (cus *Customer) create() error {
	if err := validate.Struct(cus); err != nil {
		return rest.BadRequestValid(err)
	}
	err := CheckExistByPhone(cus.Phone, cus.Email)
	if err != nil {
		return err
	}
	if strings.Trim(cus.ReferenceCode, " ") != "" {
		var count, err = CheckCusByReferentCode(cus.ReferenceCode)
		rest.AssertNil(err)
		if count == 0 {
			return errors.New("Không tồn tại mã giới thiệu!")
		}
	}
	cus.InviteCode = NewInviteCode()
	psd, err := user.Password(cus.Password).GererateHashedPassword()
	if err != nil {
		return err
	}
	cus.Password = string(psd)
	// cus.BeforeCreate("", 12)
	return nil
}

func (cus *Customer) delete() {
	cus.BeforeDelete()
}

func (cus *Customer) update() {
	cus.BeforeUpdate()
}

func CheckCusByReferentCode(refCode string) (int, error) {
	return CustomerTable.CountWhere(bson.M{
		"invite_code": refCode,
	})
}

func CheckExistByPhone(phone string, email string) error {
	var count int
	var err error
	if email != "" {
		count, err = GetCusByPhone(phone, email)
	} else {
		count, err = GetByPhone(phone)
	}
	if err != nil {
		if err.Error() != common.NOT_EXIST {
			return rest.BadRequestNotFound(errors.New("Tài khoản đã tồn tại"))
		}
		return err
	}
	if count > 0 {
		return rest.BadRequestNotFound(errors.New("Tài khoản đã tồn tại"))
	}
	return nil
}
