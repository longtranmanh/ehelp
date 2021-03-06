package mongodb

import (
	"ehelp/x/rest"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	ERR_EXIST = rest.BadRequest("Đã tồn tại")
)

type Table struct {
	*mgo.Collection
	Name   string
	Prefix string
	Length int
}

func NewTable(name, prefix string, length int) *Table {
	return &Table{
		Collection: NewCollection(name),
		Name:       name,
		Prefix:     prefix,
		Length:     length,
	}
}

func (t *Table) Create(model IModel) error {
	model.BeforeCreate(t.Prefix, t.Length)
	return t.Insert(model)
}

func (t *Table) CreateUnique(query bson.M, model IModel) error {
	count, err := t.CountWhere(query)
	if err == nil {
		if count == 0 {
			return t.Create(model)
		}
		return ERR_EXIST
	}
	return err
}

func (t *Table) CountWhere(query bson.M) (int, error) {
	query["updated_at"] = bson.M{
		"$ne": 0,
	}
	return t.Find(query).Count()
}

func (t *Table) FindWhere(query bson.M, result interface{}) error {
	query["updated_at"] = bson.M{
		"$ne": 0,
	}
	return t.Find(query).Sort("-updated_at").All(result)
}
func (t *Table) FindOne(query bson.M, result interface{}) error {
	query["updated_at"] = bson.M{
		"$ne": 0,
	}
	return t.Find(query).One(result)
}
func (t *Table) FindByID(id string, result interface{}) error {
	return t.FindId(id).One(result)
}

func (t *Table) DeleteByID(id string) error {
	return t.UpdateId(id, bson.M{"$set": bson.M{"updated_at": 0}})
}
func (t *Table) UpdateByID(id string, i IModel) error {
	i.BeforeUpdate()
	return t.UpdateId(id, i)
}
