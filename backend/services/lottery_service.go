package services

import (
	"errors"
	"math/rand"
	"time"

	"gin-backend/models"
	"gin-backend/repositories"
)

type LotteryService interface {
	GetLotteryInfo(userID uint) (map[string]interface{}, error)
	Draw(userID uint) (*models.LotteryPrize, error)
	GetUserRecords(userID uint, activityID uint) ([]models.LotteryRecord, error)
}

type lotteryService struct {
	repo repositories.LotteryRepository
}

func NewLotteryService(repo repositories.LotteryRepository) LotteryService {
	// 初始化随机种子
	rand.Seed(time.Now().UnixNano())
	return &lotteryService{repo: repo}
}

func (s *lotteryService) GetLotteryInfo(userID uint) (map[string]interface{}, error) {
	activity, err := s.repo.GetActiveActivity()
	if err != nil {
		return nil, err
	}

	if activity == nil {
		return map[string]interface{}{
			"active": false,
			"msg":    "当前无开放的抽奖活动",
		}, nil
	}

	prizes, err := s.repo.GetPrizesByActivityID(activity.ID)
	if err != nil {
		return nil, err
	}

	// 计算剩余次数
	var remain int64 = 0
	count, _ := s.repo.CountUserRecordsToday(userID)
	if int64(activity.DailyLimit)-count > 0 {
		remain = int64(activity.DailyLimit) - count
	}

	// 屏蔽敏感信息后返回给前端
	safePrizes := make([]map[string]interface{}, 0)
	for _, p := range prizes {
		safePrizes = append(safePrizes, map[string]interface{}{
			"id":        p.ID,
			"name":      p.Name,
			"type":      p.Type,
			"image_url": p.ImageUrl,
			"sort":      p.Sort,
		})
	}

	return map[string]interface{}{
		"activity_id": activity.ID,
		"active":      true,
		"remain": remain,
		"prizes": safePrizes,
	}, nil
}

func (s *lotteryService) Draw(userID uint) (*models.LotteryPrize, error) {
	activity, err := s.repo.GetActiveActivity()
	if err != nil {
		return nil, err
	}
	if activity == nil {
		return nil, errors.New("活动未开启或已结束")
	}

	// 1. 频控：检查剩余次数
	count, err := s.repo.CountUserRecordsToday(userID)
	if err != nil {
		return nil, err
	}
	if count >= int64(activity.DailyLimit) {
		return nil, errors.New("今日抽奖次数已用尽")
	}

	// 2. 获取有效奖品池 (过滤库存不足的实物奖品，但"谢谢惠顾"可能是无限库存)
	allPrizes, err := s.repo.GetPrizesByActivityID(activity.ID)
	if err != nil || len(allPrizes) == 0 {
		return nil, errors.New("奖品未配置")
	}

	validPrizes := make([]models.LotteryPrize, 0)
	var totalWeight int
	var fallbackPrize *models.LotteryPrize

	for _, p := range allPrizes {
		// 记录一下兜底奖品（比如谢谢惠顾，type=3），当没有抽中任何东西或全部无库存时使用
		if p.Type == 3 {
			fallbackPrize = &models.LotteryPrize{}
			*fallbackPrize = p
		}

		if p.Weight > 0 && (p.LeftStock > 0 || p.TotalStock == -1 || p.Type == 3) {
			validPrizes = append(validPrizes, p)
			totalWeight += p.Weight
		}
	}

	if totalWeight <= 0 {
		// 都没库存了，直接给谢谢惠顾
		if fallbackPrize != nil {
			s.recordDraw(userID, activity.ID, fallbackPrize)
			return fallbackPrize, nil
		}
		return nil, errors.New("奖品已抽完")
	}

	// 3. 权重随机抽取
	r := rand.Intn(totalWeight)
	var sum int
	var hitPrize *models.LotteryPrize

	for _, p := range validPrizes {
		sum += p.Weight
		if r < sum {
			hitPrize = &models.LotteryPrize{}
			*hitPrize = p
			break
		}
	}

	// 4. 尝试扣减库存
	if hitPrize != nil && hitPrize.Type != 3 && hitPrize.TotalStock != -1 {
		rows, _ := s.repo.DeductPrizeStock(hitPrize.ID)
		if rows == 0 {
			// 并发超卖导致没扣成功库存，降级为"谢谢惠顾"
			if fallbackPrize != nil {
				s.recordDraw(userID, activity.ID, fallbackPrize)
				return fallbackPrize, nil
			}
			return nil, errors.New("手慢了，奖品被抢光了")
		}
	}

	// 如果根本没抽中而且没落到任何奖（防万一），降级
	if hitPrize == nil && fallbackPrize != nil {
		hitPrize = fallbackPrize
	}

	// 5. 记录流水
	if hitPrize != nil {
		s.recordDraw(userID, activity.ID, hitPrize)
	}

	return hitPrize, nil
}

func (s *lotteryService) recordDraw(userID uint, activityID uint, prize *models.LotteryPrize) {
	record := &models.LotteryRecord{
		UserID:     userID,
		ActivityID: activityID,
		PrizeID:    prize.ID,
		PrizeName:  prize.Name,
		IsHit:      prize.Type != 3,
		CreatedAt:  time.Now(),
	}
	_ = s.repo.CreateRecord(record)
}

func (s *lotteryService) GetUserRecords(userID uint, activityID uint) ([]models.LotteryRecord, error) {
	return s.repo.GetUserRecords(userID, activityID)
}
