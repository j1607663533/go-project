package models

import "time"

// Payment 支付模型
type Payment struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	OrderID       uint      `json:"order_id" gorm:"not null"`
	Amount        float64   `json:"amount" gorm:"not null"`
	PaymentMethod string    `json:"payment_method" gorm:"not null"` // credit_card, debit_card, paypal, etc.
	Status        string    `json:"status" gorm:"not null"`         // pending, completed, failed, refunded
	TransactionID string    `json:"transaction_id" gorm:"unique"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// PaymentCreateRequest 创建支付请求
type PaymentCreateRequest struct {
	OrderID       uint    `json:"order_id" binding:"required" validate:"required"`
	Amount        float64 `json:"amount" binding:"required" validate:"required"`
	PaymentMethod string  `json:"payment_method" binding:"required" validate:"required"`
}

// PaymentUpdateRequest 更新支付请求
type PaymentUpdateRequest struct {
	Status        string `json:"status" binding:"omitempty" validate:"omitempty"`
	TransactionID string `json:"transaction_id" binding:"omitempty" validate:"omitempty"`
}

// PaymentResponse 支付响应
type PaymentResponse struct {
	Payment Payment `json:"payment"`
}

// ToResponse 转换为响应格式
func (p *Payment) ToResponse() PaymentResponse {
	return PaymentResponse{
		Payment: *p,
	}
}
