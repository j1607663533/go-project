package services

import (
	"gin-backend/models"
	"gin-backend/repositories"
)

// OrderService 订单业务逻辑接口
type OrderService interface {
	CreateOrder(order *models.Order) error
	GetOrderList() ([]models.Order, error)
	GetOrderListWithPage(query *models.OrderQuery) (*models.PageResponse, error)
	GetOrderById(id uint) (*models.Order, error)
	UpdateOrder(order *models.Order) error
	DeleteOrder(id uint) error
}

type orderService struct {
	orderRepo repositories.OrderRepository
}

// NewOrderService 创建订单业务逻辑实例
func NewOrderService(orderRepo repositories.OrderRepository) OrderService {
	return &orderService{
		orderRepo: orderRepo,
	}
}

// CreateOrder 创建订单
func (s *orderService) CreateOrder(order *models.Order) error {
	return s.orderRepo.CreateOrder(order)
}

// GetOrderList 获取订单列表
func (s *orderService) GetOrderList() ([]models.Order, error) {
	return s.orderRepo.GetOrderList()
}

// GetOrderListWithPage 分页获取订单列表
func (s *orderService) GetOrderListWithPage(query *models.OrderQuery) (*models.PageResponse, error) {
	page := query.GetPage()
	pageSize := query.GetPageSize()

	// 调用仓储层分页查询
	orders, total, err := s.orderRepo.GetOrderListWithPage(query)
	if err != nil {
		return nil, err
	}

	// 构建分页响应
	return models.NewPageResponse(page, pageSize, total, orders), nil
}

// GetOrderById 根据ID获取订单
func (s *orderService) GetOrderById(id uint) (*models.Order, error) {
	return s.orderRepo.GetOrderById(id)
}

// UpdateOrder 更新订单
func (s *orderService) UpdateOrder(order *models.Order) error {
	return s.orderRepo.UpdateOrder(order)
}

// DeleteOrder 删除订单
func (s *orderService) DeleteOrder(id uint) error {
	return s.orderRepo.DeleteOrder(id)
}
