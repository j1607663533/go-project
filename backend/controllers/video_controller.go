package controllers

import (
	"encoding/json"
	"fmt"
	"gin-backend/utils"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type VideoController struct{}

func NewVideoController() *VideoController {
	return &VideoController{}
}

type RemoveWatermarkRequest struct {
	URL string `json:"url" binding:"required"`
}

type VideoInfo struct {
	Title    string `json:"title"`
	Cover    string `json:"cover"`
	VideoURL string `json:"video_url"`
	Platform string `json:"platform"`
}

// RemoveWatermark 去水印接口
func (ctrl *VideoController) RemoveWatermark(c *gin.Context) {
	var req RemoveWatermarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "解析请求参数失败")
		return
	}

	url := req.URL
	var info *VideoInfo
	var err error

	// 简单的平台识别和解析逻辑
	if strings.Contains(url, "douyin.com") || strings.Contains(url, "iesdouyin.com") {
		info, err = ctrl.parseDouyin(url)
	} else {
		utils.ErrorResponse(c, http.StatusBadRequest, "暂不支持该平台或链接格式错误")
		return
	}

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("解析失败: %v", err))
		return
	}

	utils.SuccessResponse(c, info)
}

// parseDouyin 抖音解析逻辑
func (ctrl *VideoController) parseDouyin(url string) (*VideoInfo, error) {
	// 使用带 Cookie 管理的 Client
	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}

	// 1. 提取 URL
	re := regexp.MustCompile(`https?://[^\s]+`)
	matchedURL := re.FindString(url)
	if matchedURL == "" {
		return nil, fmt.Errorf("未在输入中找到有效链接")
	}

	// 2. 获取真实 URL (处理短链接重定向)
	req, _ := http.NewRequest("GET", matchedURL, nil)
	// 使用 iPhone User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1")
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("访问视频链接失败: %v", err)
	}
	defer resp.Body.Close()

	realURL := resp.Request.URL.String()
	fmt.Printf("解析到的真实 URL: %s\n", realURL)

	// 提取视频 ID
	var videoID string
	if strings.Contains(realURL, "/video/") {
		idRe := regexp.MustCompile(`/video/(\d+)`)
		matches := idRe.FindStringSubmatch(realURL)
		if len(matches) >= 2 {
			videoID = matches[1]
		}
	} else if strings.Contains(realURL, "modal_id=") {
		idRe := regexp.MustCompile(`modal_id=(\d+)`)
		matches := idRe.FindStringSubmatch(realURL)
		if len(matches) >= 2 {
			videoID = matches[1]
		}
	}

	if videoID == "" {
		return nil, fmt.Errorf("无法识别视频 ID，提取到的 URL 是: %s", realURL)
	}
	fmt.Printf("提取到视频 ID: %s\n", videoID)

	// 3. 尝试多个解析接口
	var info *VideoInfo
	
	// 策略 A: 官方 Web 接口 (通常需要 aid, aid=6383 是 Web 端的)
	apiURL := fmt.Sprintf("https://www.douyin.com/aweme/v1/web/aweme/detail/?aweme_id=%s&device_platform=webapp&aid=6383", videoID)
	info, err = ctrl.fetchVideoInfo(client, apiURL, "aweme_detail", "webapp")
	
	if err != nil || info == nil {
		fmt.Printf("策略 A 失败: %v，尝试策略 B...\n", err)
		// 策略 B: 移动端接口
		apiURL = fmt.Sprintf("https://www.iesdouyin.com/web/api/v2/aweme/iteminfo/?item_ids=%s", videoID)
		info, err = ctrl.fetchVideoInfo(client, apiURL, "item_list", "mobile")
	}

	if err != nil || info == nil {
		return nil, fmt.Errorf("解析失败: %v", err)
	}

	return info, nil
}

// fetchVideoInfo 通用抓取逻辑
func (ctrl *VideoController) fetchVideoInfo(client *http.Client, apiURL string, key string, mode string) (*VideoInfo, error) {
	req, _ := http.NewRequest("GET", apiURL, nil)
	
	if mode == "webapp" {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
		req.Header.Set("Referer", "https://www.douyin.com/")
	} else {
		req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if len(body) == 0 {
		return nil, fmt.Errorf("响应体为空 (Status: %d)", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %v", err)
	}

	var aweme map[string]interface{}
	if key == "aweme_detail" {
		if detail, ok := result["aweme_detail"].(map[string]interface{}); ok {
			aweme = detail
		}
	} else if key == "item_list" {
		if list, ok := result["item_list"].([]interface{}); ok && len(list) > 0 {
			aweme = list[0].(map[string]interface{})
		}
	}

	if aweme == nil {
		return nil, fmt.Errorf("返回数据中未找到视频详情")
	}

	video := aweme["video"].(map[string]interface{})
	desc := aweme["desc"].(string)
	
	// 提取无水印视频地址
	playAddr := video["play_addr"].(map[string]interface{})
	urlList := playAddr["url_list"].([]interface{})
	if len(urlList) == 0 {
		return nil, fmt.Errorf("未找到播放链接")
	}

	videoURL := urlList[0].(string)
	// 核心去水印逻辑：替换 playwm 为 play
	videoURL = strings.Replace(videoURL, "playwm", "play", -1)
	if strings.HasPrefix(videoURL, "//") {
		videoURL = "https:" + videoURL
	}

	// 提取封面
	cover := ""
	if originCover, ok := video["origin_cover"].(map[string]interface{}); ok {
		if cList, ok := originCover["url_list"].([]interface{}); ok && len(cList) > 0 {
			cover = cList[0].(string)
		}
	}

	return &VideoInfo{
		Title:    desc,
		Cover:    cover,
		VideoURL: videoURL,
		Platform: "抖音",
	}, nil
}
