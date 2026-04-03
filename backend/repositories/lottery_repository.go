package repositories

import (
	"gin-backend/models"

	"gorm.io/gorm"
)

type LotteryRepository interface {
	GetActiveActivity() (*models.LotteryActivity, error)
	GetPrizesByActivityID(activityID uint) ([]models.LotteryPrize, error)
	CountUserRecordsToday(userID uint) (int64, error)
	CreateRecord(record *models.LotteryRecord) error
	DeductPrizeStock(prizeID uint) (int64, error)
	GetUserRecords(userID uint, activityID uint) ([]models.LotteryRecord, error)
	GetPublicRecords(activityID uint, limit int) ([]models.LotteryRecord, error)
	GetPrizeByID(prizeID uint) (*models.LotteryPrize, error)
}

type lotteryRepository struct {
	db *gorm.DB
}

func NewLotteryRepository(db *gorm.DB) LotteryRepository {
	return &lotteryRepository{db: db}
}

func (r *lotteryRepository) GetActiveActivity() (*models.LotteryActivity, error) {
	var activity models.LotteryActivity
	// 取第一个开启状态的有效活动
	err := r.db.Where("status = 1 AND now() >= start_time AND now() <= end_time").First(&activity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &activity, nil
}

func (r *lotteryRepository) GetPrizesByActivityID(activityID uint) ([]models.LotteryPrize, error) {
	var prizes []models.LotteryPrize
	err := r.db.Where("activity_id = ?", activityID).Order("sort asc").Find(&prizes).Error
	return prizes, err
}

func (r *lotteryRepository) CountUserRecordsToday(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.LotteryRecord{}).
		Where("user_id = ? AND DATE(created_at) = CURDATE()", userID).
		Count(&count).Error
	return count, err
}

func (r *lotteryRepository) CreateRecord(record *models.LotteryRecord) error {
	return r.db.Create(record).Error
}

func (r *lotteryRepository) DeductPrizeStock(prizeID uint) (int64, error) {
	// 使用乐观锁防止超卖：只有剩余库存大于0时才扣减
	result := r.db.Model(&models.LotteryPrize{}).
		Where("id = ? AND left_stock > 0", prizeID).
		UpdateColumn("left_stock", gorm.Expr("left_stock - 1"))
	return result.RowsAffected, result.Error
}

func (r *lotteryRepository) GetUserRecords(userID uint, activityID uint) ([]models.LotteryRecord, error) {
	var records []models.LotteryRecord
	query := r.db.Where("user_id = ?", userID)
	if activityID > 0 {
		query = query.Where("activity_id = ?", activityID)
	}
	err := query.Order("id desc").Find(&records).Error
	return records, err
}

func (r *lotteryRepository) GetPublicRecords(activityID uint, limit int) ([]models.LotteryRecord, error) {
	var records []models.LotteryRecord
	// 只返回真正中奖的流水（排除谢谢惠顾）
	query := r.db.Where("activity_id = ? AND is_hit = ?", activityID, true)
	err := query.Order("id desc").Limit(limit).Find(&records).Error
	return records, err
}

func (r *lotteryRepository) GetPrizeByID(prizeID uint) (*models.LotteryPrize, error) {
	var prize models.LotteryPrize
	err := r.db.First(&prize, prizeID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &prize, nil
}
