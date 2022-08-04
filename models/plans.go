package models

import "time"

type AvailablePlans struct {
	Id          int `gorm:"primaryKey"`
	Name        string
	Description string
	Price       float32
	Recurrence  int
	PriceId     string `gorm:"column:priceId"`
}

type UserPlans struct {
	Id             string `json:"id"`
	UserId         string `gorm:"column:userId"`
	PlanId         int    `gorm:"column:planId"`
	CustomerId     string `gorm:"column:customerId"`
	PriceId        string `gorm:"column:priceId"`
	SubscriptionId string `gorm:"column:subscriptionId"`
	Status         string
	CreatedTs      time.Time
	LastModifiedTs time.Time
}
