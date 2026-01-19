package models

import "time"

// Order 订单模型
type Order struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	ProductID uint      `json:"product_id" gorm:"not null"`
	Quantity  uint      `json:"quantity" gorm:"not null"`
	Total     float64   `json:"total" gorm:"not null"`
	Status    string    `json:"status" gorm:"not null"`
	PaymentID uint      `json:"payment_id" gorm:"not null"`
	Payment   Payment   `json:"payment" gorm:"foreignKey:PaymentID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderCreateRequest struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Quantity  uint    `json:"quantity" binding:"required"`
	Total     float64 `json:"total"`
	Status    string  `json:"status" binding:"required"`
	PaymentID uint    `json:"payment_id"`
}

// OrderQuery 订单查询请求
type OrderQuery struct {
	PageRequest
	ProductID string `form:"product_id"`
	Status    string `form:"status"`
}

// OrderUpdateRequest 更新订单请求
type OrderUpdateRequest struct {
	Quantity  uint    `json:"quantity" binding:"omitempty" validate:"omitempty"`
	Total     float64 `json:"total" binding:"omitempty" validate:"omitempty"`
	Status    string  `json:"status" binding:"omitempty" validate:"omitempty"`
	PaymentID uint    `json:"payment_id" binding:"omitempty" validate:"omitempty"`
}

// OrderRequest 订单请求
type OrderRequest struct {
	OrderID   uint    `json:"order_id" binding:"required" validate:"required"`
	UserID    uint    `json:"user_id" binding:"required" validate:"required"`
	ProductID uint    `json:"product_id" binding:"required" validate:"required"`
	Quantity  uint    `json:"quantity" binding:"omitempty" validate:"omitempty"`
	Total     float64 `json:"total" binding:"omitempty" validate:"omitempty"`
	Status    string  `json:"status" binding:"omitempty" validate:"omitempty"`
	PaymentID uint    `json:"payment_id" binding:"omitempty" validate:"omitempty"`
}

// OrderResponse 订单响应
type OrderResponse struct {
	Order Order `json:"order"`
}

// ToResponse 转换为响应格式
func (o *Order) ToResponse() OrderResponse {
	return OrderResponse{
		Order: *o,
	}
}
