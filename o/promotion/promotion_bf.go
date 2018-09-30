package promotion

import (
	"ehelp/x/rest"
	validator "gopkg.in/go-playground/validator.v9"
	//"gopkg.in/mgo.v2/bson"
)

var validate = validator.New()

func (tok *Promotion) create() error {
	if err := validate.Struct(tok); err != nil {
		return rest.BadRequestValid(err)
	}
	return nil
}

func (tok *Promotion) delete() {
	tok.BeforeDelete()
}

func (tok *Promotion) update() {
	tok.BeforeUpdate()
}
