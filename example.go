package main

import "time"

type AffiliateLogs struct {
	ID             int        `db:"id" gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID         *int       `db:"user_id" gorm:"column:user_id" json:"userID"`
	GuestID        *int       `db:"guest_id" gorm:"column:guest_id" json:"guestID"`
	ReferredByUser int        `db:"referred_by_user" gorm:"column:referred_by_user" json:"referredByUser"`
	Amount         *float64   `db:"amount" gorm:"column:amount" json:"amount"`
	OrderID        *int       `db:"order_id" gorm:"column:order_id" json:"orderID"`
	OrderDetailID  *int       `db:"order_detail_id" gorm:"column:order_detail_id" json:"orderDetailID"`
	AffiliateType  string     `db:"affiliate_type" gorm:"column:affiliate_type" json:"affiliateType"`
	Status         int        `db:"status" gorm:"column:status" json:"status"`
	CreatedAt      *time.Time `db:"created_at" gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      *time.Time `db:"updated_at" gorm:"column:updated_at" json:"updatedAt"`
}
