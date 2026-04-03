package controllers

import (
	"net/http"
	"time"

	"gin-backend/config"
	"gin-backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AdminGetRecords 获取抽奖记录（支持分页和筛选）
func AdminGetRecords(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	title := ctx.Query("title")
	statusStr := ctx.Query("status")

	offset := (page - 1) * pageSize

	type RecordResult struct {
		models.LotteryRecord
		ActivityTitle  string `json:"activity_title"`
		ActivityStatus int    `json:"activity_status"`
		Username       string `json:"username"`
	}

	query := config.DB.Table("lottery_records r").
		Select("r.*, a.title as activity_title, a.status as activity_status, u.username").
		Joins("left join lottery_activities a on a.id = r.activity_id").
		Joins("left join users u on u.id = r.user_id")

	if title != "" {
		query = query.Where("a.title LIKE ?", "%"+title+"%")
	}
	if statusStr != "" {
		query = query.Where("a.status = ?", statusStr)
	}

	var total int64
	query.Count(&total)

	var records []RecordResult
	query.Order("r.id desc").Offset(offset).Limit(pageSize).Find(&records)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total":   total,
			"records": records,
		},
	})
}

// AdminGetActivities 获取活动的列表
func AdminGetActivities(ctx *gin.Context) {
	var activities []models.LotteryActivity
	config.DB.Order("id desc").Find(&activities)

	// 获取所有的奖品并分类组装，方便前端展示
	var allPrizes []models.LotteryPrize
	config.DB.Find(&allPrizes)

	prizeMap := make(map[uint][]models.LotteryPrize)
	for _, p := range allPrizes {
		prizeMap[p.ActivityID] = append(prizeMap[p.ActivityID], p)
	}

	result := make([]map[string]interface{}, 0)
	for _, a := range activities {
		result = append(result, map[string]interface{}{
			"activity": a,
			"prizes":   prizeMap[a.ID],
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// AdminToggleStatus 切换启动禁用状态
func AdminToggleStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	var req struct {
		Status int `json:"status"` // 1开启 0关闭
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var activity models.LotteryActivity
	if err := config.DB.First(&activity, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "未找到活动"})
		return
	}

	tx := config.DB.Begin()
	// 如果设置的状态是启用 (1)，先将其他所有的活动都改为禁用 (0)
	if req.Status == 1 {
		tx.Model(&models.LotteryActivity{}).Where("id != ?", id).Update("status", 0)
	}
	
	// 更新当前状态
	tx.Model(&activity).Update("status", req.Status)
	tx.Commit()

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "状态更新成功"})
}

// AdminSaveConfig 保存活动和奖品配置 (新建/编辑)
func AdminSaveConfig(ctx *gin.Context) {
	var req struct {
		ID         uint                  `json:"id"`
		Title      string                `json:"title"`
		StartTime  string                `json:"start_time"`
		EndTime    string                `json:"end_time"`
		DailyLimit int                   `json:"daily_limit"`
		Status     int                   `json:"status"`
		Prizes     []models.LotteryPrize `json:"prizes"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	startTime, _ := time.ParseInLocation(time.RFC3339, req.StartTime, time.Local)
	endTime, _ := time.ParseInLocation(time.RFC3339, req.EndTime, time.Local)

	tx := config.DB.Begin()

	var activity models.LotteryActivity
	if req.ID > 0 {
		tx.First(&activity, req.ID)
	}

	activity.Title = req.Title
	activity.StartTime = startTime
	activity.EndTime = endTime
	activity.DailyLimit = req.DailyLimit
	activity.Status = req.Status
	
	if req.Status == 1 {
		// 保证只有一个处于启用状态
		tx.Model(&models.LotteryActivity{}).Where("id != ?", req.ID).Update("status", 0)
	}

	if activity.ID == 0 {
		tx.Create(&activity)
	} else {
		tx.Save(&activity)
	}

	if req.ID > 0 {
		// 清理旧的奖品
		tx.Where("activity_id = ?", activity.ID).Delete(&models.LotteryPrize{})
	}

	// 插入新奖品
	for _, p := range req.Prizes {
		newPrize := p
		newPrize.ID = 0
		newPrize.ActivityID = activity.ID
		newPrize.LeftStock = newPrize.TotalStock
		tx.Create(&newPrize)
	}

	tx.Commit()

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "保存成功"})
}

// AdminDeleteActivity 删除活动
func AdminDeleteActivity(ctx *gin.Context) {
	id := ctx.Param("id")
	
	tx := config.DB.Begin()
	// 删除活动相关奖品
	tx.Where("activity_id = ?", id).Delete(&models.LotteryPrize{})
	// 删除活动
	tx.Where("id = ?", id).Delete(&models.LotteryActivity{})
	tx.Commit()

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}
