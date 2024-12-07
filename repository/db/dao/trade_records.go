package dao

import (
	"context"
	//"errors"
	//"fmt"
	"github.com/kasiforce/trade/repository/db/model"
	"github.com/kasiforce/trade/types"
	"gorm.io/gorm"
	//"time"
)

type TradeRecords struct {
	*gorm.DB
}

func NewTradeRecordsByDB(db *gorm.DB) *TradeRecords {
	return &TradeRecords{db}
}

func NewTradeRecords(ctx context.Context) *TradeRecords {
	return &TradeRecords{NewDBClient(ctx)}
}

//	func (tr *TradeRecords) FindAll() (tradeRecords []*model.TradeRecords, err error) {
//		err = tr.DB.Model(&model.TradeRecords{}).Find(&tradeRecords).Error
//		return
//	}
//
//	func (tr *TradeRecords) FindByID(id int) (t *model.TradeRecords, err error) {
//		err = tr.DB.Model(&model.TradeRecords{}).Where("tradeID = ?", id).First(&t).Error
//		return
//	}
//
//	func (tr *TradeRecords) FindBySellerID(sellerID int) (tradeRecords []*model.TradeRecords, err error) {
//		err = tr.DB.Model(&model.TradeRecords{}).Where("sellerID = ?", sellerID).Find(&tradeRecords).Error
//		return
//	}
//
//	func (tr *TradeRecords) FindByBuyerID(buyerID int) (tradeRecords []*model.TradeRecords, err error) {
//		err = tr.DB.Model(&model.TradeRecords{}).Where("buyerID = ?", buyerID).Find(&tradeRecords).Error
//		return
//	}
//
//	func (tr *TradeRecords) FindByGoodsID(goodsID int) (tradeRecords []*model.TradeRecords, err error) {
//		err = tr.DB.Model(&model.TradeRecords{}).Where("goodsID = ?", goodsID).Find(&tradeRecords).Error
//		return
//	}
//
//	func (tr *TradeRecords) CreateTradeRecord(t *model.TradeRecords) (err error) {
//		err = tr.DB.Model(&model.TradeRecords{}).Create(&t).Error
//		return
//	}
//
//	func (tr *TradeRecords) UpdateTradeRecord(id int, t *model.TradeRecords) (err error) {
//		err = tr.DB.Model(&model.TradeRecords{}).Where("tradeID = ?", id).Updates(&t).Error
//		return
//	}
//
//	func (tr *TradeRecords) DeleteTradeRecord(id int) (err error) {
//		err = tr.DB.Model(&model.TradeRecords{}).Where("tradeID = ?", id).Delete(&model.TradeRecords{}).Error
//		return
//	}
//
// GetAllOrders 获取所有订单
func (c *TradeRecords) GetAllOrders(req types.ShowOrdersReq) (r []types.OrderInfo, total int64, err error) {
	query := c.DB.Model(&model.TradeRecords{})

	if req.SearchQuery != "" {
		query = query.Where("tradeID = ?", req.SearchQuery)
	}

	err = query.Count(&total).Error
	if err != nil {
		return
	}

	err = query.
		Joins("left join users as seller on seller.userID = trade_records.sellerID").
		Joins("left join users as buyer on buyer.userID = trade_records.buyerID").
		Joins("left join goods on goods.goodsID = trade_records.goodsID").
		Offset((req.PageNum - 1) * req.PageSize).Limit(req.PageSize).
		Select("trade_records.tradeID as TradeID," +
			"trade_records.sellerID as SellerID," +
			"trade_records.buyerID as BuyerID," +
			"seller.userName as SellerName," +
			"buyer.userName as BuyerName," +
			"trade_records.goodsID as GoodsID," +
			"goods.goodsName as GoodsName," +
			"trade_records.turnoverAmount as Price," +
			"trade_records.payMethod as DeliveryMethod," +
			"trade_records.shippingCost as ShippingCost," +
			"trade_records.shippingAddress as ShippingAddress," +
			"trade_records.deliveryAddress as SenderAddress," +
			"trade_records.orderTime as OrderTime," +
			"trade_records.payTime as PayTime," +
			"trade_records.shippingTime as ShippingTime," +
			"trade_records.turnoverTime as TurnoverTime," +
			"trade_records.status as Status").
		Find(&r).Error

	return
}
