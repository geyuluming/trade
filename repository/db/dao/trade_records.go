package dao

import (
	"context"
	//"errors"
	//"fmt"
	"github.com/kasiforce/trade/repository/db/model"
	"github.com/kasiforce/trade/types"
	"gorm.io/gorm"
	"time"
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

	var orders []struct {
		TradeID            int
		SellerID           int
		BuyerID            int
		SellerName         string
		BuyerName          string
		GoodsID            int
		GoodsName          string
		Price              float64
		DeliveryMethod     string
		ShippingCost       float64
		ShippingProvince   string
		ShippingCity       string
		ShippingArea       string
		ShippingDetailArea string
		DeliveryProvince   string
		DeliveryCity       string
		DeliveryArea       string
		DeliveryDetailArea string
		OrderTime          time.Time
		PayTime            time.Time
		ShippingTime       time.Time
		TurnoverTime       time.Time
		Status             string
	}

	err = query.
		Joins("left join users as seller on seller.userID = trade_records.sellerID").
		Joins("left join users as buyer on buyer.userID = trade_records.buyerID").
		Joins("left join goods on goods.goodsID = trade_records.goodsID").
		Joins("left join address as shippingAddr on shippingAddr.addrID = trade_records.shippingAddrID").
		Joins("left join address as deliveryAddr on deliveryAddr.addrID = trade_records.deliveryAddrID").
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
			"shippingAddr.province as ShippingProvince," +
			"shippingAddr.city as ShippingCity," +
			"shippingAddr.districts as ShippingArea," +
			"shippingAddr.address as ShippingDetailArea," +
			"deliveryAddr.province as DeliveryProvince," +
			"deliveryAddr.city as DeliveryCity," +
			"deliveryAddr.districts as DeliveryArea," +
			"deliveryAddr.address as DeliveryDetailArea," +
			"trade_records.orderTime as OrderTime," +
			"trade_records.payTime as PayTime," +
			"trade_records.shippingTime as ShippingTime," +
			"trade_records.turnoverTime as TurnoverTime," +
			"trade_records.status as Status").
		Scan(&orders).Error

	if err != nil {
		return
	}

	for _, order := range orders {
		r = append(r, types.OrderInfo{
			TradeID:        order.TradeID,
			SellerID:       order.SellerID,
			BuyerID:        order.BuyerID,
			SellerName:     order.SellerName,
			BuyerName:      order.BuyerName,
			GoodsID:        order.GoodsID,
			GoodsName:      order.GoodsName,
			Price:          order.Price,
			DeliveryMethod: order.DeliveryMethod,
			ShippingCost:   order.ShippingCost,
			SenderAddress: types.AddressDetail{
				Province:   order.DeliveryProvince,
				City:       order.DeliveryCity,
				Area:       order.DeliveryArea,
				DetailArea: order.DeliveryDetailArea,
			},
			ShippingAddress: types.AddressDetail{
				Province:   order.ShippingProvince,
				City:       order.ShippingCity,
				Area:       order.ShippingArea,
				DetailArea: order.ShippingDetailArea,
			},
			OrderTime:    order.OrderTime,
			PayTime:      order.PayTime,
			ShippingTime: order.ShippingTime,
			TurnoverTime: order.TurnoverTime,
			Status:       order.Status,
		})
	}

	return
}

// UpdateOrderStatus 修改订单状态
func (c *TradeRecords) UpdateOrderStatus(req types.UpdateOrderStatusReq) (resp interface{}, err error) {
	// 更新订单状态
	err = c.DB.Model(&model.TradeRecords{}).Where("tradeID = ?", req.ID).Update("status", req.Status).Error
	if err != nil {
		return
	}

	// 如果存在退款理由，插入退款申诉
	if req.RefundReason != "" {
		refundComplaint := model.RefundComplaint{
			TradeID: req.ID,
			CReason: req.RefundReason,
			CTime:   time.Now(),
			CStatus: 0,
		}
		err = c.DB.Create(&refundComplaint).Error
		if err != nil {
			return
		}
	}

	// 如果存在评价内容，创建评论
	if req.Comment != "" {
		var tradeRecord model.TradeRecords
		err = c.DB.Where("tradeID = ?", req.ID).First(&tradeRecord).Error
		if err != nil {
			return
		}

		comment := model.Comment{
			GoodsID:        tradeRecord.GoodsID,
			CommentatorID:  tradeRecord.BuyerID,
			CommentContent: req.Comment,
			CommentTime:    time.Now(),
		}
		err = c.DB.Create(&comment).Error
		if err != nil {
			return
		}
	}

	resp = types.UpdateOrderStatusResp{
		Status: req.Status,
	}
	return
}

// UpdateOrderAddress 修改订单地址
func (c *TradeRecords) UpdateOrderAddress(req types.UpdateOrderAddressReq) (resp interface{}, err error) {
	// 更新订单地址
	err = c.DB.Model(&model.TradeRecords{}).Where("tradeID = ?", req.ID).Updates(map[string]interface{}{
		//"deliveryAddrID": req.Province + req.City + req.Area + req.DetailArea,
		"deliveryAddrID": req.AddrID,
	}).Error
	if err != nil {
		return
	}

	resp = types.UpdateOrderAddressResp{}
	return
}

// CreateOrder 生成订单
func (c *TradeRecords) CreateOrder(req types.CreateOrderReq) (resp interface{}, err error) {
	// 创建发货地址
	senderAddress := model.Address{
		Province:   req.SenderAddress.Province,
		City:       req.SenderAddress.City,
		Area:       req.SenderAddress.Area,
		DetailArea: req.SenderAddress.DetailArea,
		Name:       req.SenderAddress.Name,
		Tel:        req.SenderAddress.Tel,
	}
	err = c.DB.Create(&senderAddress).Error
	if err != nil {
		return
	}

	// 创建收货地址
	shippingAddress := model.Address{
		Province:   req.ShippingAddress.Province,
		City:       req.ShippingAddress.City,
		Area:       req.ShippingAddress.Area,
		DetailArea: req.ShippingAddress.DetailArea,
		Name:       req.ShippingAddress.Name,
		Tel:        req.ShippingAddress.Tel,
	}
	err = c.DB.Create(&shippingAddress).Error
	if err != nil {
		return
	}

	// 创建订单
	order := model.TradeRecords{
		SellerID:       req.SellerID,
		GoodsID:        req.GoodsID,
		TurnoverAmount: req.Price,
		PayMethod:      req.DeliveryMethod,
		ShippingCost:   req.ShippingCost,
		ShippingAddrID: senderAddress.AddressID,
		DeliveryAddrID: shippingAddress.AddressID,
		OrderTime:      time.Now(),
		Status:         "未发货",
	}
	err = c.DB.Create(&order).Error
	if err != nil {
		return
	}

	resp = types.CreateOrderResp{
		TradeID: order.TradeID,
	}
	return
}
