package repositories

import (
	"gin-backend/models"

	"gorm.io/gorm"
)

// OrderRepository 订单仓储接口
type OrderRepository interface {
	CreateOrder(order *models.Order) error
	GetOrderList() ([]models.Order, error)
	GetOrderListWithPage(query *models.OrderQuery) ([]models.Order, int64, error) // 分页查询（带筛选）
	GetOrderById(id uint) (*models.Order, error)
	UpdateOrder(order *models.Order) error
	DeleteOrder(id uint) error
}

type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository 创建订单仓储实例
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

// CreateOrder 创建订单
func (r *orderRepository) CreateOrder(order *models.Order) error {
	return r.db.Create(order).Error
}

// GetOrderList 获取订单列表
func (r *orderRepository) GetOrderList() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "nickname")
	}).Preload("Payment").Find(&orders).Error
	return orders, err
}

// GetOrderListWithPage 分页获取订单列表（支持筛选）
func (r *orderRepository) GetOrderListWithPage(query *models.OrderQuery) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	db := r.db.Model(&models.Order{})

	// 应用筛选条件
	if query.ProductID != "" {
		db = db.Where("CAST(product_id AS CHAR) LIKE ?", "%"+query.ProductID+"%")
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	// 查询总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 偏移量
	offset := query.GetOffset()
	pageSize := query.GetPageSize()

	// 分页查询
	err := db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "nickname")
	}).Preload("Payment").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&orders).Error

	return orders, total, err
}

// GetOrderById 根据ID获取订单
func (r *orderRepository) GetOrderById(id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "nickname")
	}).Preload("Payment").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// UpdateOrder 更新订单
func (r *orderRepository) UpdateOrder(order *models.Order) error {
	return r.db.Save(order).Error
}

// DeleteOrder 删除订单
func (r *orderRepository) DeleteOrder(id uint) error {
	return r.db.Delete(&models.Order{}, id).Error
}
