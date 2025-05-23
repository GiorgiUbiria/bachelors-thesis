package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	OrderID         uint           `gorm:"not null" json:"order_id"`
	Order           Order          `gorm:"foreignKey:OrderID" json:"order"`
	Amount          float64        `gorm:"not null" json:"amount"`
	Currency        string         `gorm:"default:'USD'" json:"currency"`
	PaymentMethod   string         `gorm:"not null" json:"payment_method"`  // credit_card, debit_card, paypal, etc.
	Status          string         `gorm:"default:'pending'" json:"status"` // pending, completed, failed, refunded
	TransactionID   string         `gorm:"uniqueIndex" json:"transaction_id"`
	GatewayResponse string         `json:"gateway_response,omitempty"`
	ProcessedAt     *time.Time     `json:"processed_at,omitempty"`
	RefundedAt      *time.Time     `json:"refunded_at,omitempty"`
	RefundAmount    float64        `gorm:"default:0" json:"refund_amount"`
	RefundReason    string         `json:"refund_reason,omitempty"`

	// Card details (encrypted/tokenized in production)
	CardLast4       string `json:"card_last4,omitempty"`
	CardBrand       string `json:"card_brand,omitempty"`
	CardExpiryMonth int    `json:"card_expiry_month,omitempty"`
	CardExpiryYear  int    `json:"card_expiry_year,omitempty"`
}

// PaymentAttempt tracks payment attempts for analytics
type PaymentAttempt struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	PaymentID      uint      `gorm:"not null" json:"payment_id"`
	Payment        Payment   `gorm:"foreignKey:PaymentID" json:"payment"`
	AttemptNumber  int       `gorm:"not null" json:"attempt_number"`
	Status         string    `gorm:"not null" json:"status"` // success, failed, timeout
	ErrorCode      string    `json:"error_code,omitempty"`
	ErrorMessage   string    `json:"error_message,omitempty"`
	ProcessingTime float64   `json:"processing_time"` // in milliseconds
	CreatedAt      time.Time `json:"created_at"`
}
