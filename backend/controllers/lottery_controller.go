package controllers

import (
	"net/http"
	"strconv"

	"gin-backend/config"
	"gin-backend/services"

	"github.com/gin-gonic/gin"
)

type LotteryController struct {
	lotteryService services.LotteryService
}

func NewLotteryController(lotteryService services.LotteryService) *LotteryController {
	return &LotteryController{lotteryService: lotteryService}
}

// GetInfo 获取当前抽奖活动信息
func (c *LotteryController) GetInfo(ctx *gin.Context) {
	userId, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权"})
		return
	}

	info, err := c.lotteryService.GetLotteryInfo(userId.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": info})
}

// Draw 进行抽奖
func (c *LotteryController) Draw(ctx *gin.Context) {
	userId, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权"})
		return
	}

	prize, err := c.lotteryService.Draw(userId.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": map[string]interface{}{
			"prize_id": prize.ID,
			"name":     prize.Name,
			"type":     prize.Type,
		},
	})
}

// GetRecords 获取用户的抽奖记录
func (c *LotteryController) GetRecords(ctx *gin.Context) {
	userId, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权"})
		return
	}

	activityIdStr := ctx.Query("activity_id")
	var activityId uint
	if activityIdStr != "" {
		if id, err := strconv.Atoi(activityIdStr); err == nil {
			activityId = uint(id)
		}
	}

	records, err := c.lotteryService.GetUserRecords(userId.(uint), activityId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": records})
}

// GetPublicRecords 获取公开的中奖列表大盘（所有人的真实中奖记录）
func (c *LotteryController) GetPublicRecords(ctx *gin.Context) {
	activityIdStr := ctx.Query("activity_id")
	var activityId uint
	if activityIdStr != "" {
		if id, err := strconv.Atoi(activityIdStr); err == nil {
			activityId = uint(id)
		}
	}

	if activityId == 0 {
		ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": []interface{}{}})
		return
	}

	type PublicRecord struct {
		PrizeName string `json:"prize_name"`
		Username  string `json:"username"`
		CreatedAt string `json:"created_at"`
	}
	var records []PublicRecord

	// 直接从底层查询构建公开展示的内容
	config.DB.Table("lottery_records r").
		Select("r.prize_name, u.username, r.created_at").
		Joins("left join users u on u.id = r.user_id").
		Where("r.activity_id = ? AND r.is_hit = ?", activityId, true).
		Order("r.id desc").
		Limit(50).
		Scan(&records)

	// 脱敏用户名
	for i, r := range records {
		if len(r.Username) > 2 {
			runes := []rune(r.Username)
			records[i].Username = string(runes[0]) + "***" + string(runes[len(runes)-1])
		} else if len(r.Username) > 0 {
			records[i].Username = string([]rune(r.Username)[0]) + "***"
		} else {
			records[i].Username = "匿名用户"
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": records})
}
