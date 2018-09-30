package config_system

import (
	"ehelp/common"
	"ehelp/x/db/mongodb"
	"errors"
	"gopkg.in/mgo.v2/bson"
)

const (
	HOUR_MONEY                 = "hour_money"
	BONUS_TOOL_SERVICE_MONEY   = "bonus_tool_service_money"
	BONUS_TOOL_SERVICE_PERCENT = "bonus_tool_service_percent"
)

type ConfigSystem struct {
	mongodb.BaseModel `bson:",inline"`
	Name              string `bson:"name" json:"name" validate:"required"`
	Description       string `bson:"description" json:"description"`
	Value             string `bson:"value" json:"value" validate:"required"`
}

var ConfigSystemTable = mongodb.NewTable("config_system", "cnf", 12)

func GetAllConfig() ([]*ConfigSystem, error) {
	var confs []*ConfigSystem
	var err = ConfigSystemTable.FindWhere(bson.M{}, confs)
	if err != nil {
		if err.Error() == common.NOT_EXIST {
			return nil, errors.New("Không có data")
		}
		return nil, err
	}
	return confs, nil
}
