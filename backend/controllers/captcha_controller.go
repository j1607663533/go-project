package controllers

import (
	"bytes"
	"encoding/base64"
	"gin-backend/models"
	"gin-backend/utils"
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

// CaptchaController 验证码控制器
type CaptchaController struct{}

// NewCaptchaController 创建验证码控制器实例
func NewCaptchaController() *CaptchaController {
	return &CaptchaController{}
}

// GetCaptcha 获取验证码（返回 base64 图片）
// 支持刷新：GET /api/v1/captcha?refresh=captcha_id
func (ctrl *CaptchaController) GetCaptcha(c *gin.Context) {
	var captchaID string

	// 检查是否是刷新请求
	refreshID := c.Query("refresh")
	if refreshID != "" {
		// 刷新现有验证码
		if utils.ReloadCaptcha(refreshID) {
			captchaID = refreshID
		} else {
			// 如果刷新失败（验证码已过期），生成新的
			captchaID = utils.GenerateCaptcha()
		}
	} else {
		// 生成新的验证码
		captchaID = utils.GenerateCaptcha()
	}

	// 创建一个 buffer 来存储图片
	var buf bytes.Buffer

	// 生成验证码图片
	err := captcha.WriteImage(&buf, captchaID, captcha.StdWidth, captcha.StdHeight)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成验证码失败",
		})
		return
	}

	// 将图片转换为 base64
	base64Image := base64.StdEncoding.EncodeToString(buf.Bytes())

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"captcha_id":    captchaID,
			"captcha_image": "data:image/png;base64," + base64Image,
		},
	})
}

// ServeCaptcha 提供验证码图片
func (ctrl *CaptchaController) ServeCaptcha(c *gin.Context) {
	captchaID := c.Param("id")

	// 移除 .png 后缀
	if len(captchaID) > 4 && captchaID[len(captchaID)-4:] == ".png" {
		captchaID = captchaID[:len(captchaID)-4]
	}

	// 设置响应头
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Header("Content-Type", "image/png")

	// 输出验证码图片
	err := captcha.WriteImage(c.Writer, captchaID, captcha.StdWidth, captcha.StdHeight)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "验证码不存在或已过期",
		})
		return
	}
}

// RefreshCaptcha 刷新验证码
func (ctrl *CaptchaController) RefreshCaptcha(c *gin.Context) {
	captchaID := c.Param("id")

	// 重新加载验证码
	if !utils.ReloadCaptcha(captchaID) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "验证码不存在或已过期",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "验证码已刷新",
		"data": models.CaptchaResponse{
			CaptchaID:  captchaID,
			CaptchaURL: "/api/v1/captcha/" + captchaID + ".png",
		},
	})
}

// VerifyCaptcha 验证验证码（用于测试）
func (ctrl *CaptchaController) VerifyCaptcha(c *gin.Context) {
	var req struct {
		CaptchaID string `json:"captcha_id" binding:"required"`
		Captcha   string `json:"captcha" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"errors":  utils.FormatValidationErrors(err),
		})
		return
	}

	isValid := utils.VerifyCaptcha(req.CaptchaID, req.Captcha)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"valid": isValid,
		},
	})
}
