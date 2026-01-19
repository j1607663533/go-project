package utils

import (
	"fmt"
	"time"

	"github.com/dchest/captcha"
)

const (
	// CaptchaPrefix Redis 验证码键前缀
	CaptchaPrefix = "captcha:"
	// CaptchaExpiration 验证码过期时间
	CaptchaExpiration = 10 * time.Minute
)

// GenerateCaptchaRedis 生成验证码并存储到 Redis
func GenerateCaptchaRedis() (string, string, error) {
	// 生成验证码 ID 和答案
	id := captcha.New()

	// 获取验证码答案
	digits := captcha.RandomDigits(captcha.DefaultLen)
	answer := ""
	for _, d := range digits {
		answer += fmt.Sprintf("%d", d)
	}

	// 存储到 Redis
	key := CaptchaPrefix + id
	err := CacheSetString(key, answer, CaptchaExpiration)
	if err != nil {
		return "", "", err
	}

	return id, answer, nil
}

// VerifyCaptchaRedis 验证 Redis 中的验证码
func VerifyCaptchaRedis(id, answer string) bool {
	key := CaptchaPrefix + id

	// 从 Redis 获取答案
	correctAnswer, err := CacheGetString(key)
	if err != nil {
		return false
	}

	// 验证成功后删除（一次性使用）
	if correctAnswer == answer {
		CacheDel(key)
		return true
	}

	return false
}

// RefreshCaptchaRedis 刷新验证码
func RefreshCaptchaRedis(id string) (string, error) {
	key := CaptchaPrefix + id

	// 检查验证码是否存在
	exists, err := CacheExists(key)
	if err != nil || !exists {
		return "", fmt.Errorf("验证码不存在或已过期")
	}

	// 生成新的答案
	digits := captcha.RandomDigits(captcha.DefaultLen)
	answer := ""
	for _, d := range digits {
		answer += fmt.Sprintf("%d", d)
	}

	// 更新 Redis
	err = CacheSetString(key, answer, CaptchaExpiration)
	if err != nil {
		return "", err
	}

	return answer, nil
}

// 保留原有的内存存储方式作为备用
var captchaStore = captcha.NewMemoryStore(captcha.CollectNum, captcha.Expiration)

func init() {
	captcha.SetCustomStore(captchaStore)
}

// GenerateCaptcha 生成验证码（兼容旧版本）
func GenerateCaptcha() string {
	return captcha.New()
}

// GenerateCaptchaWithLength 生成指定长度的验证码
func GenerateCaptchaWithLength(length int) string {
	return captcha.NewLen(length)
}

// VerifyCaptcha 验证验证码（兼容旧版本）
func VerifyCaptcha(id, answer string) bool {
	// 优先使用 Redis
	if VerifyCaptchaRedis(id, answer) {
		return true
	}
	// 降级到内存存储
	return captcha.VerifyString(id, answer)
}

// ReloadCaptcha 重新加载验证码
func ReloadCaptcha(id string) bool {
	return captcha.Reload(id)
}
