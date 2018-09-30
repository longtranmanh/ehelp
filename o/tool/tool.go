package tool

import (
	//"ehelp/common"
	"ehelp/x/db/mongodb"
	"ehelp/x/rest"
	//"errors"
	"github.com/golang/glog"
	validator "gopkg.in/go-playground/validator.v9"
	"gopkg.in/mgo.v2/bson"
)

var ToolTable = mongodb.NewTable("tool", "tbx", 12)

type Tool struct {
	mongodb.BaseModel `bson:",inline"`
	Name              string `bson:"name" json:"name" validate:"required"`
	Price             int    `bson:"price" json:"price" validate:"required"`
}

var validate = validator.New()

func (t *Tool) Create() error {
	if err := validate.Struct(t); err != nil {
		glog.Error(err)
		return rest.BadRequestValid(err)
	}
	return ToolTable.CreateUnique(bson.M{"name": t.Name}, t)
}
func (t *Tool) Update() error {
	return ToolTable.UpdateId(t.ID, t)
}
func DeleteToolByID(id string) error {
	return ToolTable.DeleteByID(id)
}

func GetTools() ([]*Tool, error) {
	var tools []*Tool
	err := ToolTable.FindWhere(bson.M{}, &tools)
	if err != nil {
		return nil, err
	}
	return tools, nil
}

func GetToolByArrayID(arrID []string) (tools []*Tool, err error) {
	err = ToolTable.FindWhere(bson.M{"_id": bson.M{
		"$in": arrID,
	}}, &tools)
	if err != nil {
		return nil, err
	}
	return tools, nil
}
