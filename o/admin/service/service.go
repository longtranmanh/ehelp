package service

import (
	"ehelp/x/db/mongodb"
	"ehelp/x/rest"
	"github.com/golang/glog"
	validator "gopkg.in/go-playground/validator.v9"
	"gopkg.in/mgo.v2/bson"
)

var ServiceTable = mongodb.NewTable("service", "srv", 12)
var ToolTable = mongodb.NewTable("tool", "tbx", 12)

type Service struct {
	mongodb.BaseModel `bson:",inline"`
	Name              string   `bson:"name" json:"name" validate:"required"`
	PricePerHour      int      `bson:"price_per_hour" json:"price_per_hour" validate:"required"`
	NodeServices      []string `bson:"node_services" json:"node_services"`
	Tools             []string `bson:"tools" json:"tools"`
}
type Tool struct {
	mongodb.BaseModel `bson:",inline"`
	Name              string `bson:"name" json:"name" validate:"required"`
	Price             int    `bson:"price" json:"price" validate:"required"`
}

var validate = validator.New()

func (s *Service) Create() error {
	if err := validate.Struct(s); err != nil {
		glog.Error(err)
		return rest.BadRequestValid(err)
	}
	return ServiceTable.Create(s)
}

func DeleteByID(id string) error {
	return ServiceTable.UpdateId(id, bson.M{"$set": bson.M{"update_at": 0}})
}
func (t *Tool) Create() error {
	if err := validate.Struct(t); err != nil {
		glog.Error(err)
		return rest.BadRequestValid(err)
	}
	return ToolTable.Create(t)
}
func GetServices() ([]*Service, error) {
	var services []*Service
	err := ServiceTable.FindWhere(bson.M{}, &services)
	if err != nil {
		return nil, err
	}
	return services, nil
}
func GetToolServices() ([]*Service, error) {
	var services []*Service
	err := ServiceTable.Pipe([]bson.M{
		bson.M{
			"$lookup": bson.M{
				"from":         "tool",
				"localField":   "tools",
				"foreignField": "_id",
				"as":           "tools",
			},
		},
		bson.M{
			"$project": bson.M{
				"_id":            "$_id",
				"name":           "$name",
				"price_per_hour": "$price_per_hour",
				"tools":          "$tools.name",
			},
		},
	}).All(&services)
	if err != nil {
		return nil, err
	}
	return services, nil
}
func GetTools() ([]*Tool, error) {
	var tools []*Tool
	err := ToolTable.FindWhere(bson.M{}, &tools)
	if err != nil {
		return nil, err
	}
	return tools, nil
}
