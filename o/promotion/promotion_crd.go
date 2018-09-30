package promotion

import (
	"ehelp/x/rest"
	"gopkg.in/mgo.v2/bson"
)

func (tok *Promotion) CratePromotion() *Promotion {
	rest.AssertNil(tok.create())
	rest.AssertNil(PromotionTable.Create(tok))
	return tok
}

func UpdatePromotion(promotionID string) error {
	var res, err = GetPromotionById(promotionID)
	rest.AssertNil(err)
	res.update()
	return PromotionTable.UpdateId(res.ID, res)
}

func (p *PromotionHistory) CreatePrmHst() error {
	p.IsActive = true
	var err = PromotionHistoryTable.Create(p)
	if err != nil {
		//glog.Info(err)
		return err
	}
	return nil
}
