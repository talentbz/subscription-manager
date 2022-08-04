package models

import "time"

type Transactions struct {
	Id             string `json:"id"`
	UserId         string `gorm:"column:userId"`
	PriceId        string `gorm:"column:priceId"`
	UserPlanId     string `gorm:"column:userplanId"`
	SessionId      string `gorm:"column:sessionId"`
	CustomerId     string `gorm:"column:customerId"`
	Status         string
	CreatedTs      time.Time
	LastModifiedTs time.Time
}
