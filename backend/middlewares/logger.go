package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger 请求日志中间件
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 读取请求体
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			// 重新设置请求体，因为读取后会被消耗
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// 打印请求信息
		log.Printf("\n========== 请求开始 ==========")
		log.Printf("时间: %s", startTime.Format("2006-01-02 15:04:05"))
		log.Printf("方法: %s", c.Request.Method)
		log.Printf("路径: %s", c.Request.URL.Path)
		log.Printf("完整URL: %s", c.Request.URL.String())
		log.Printf("客户端IP: %s", c.ClientIP())

		// 打印查询参数
		if len(c.Request.URL.Query()) > 0 {
			log.Printf("查询参数: %v", c.Request.URL.Query())
		}

		// 打印请求头
		log.Printf("请求头:")
		for key, values := range c.Request.Header {
			// 过滤敏感信息
			if key == "Authorization" {
				log.Printf("  %s: [已隐藏]", key)
			} else {
				log.Printf("  %s: %v", key, values)
			}
		}

		// 打印请求体
		if len(bodyBytes) > 0 {
			// 尝试格式化 JSON
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, bodyBytes, "", "  "); err == nil {
				log.Printf("请求体 (JSON):\n%s", prettyJSON.String())
			} else {
				log.Printf("请求体 (原始): %s", string(bodyBytes))
			}
		}

		log.Printf("==============================\n")

		// 处理请求
		c.Next()

		// 计算耗时
		duration := time.Since(startTime)

		// 打印响应信息
		log.Printf("\n========== 请求结束 ==========")
		log.Printf("路径: %s", c.Request.URL.Path)
		log.Printf("状态码: %d", c.Writer.Status())
		log.Printf("耗时: %v", duration)
		log.Printf("==============================\n")
	}
}
