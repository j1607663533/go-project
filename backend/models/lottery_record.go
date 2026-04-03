package models

import "time"

// LotteryRecord 抽奖流水表
type LotteryRecord struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	ActivityID uint      `gorm:"not null;index" json:"activity_id"`
	PrizeID    uint      `gorm:"not null" json:"prize_id"`
	PrizeName  string    `gorm:"size:255;not null" json:"prize_name"` // 奖品名称快照
	IsHit      bool      `gorm:"not null;default:false" json:"is_hit"`// 是否真正中奖
	CreatedAt  time.Time `json:"created_at"`
}
