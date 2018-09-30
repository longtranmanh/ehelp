package service

import (
	"ehelp/o/tool"
	"ehelp/x/db/mongodb"
	"ehelp/x/rest"
	"errors"
	"github.com/golang/glog"
	validator "gopkg.in/go-playground/validator.v9"
	"gopkg.in/mgo.v2/bson"
)

const NOT_EXIST = "not found"

var ServiceTable = mongodb.NewTable("service", "srv", 12)

type Service struct {
	mongodb.BaseModel `bson:",inline"`
	Name              string   `bson:"name" json:"name" validate:"required"`
	IconLink          string   `bson:"icon_link" json:"icon_link"`
	NodeServices      []string `bson:"node_services" json:"node_services"`
	PricePerHour      int      `bson:"price_per_hour" json:"price_per_hour" validate:"required"`
	Tools             []string `bson:"tools" json:"tools"`
}

var validate = validator.New()

func (s *Service) Create() error {
	if err := validate.Struct(s); err != nil {
		glog.Error(err)
		return rest.BadRequestValid(err)
	}
	return ServiceTable.CreateUnique(bson.M{"name": s.Name}, s)
}

func UpdateIconLink(id string, iconLink string) error {
	var data = bson.M{
		"icon_link": iconLink,
	}
	return ServiceTable.UpdateId(id, bson.M{"$set": data})
}

func (s *Service) Update() error {
	return ServiceTable.UpdateByID(s.ID, s)
}

func GetServices() ([]*Service, error) {
	var services []*Service
	err := ServiceTable.FindWhere(bson.M{}, &services)
	if err != nil {
		return nil, err
	}
	return services, nil
}
func GetToolServices(types string) ([]*Service, error) {
	var services []*Service
	var match = bson.M{
		"$match": bson.M{},
	}
	if types == "" {
		match["$match"] = bson.M{
			"updated_at": bson.M{
				"$ne": 0,
			},
		}
	}
	err := ServiceTable.Pipe([]bson.M{
		match,
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
				"node_services":  "$node_services",
				"tools":          "$tools.name",
			},
		},
	}).All(&services)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	if services == nil {
		return []*Service{}, nil
	}
	return services, nil
}

type ServiceTool struct {
	mongodb.BaseModel `bson:",inline"`
	NodeServices      []string    `bson:"node_services" json:"node_services"`
	Name              string      `bson:"name" json:"name" validate:"required"`
	PricePerHour      int         `bson:"price_per_hour" json:"price_per_hour" validate:"required"`
	IconLink          string      `bson:"icon_link" json:"icon_link"`
	Tools             []tool.Tool `bson:"tools" json:"tools"`
}

func GetServiceAndTool(idServices []string) ([]*ServiceTool, error) {
	var services []*ServiceTool

	err := ServiceTable.Pipe([]bson.M{
		bson.M{
			"$match": bson.M{
				"_id": bson.M{"$in": idServices},
				"updated_at": bson.M{
					"$ne": 0,
				},
			}},
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
				"node_services":  "$node_services",
				"tools":          "$tools",
			},
		},
	}).All(&services)
	if err != nil {
		if err.Error() == NOT_EXIST {
			return nil, errors.New("Dịch vụ không hỗ trợ!")
		}
	}
	return services, nil
}

func DeleteServiceByID(id string) error {
	return ServiceTable.DeleteByID(id)
}

func GetAllServiceAndTool() ([]*ServiceTool, error) {
	var services []*ServiceTool
	var match = bson.M{}
	match["updated_at"] = bson.M{"$ne": 0}
	err := ServiceTable.Pipe([]bson.M{
		bson.M{"$match": match},
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
				"_id":            1,
				"name":           1,
				"node_services":  1,
				"price_per_hour": 1,
				"icon_link":      1,
				"tools":          "$tools",
			},
		},
	}).All(&services)
	if err != nil {
		if err.Error() == NOT_EXIST {
			return nil, errors.New("Dịch vụ không hỗ trợ!")
		}
		return nil, err
	}
	return services, nil
}

func GetByID(idSer string) (srs *Service, err error) {
	err = ServiceTable.FindByID(idSer, &srs)
	if err != nil {
		if err.Error() == NOT_EXIST {
			err = errors.New("Không tồn tại service")
		}
		return nil, err
	}
	return srs, nil
}
