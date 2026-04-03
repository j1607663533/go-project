package models

import "time"

// LotteryPrize 抽奖奖品表
type LotteryPrize struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ActivityID uint      `gorm:"not null;index" json:"activity_id"`
	Name       string    `gorm:"size:255;not null" json:"name"`
	Type       int       `gorm:"not null;default:1" json:"type"` // 奖品类型：1-实物，2-积分，3-谢谢惠顾
	ImageUrl   string    `gorm:"size:255" json:"image_url"`
	TotalStock int       `gorm:"not null;default:0" json:"total_stock"` // 初始库存
	LeftStock  int       `gorm:"not null;default:0" json:"left_stock"`  // 剩余库存
	Weight     int       `gorm:"not null;default:0" json:"weight"`      // 中奖权重
	Sort       int       `gorm:"not null;default:0" json:"sort"`        // 转盘排序
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
