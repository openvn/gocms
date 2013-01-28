package dbctx

import (
	"errors"
	"labix.org/v2/mgo/bson"
)

func (ctx *DBCtx) SaveHouse(h *House) error {
	if len(h.Member) < 1 {
		return errors.New("dbctx: must add one host for each house")
	}
	h.HouseId = bson.NewObjectId()
	h.Member[0].Role = DefaultRole[0]

	h.Member[0].AvgIncome = 0
	if len(h.Member[0].Incomes) > 0 {
		for _, v := range h.Member[0].Incomes {
			h.Member[0].AvgIncome = h.Member[0].AvgIncome + v.Amount
		}
	}
	return ctx.houseColl.Insert(h)
}
