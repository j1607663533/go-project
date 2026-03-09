package services

import (
	"fmt"
	"gin-backend/config"
	"net/http"
	"sync"
	"time"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/officialaccount"
	"github.com/silenceper/wechat/v2/officialaccount/basic"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/silenceper/wechat/v2/officialaccount/server"
)

// WechatStatus 微信登录状态
type WechatStatus string

const (
	StatusPending  WechatStatus = "PENDING"  // 等待扫码
	StatusScanning WechatStatus = "SCANNING" // 已扫码，等待确认
	StatusSuccess  WechatStatus = "SUCCESS"  // 登录成功
	StatusExpired  WechatStatus = "EXPIRED"  // 二维码过期
)

// WechatSession 微信登录会话
type WechatSession struct {
	SceneID  string       `json:"scene_id"`
	Status   WechatStatus `json:"status"`
	UserID   uint         `json:"user_id,omitempty"`
	OpenID   string       `json:"openid,omitempty"`
	ExpireAt time.Time    `json:"expire_at"`
}

// WechatService 微信业务逻辑接口
type WechatService interface {
	GetQRCode() (*WechatSession, string, error)
	CheckStatus(sceneID string) (*WechatSession, error)
	GetServer(req *http.Request, writer http.ResponseWriter) *server.Server
	HandleCallback(msg message.MixMessage) *message.Reply
	MockScan(sceneID string, userID uint) error
}

type wechatService struct {
	sessions sync.Map
	oa       *officialaccount.OfficialAccount
}

func NewWechatService() WechatService {
	s := &wechatService{}

	if config.AppConfig.Wechat.AppID != "" {
		wc := wechat.NewWechat()
		cfg := &offConfig.Config{
			AppID:          config.AppConfig.Wechat.AppID,
			AppSecret:      config.AppConfig.Wechat.AppSecret,
			Token:          config.AppConfig.Wechat.Token,
			EncodingAESKey: config.AppConfig.Wechat.EncodingAESKey,
		}
		s.oa = wc.GetOfficialAccount(cfg)
	}

	go s.gc()
	return s
}

func (s *wechatService) GetQRCode() (*WechatSession, string, error) {
	sceneStr := fmt.Sprintf("%d", time.Now().UnixNano())

	var qrURL string
	if s.oa != nil {
		// 使用 basic 包中的 Request
		ticker, err := s.oa.GetBasic().GetQRTicket(&basic.Request{
			ExpireSeconds: 2592000,
			ActionName:    "QR_STR_SCENE",
			ActionInfo: struct {
				Scene struct {
					SceneStr string `json:"scene_str,omitempty"`
					SceneID  int    `json:"scene_id,omitempty"`
				} `json:"scene"`
			}{
				Scene: struct {
					SceneStr string `json:"scene_str,omitempty"`
					SceneID  int    `json:"scene_id,omitempty"`
				}{
					SceneStr: sceneStr,
				},
			},
		})
		if err == nil {
			qrURL = "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + ticker.Ticket
		} else {
			fmt.Printf("GetQRTicket error: %v\n", err)
		}
	}

	if qrURL == "" {
		qrURL = fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=300x300&data=LOGIN_SCENE_%s", sceneStr)
	}

	session := &WechatSession{
		SceneID:  sceneStr,
		Status:   StatusPending,
		ExpireAt: time.Now().Add(5 * time.Minute),
	}
	s.sessions.Store(sceneStr, session)

	return session, qrURL, nil
}

func (s *wechatService) CheckStatus(sceneID string) (*WechatSession, error) {
	val, ok := s.sessions.Load(sceneID)
	if !ok {
		return nil, fmt.Errorf("会话不存在或已过期")
	}

	session := val.(*WechatSession)
	if time.Now().After(session.ExpireAt) {
		session.Status = StatusExpired
		return session, nil
	}

	return session, nil
}

func (s *wechatService) GetServer(req *http.Request, writer http.ResponseWriter) *server.Server {
	if s.oa == nil {
		return nil
	}
	return s.oa.GetServer(req, writer)
}

func (s *wechatService) HandleCallback(msg message.MixMessage) *message.Reply {
	var sceneID string

	if msg.Event == message.EventScan {
		sceneID = msg.EventKey
	} else if msg.Event == message.EventSubscribe {
		// 关注并扫码
		if len(msg.EventKey) > 8 {
			sceneID = msg.EventKey[8:]
		}
	}

	if sceneID != "" {
		if val, ok := s.sessions.Load(sceneID); ok {
			session := val.(*WechatSession)
			session.Status = StatusSuccess
			session.OpenID = string(msg.FromUserName)
			session.UserID = 1 // 实际需查询数据库
			s.sessions.Store(sceneID, session)

			return &message.Reply{
				MsgType: message.MsgTypeText,
				MsgData: message.NewText("登录成功！"),
			}
		}
	}

	return nil
}

func (s *wechatService) MockScan(sceneID string, userID uint) error {
	val, ok := s.sessions.Load(sceneID)
	if !ok {
		return fmt.Errorf("会话不存在")
	}

	session := val.(*WechatSession)
	session.Status = StatusSuccess
	session.UserID = userID
	s.sessions.Store(sceneID, session)
	return nil
}

func (s *wechatService) gc() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		s.sessions.Range(func(key, value interface{}) bool {
			session := value.(*WechatSession)
			if time.Now().After(session.ExpireAt) {
				s.sessions.Delete(key)
			}
			return true
		})
	}
}
