package models

import (
	"time"

	"gorm.io/gorm"
)

type InvoiceItem struct {
	gorm.Model
	InvoiceID uint    `json:"invoice_id" gorm:"not null"`
	ItemID    uint    `json:"item_id" gorm:"not null"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	Price     float64 `json:"price" gorm:"not null"`
}

type Invoice struct {
	gorm.Model
	InvoiceNumber string        `json:"invoice_number" gorm:"unique;not null"`
	CustomerName  string        `json:"customer_name" gorm:"not null"`
	CustomerPhone string        `json:"customer_phone"`
	Date          time.Time     `json:"date" gorm:"not null"`
	Items         []InvoiceItem `json:"items" gorm:"foreignKey:InvoiceID"`
	Total         float64       `json:"total" gorm:"not null"`
}
