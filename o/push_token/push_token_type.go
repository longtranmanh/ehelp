package push_token

import (
	"ehelp/x/db/mongodb"

	mgo "gopkg.in/mgo.v2"
)

type PushToken struct {
	mongodb.BaseModel `bson:",inline"`
	UserId            string `bson:"user_id" json:"user_id" validate:"required"`
	Role              int    `bson:"role" json:"role" validate:"required"`
	IsRevoke          bool   `bson:"is_revoke" json:"is_revoke"`
	DeviceId          string `bson:"device_id" json:"device_id"`
	Platform          string `bson:"platform" json:"platform"`
	VersionApp        string `bson:"version_app" json:"version_app"`
	PushToken         string `bson:"push_token" json:"push_token"`
}

var PushTokenTable = mongodb.NewTable("push_token", "k", 80)

func pushTokenTable() *mgo.Collection {
	return mongodb.NewCollection("push_token")
}
