package area_city

import (
	"ehelp/x/db/mongodb"
	"ehelp/x/rest"
	validator "gopkg.in/go-playground/validator.v9"
	"gopkg.in/mgo.v2/bson"
)

type AreaCity struct {
	mongodb.BaseModel `bson:",inline"`
	NameArea          string `bson:"name_area" json"name_area" validate:"required"`
}

var AreaCityTable = mongodb.NewTable("area_city", "arc", 12)
var validate = validator.New()

func (s *AreaCity) CreateAreaCity() error {
	if err := validate.Struct(s); err != nil {
		return rest.BadRequestValid(err)
	}
	return AreaCityTable.Create(s)
}

func DeleteByIDAreaCity(id string) error {
	return AreaCityTable.UpdateId(id, bson.M{"$set": bson.M{"update_at": 0}})
}

func GetAreaCitys() ([]*AreaCity, error) {
	var areaCitys []*AreaCity
	err := AreaCityTable.FindWhere(bson.M{}, &areaCitys)
	if err != nil {
		return nil, err
	}
	return areaCitys, nil
}

func GetAreaSearch(area string) (*AreaCity, error) {
	var areaCity *AreaCity
	err := AreaCityTable.FindWhere(bson.M{
		"name_area": area,
	}, &areaCity)
	if err != nil {
		return nil, err
	}
	return areaCity, nil
}
