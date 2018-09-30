package mongodb

import (
	"ehelp/x/rest/math"
	"time"
)

type IModel interface {
	BeforeCreate(prefix string, length int)
	BeforeUpdate()
	BeforeDelete()
}
type BaseModel struct {
	ID        string `json:"id" bson:"_id"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
	UpdatedAt int64  `json:"updated_at" bson:"updated_at"`
}

func getTimeNowVietNam() time.Time {
	var LocNowVietNam, _ = time.LoadLocation("Asia/Ho_Chi_Minh")
	return time.Now().In(LocNowVietNam)
}

func (b *BaseModel) BeforeCreate(prefix string, length int) {
	b.ID = math.RandString(prefix, length)
	b.CreatedAt = getTimeNowVietNam().Unix()
	b.UpdatedAt = getTimeNowVietNam().Unix()
}

func (b *BaseModel) BeforeUpdate() {
	b.UpdatedAt = getTimeNowVietNam().Unix()
}

func (b *BaseModel) BeforeDelete() {
	b.UpdatedAt = 0
}
