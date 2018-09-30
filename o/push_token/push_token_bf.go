package push_token

import (
	"ehelp/x/rest"
	"ehelp/x/rest/validator"

	"gopkg.in/mgo.v2/bson"
)

func (tok *PushToken) create() {
	rest.AssertNil(validator.Validate(tok))
	tok.IsRevoke = false
	// tok.BeforeCreate("k", 80)
}

func (tok *PushToken) delete() {
	tok.BeforeDelete()
}

func (tok *PushToken) update() {
	tok.BeforeUpdate()
	tok.IsRevoke = true
}

func CheckTokenRevoke(token string) (error, *PushToken) {
	var tok PushToken
	return PushTokenTable.FindOne(bson.M{
		"id":        token,
		"is_revoke": true,
	}, &tok), &tok
}

func CheckTokenByUserId(userId string, role int) (error, *PushToken) {
	var tok PushToken
	return PushTokenTable.FindOne(bson.M{
		"user_id":   userId,
		"role":      role,
		"is_revoke": false,
	}, &tok), &tok
}
