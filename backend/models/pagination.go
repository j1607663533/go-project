package models

// PageRequest 分页请求
type PageRequest struct {
	Page     int `json:"page" form:"page" binding:"omitempty,min=1"`           // 页码，从1开始
	PageSize int `json:"page_size" form:"page_size" binding:"omitempty,min=1"` // 每页数量
}

// GetPage 获取页码，默认为1
func (p *PageRequest) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

// GetPageSize 获取每页数量，默认为10
func (p *PageRequest) GetPageSize() int {
	if p.PageSize <= 0 {
		return 10
	}
	if p.PageSize > 100 {
		return 100 // 最大100条
	}
	return p.PageSize
}

// GetOffset 计算偏移量
func (p *PageRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}

// PageResponse 分页响应
type PageResponse struct {
	Page       int         `json:"page"`        // 当前页码
	PageSize   int         `json:"page_size"`   // 每页数量
	Total      int64       `json:"total"`       // 总记录数
	TotalPages int         `json:"total_pages"` // 总页数
	Data       interface{} `json:"data"`        // 数据列表
}

// NewPageResponse 创建分页响应
func NewPageResponse(page, pageSize int, total int64, data interface{}) *PageResponse {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &PageResponse{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		Data:       data,
	}
}
