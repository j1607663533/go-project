package models

import "time"

// Article 文章模型
type Article struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"size:200;not null"`
	Content   string    `json:"content" gorm:"type:text"`
	Tags      []Tag     `json:"tags" gorm:"many2many:article_tags;"` // 多对多关系
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Tag 标签模型
type Tag struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:50;unique;not null"`
	Articles  []Article `json:"articles" gorm:"many2many:article_tags;"` // 反向多对多关系
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ArticleCreateRequest 创建文章请求
type ArticleCreateRequest struct {
	Title   string `json:"title" binding:"required,max=200"`
	Content string `json:"content" binding:"required"`
	TagIDs  []uint `json:"tag_ids"`
}

// ArticleResponse 文章响应
type ArticleResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []Tag     `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
}
