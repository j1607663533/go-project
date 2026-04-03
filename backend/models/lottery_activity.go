package models

import "time"

// LotteryActivity 抽奖活动表
type LotteryActivity struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title      string    `gorm:"size:255;not null" json:"title"`
	StartTime  time.Time `gorm:"not null" json:"start_time"`
	EndTime    time.Time `gorm:"not null" json:"end_time"`
	DailyLimit int       `gorm:"not null;default:1" json:"daily_limit"` // 每日抽奖次数限制
	Status     int       `gorm:"not null;default:0" json:"status"`      // 状态：0-停用，1-开启
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
