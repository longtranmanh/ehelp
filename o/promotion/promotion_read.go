package promotion

import (
	"ehelp/common"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

var errNotExistPromotion = errors.New("Mã khuyến mại không tồn tại.")
var errNotExistPromotionDay = errors.New("Mã khuyến mại chưa đến ngày sử dụng.")
var errOutExistPromotion = errors.New("Mã khuyến mại đã được sử dụng hết trong hôm nay.")
var errOutExistPromotionAll = errors.New("Mã khuyến mại đã hết hạn sử dụng.")
var errOutHour = errors.New("Ngoài khung giờ khuyến mại.")

func GetPromotionById(proId string) (*Promotion, error) {
	var prm *Promotion
	var err = PromotionTable.FindByID(proId, prm)
	if err != nil && err.Error() != ERR_NOT_EXIST {
		return nil, err
	}
	return prm, nil
}

func GetAllPromotion() ([]*Promotion, error) {
	var prms []*Promotion
	var err = PromotionTable.FindWhere(bson.M{}, &prms)
	if err != nil && err.Error() != common.NOT_EXIST {
		return nil, err
	}
	return prms, nil
}

func ValidPromotion(prmID string, prmsCustomer []string, cusId string) (*Promotion, error) {
	var isCheckExist = false
	for _, item := range prmsCustomer {
		item = strings.ToUpper(item)
		if prmID == item {
			isCheckExist = true
			break
		}
	}
	if isCheckExist == false {
		var count, err = PromotionHistoryTable.CountWhere(bson.M{
			"customer_id":  cusId,
			"promotion_id": prmID,
		})
		if err != nil && err.Error() != common.NOT_EXIST {
			//glog.Info(err)
			return nil, errNotExistPromotion
		}

		if count > 0 {
			return nil, errOutExistPromotionAll
		}
		return nil, errNotExistPromotion
	}
	var prmAndHst, err = GetProAndCheckHst(prmID, cusId)
	return prmAndHst, nil
}
func GetProAndCheckHst(prmId string, cusID string) (*Promotion, error) {
	var promotionAndHst *PromotionAndHst
	err := PromotionTable.Pipe([]bson.M{
		bson.M{
			"$match": bson.M{
				"_id":     prmId,
				"user_id": cusID,
			}},
		bson.M{
			"$lookup": bson.M{
				"from":         "promotion_hst",
				"localField":   "promotion_id",
				"foreignField": "_id",
				"as":           "prm_hst",
			},
		},
	}).One(promotionAndHst)
	if err != nil {
		if err.Error() == common.NOT_EXIST {
			return nil, errNotExistPromotion
		}
		return nil, err
	}
	var timeNowNotAsia = common.GetTimeNowVietNam().Unix()
	var timeHour = common.HourMinute()
	if promotionAndHst.CreatedAt > timeNowNotAsia {
		return nil, errNotExistPromotionDay
	}
	if promotionAndHst.CreatedAt < timeNowNotAsia {
		return nil, errOutExistPromotionAll
	}

	if float32(timeHour) < promotionAndHst.HourStart || float32(timeHour) > promotionAndHst.HourEnd {
		return nil, errOutHour
	}
	count := len(promotionAndHst.PrmsHst)

	if count >= promotionAndHst.NumberOfOrder {
		return nil, errOutExistPromotionAll
	}

	var promotion = Promotion{
		Content:       promotionAndHst.Content,
		Description:   promotionAndHst.Description,
		Discount:      promotionAndHst.Discount,
		HourEnd:       promotionAndHst.HourEnd,
		HourStart:     promotionAndHst.HourStart,
		NumberOfOrder: promotionAndHst.NumberOfOrder,
		UrlWeb:        promotionAndHst.UrlWeb,
		Title:         promotionAndHst.Title,
	}
	promotion.ID = promotionAndHst.ID
	promotion.CreatedAt = promotionAndHst.CreatedAt

	return &promotion, nil
}
