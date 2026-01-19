package controllers

import (
	"fmt"
	"gin-backend/models"
	"gin-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// OrderController 订单控制器
type OrderController struct {
	orderService services.OrderService
}

// NewOrderController 创建订单控制器实例
func NewOrderController(orderService services.OrderService) *OrderController {
	return &OrderController{
		orderService: orderService,
	}
}

// GetOrderList 获取订单列表
func (ctrl *OrderController) GetOrderList(c *gin.Context) {
	orders, err := ctrl.orderService.GetOrderList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取订单列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    orders,
	})
}

// GetOrderListWithPage 分页获取订单列表
func (ctrl *OrderController) GetOrderListWithPage(c *gin.Context) {
	// 绑定分页和筛选参数
	var query models.OrderQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "查询参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 调用服务层获取数据
	pageResp, err := ctrl.orderService.GetOrderListWithPage(&query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取订单列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    pageResp,
	})
}

// GetOrderById 获取订单详情
func (ctrl *OrderController) GetOrderById(c *gin.Context) {
	orderId := c.Param("id")
	id, err := strconv.Atoi(orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "获取订单详情失败",
			"error":   err.Error(),
		})
		return
	}

	order, err := ctrl.orderService.GetOrderById(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取订单详情失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    order,
	})
}

// CreateOrder 创建订单
func (ctrl *OrderController) CreateOrder(c *gin.Context) {
	// 从上下文中获取 userID (由 AuthMiddleware 设置)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录或登录已失效",
		})
		return
	}

	var req models.OrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 打印详细错误
		fmt.Printf("创建订单绑定失败: %+v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "创建订单参数验证失败",
			"error":   err.Error(),
		})
		return
	}
	fmt.Printf("收到创建订单请求: %+v\n", req)

	// 将请求转换为订单模型
	order := &models.Order{
		UserID:    userID.(uint), // 类型转换
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Total:     req.Total,
		Status:    req.Status,
		PaymentID: req.PaymentID,
	}

	err := ctrl.orderService.CreateOrder(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建订单失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    order,
	})
}

// UpdateOrder 更新订单
func (ctrl *OrderController) UpdateOrder(c *gin.Context) {
	orderId := c.Param("id")
	id, err := strconv.Atoi(orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的订单ID",
		})
		return
	}

	// 检查订单是否存在
	existingOrder, err := ctrl.orderService.GetOrderById(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "订单不存在",
		})
		return
	}

	// 从上下文中获取 userID (由 AuthMiddleware 设置)
	currentUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录或登录已失效",
		})
		return
	}

	// 权限检查：只能修改自己的订单
	if existingOrder.UserID != currentUserID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "没有权限修改此订单",
		})
		return
	}

	var req models.OrderCreateRequest // 复用创建请求结构
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("更新订单绑定失败: %+v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "更新订单参数验证失败",
			"error":   err.Error(),
		})
		return
	}
	fmt.Printf("收到更新订单请求: %+v\n", req)

	// 更新字段
	existingOrder.ProductID = req.ProductID
	existingOrder.Quantity = req.Quantity
	existingOrder.Total = req.Total
	existingOrder.Status = req.Status
	existingOrder.PaymentID = req.PaymentID

	err = ctrl.orderService.UpdateOrder(existingOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新订单失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data":    existingOrder,
	})
}

// DeleteOrder 删除订单
func (ctrl *OrderController) DeleteOrder(c *gin.Context) {
	orderId := c.Param("id")
	id, err := strconv.Atoi(orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的订单ID",
		})
		return
	}

	// 检查订单是否存在
	existingOrder, err := ctrl.orderService.GetOrderById(uint(id))

	fmt.Printf("收到删除订单请求: %+v\n", existingOrder)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "订单不存在",
		})
		return
	}

	// 从上下文中获取 userID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录或登录已失效",
		})
		return
	}

	// 权限检查
	if existingOrder.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "没有权限删除此订单",
		})
		return
	}

	err = ctrl.orderService.DeleteOrder(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除订单失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}
