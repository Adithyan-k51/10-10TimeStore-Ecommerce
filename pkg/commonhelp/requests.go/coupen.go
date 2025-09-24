package requests

import "time"

type Coupon struct {
	Code                 string    `json:"code" binding:"required"`
	DiscountPercent      float64   `json:"discount_percent" binding:"required"`
	UsageLimits          int       `json:"usage_limits" binding:"required"`
	MaximumDiscountPrice float64   `json:"maximum_discount_price" binding:"required"`
	MinimumPurchasePrice float64   `json:"minimum_purchase_price" binding:"required"`
	ExpiryDate           time.Time `json:"expiry_date" binding:"required"`
}
