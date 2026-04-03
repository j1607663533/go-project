import React, { useState } from 'react';
import { 
  Card, 
  Input, 
  Button, 
  Space, 
  message, 
  Typography, 
  Divider, 
  Tag, 
  Empty,
  Spin,
  Breadcrumb,
  Alert
} from 'antd';
import { 
  RocketOutlined, 
  DownloadOutlined, 
  VideoCameraOutlined,
  LinkOutlined,
  QuestionCircleOutlined,
  CheckCircleOutlined
} from '@ant-design/icons';
import { removeWatermark } from '../api/video';

const { Title, Text, Paragraph } = Typography;

const VideoTools = () => {
  const [url, setUrl] = useState('');
  const [loading, setLoading] = useState(false);
  const [videoInfo, setVideoInfo] = useState(null);

  const handleParse = async () => {
    if (!url.trim()) {
      message.warning('请输入视频分享链接');
      return;
    }

    setLoading(true);
    setVideoInfo(null);
    try {
      const res = await removeWatermark(url);
      // 根据 utils/request.js 的封装，res 已经是 res.data
      setVideoInfo(res);
      message.success('解析成功');
    } catch (error) {
      console.error(error);
      message.error(error.message || '解析失败，请检查链接是否正确');
    } finally {
      setLoading(false);
    }
  };

  const handleDownload = () => {
    if (videoInfo && videoInfo.video_url) {
      window.open(videoInfo.video_url, '_blank');
    }
  };

  const handleCopyLink = () => {
    if (videoInfo && videoInfo.video_url) {
      navigator.clipboard.writeText(videoInfo.video_url);
      message.success('视频直链已复制到剪贴板');
    }
  };

  return (
    <div style={{ padding: '24px', maxWidth: '1000px', margin: '0 auto' }}>
      <Breadcrumb style={{ marginBottom: '16px' }}>
        <Breadcrumb.Item>首页</Breadcrumb.Item>
        <Breadcrumb.Item>视频工具</Breadcrumb.Item>
        <Breadcrumb.Item>短视频去水印</Breadcrumb.Item>
      </Breadcrumb>

      <div style={{ textAlign: 'center', marginBottom: '40px' }}>
        <Title level={2}>
          <VideoCameraOutlined /> 全平台短视频去水印
        </Title>
        <Text type="secondary">
          支持抖音、TikTok 等平台（更多平台持续更新中...）
        </Text>
      </div>

      <Card 
        className="glass-morphism"
        style={{ 
          borderRadius: '16px', 
          boxShadow: '0 8px 32px rgba(0,0,0,0.05)',
          marginBottom: '32px'
        }}
      >
        <Space direction="vertical" style={{ width: '100%' }} size="large">
          <Alert
            message="使用提示"
            description="直接粘贴分享链接即可，程序会自动提取其中的 URL 部分。"
            type="info"
            showIcon
            icon={<QuestionCircleOutlined />}
            style={{ borderRadius: '12px' }}
          />
          
          <div style={{ display: 'flex', gap: '12px' }}>
            <Input 
              size="large"
              placeholder="请粘贴短视频分享链接..." 
              prefix={<LinkOutlined style={{ color: '#bfbfbf' }} />}
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              onPressEnter={handleParse}
              style={{ borderRadius: '10px' }}
            />
            <Button 
              type="primary" 
              size="large" 
              icon={<RocketOutlined />}
              loading={loading}
              onClick={handleParse}
              style={{ 
                borderRadius: '10px',
                background: 'linear-gradient(135deg, #6366f1 0%, #a855f7 100%)',
                border: 'none',
                height: '40px'
              }}
            >
              立即解析
            </Button>
          </div>
        </Space>
      </Card>

      {loading && (
        <div style={{ textAlign: 'center', padding: '100px 0' }}>
          <Spin size="large" tip="正在极速解析中..." />
        </div>
      )}

      {videoInfo && (
        <Card
          className="glass-morphism"
          style={{ 
            borderRadius: '24px', 
            overflow: 'hidden',
            border: '1px solid rgba(255,255,255,0.2)'
          }}
        >
          <div style={{ display: 'flex', flexDirection: window.innerWidth < 768 ? 'column' : 'row', gap: '24px' }}>
            {/* 封面图展示 */}
            <div style={{ flex: '0 0 300px' }}>
               <div style={{ 
                 position: 'relative',
                 borderRadius: '16px',
                 overflow: 'hidden',
                 boxShadow: '0 4px 12px rgba(0,0,0,0.1)'
               }}>
                  <img 
                    src={videoInfo.cover} 
                    alt="封面" 
                    style={{ width: '100%', height: 'auto', display: 'block' }} 
                  />
                  <Tag 
                    color="blue" 
                    style={{ 
                      position: 'absolute', 
                      top: '12px', 
                      left: '12px',
                      borderRadius: '6px'
                    }}
                  >
                    {videoInfo.platform}
                  </Tag>
               </div>
            </div>

            {/* 信息展示 */}
            <div style={{ flex: 1, display: 'flex', flexDirection: 'column', justifyContent: 'space-between' }}>
              <div>
                <Title level={4} style={{ marginBottom: '16px' }}>{videoInfo.title}</Title>
                <Divider style={{ margin: '12px 0' }} />
                <Space style={{ marginBottom: '24px' }}>
                  <CheckCircleOutlined style={{ color: '#52c41a' }} />
                  <Text type="secondary">已成功移除画质水印</Text>
                </Space>
              </div>

              <Space size="middle" wrap>
                <Button 
                  type="primary" 
                  size="large" 
                  icon={<DownloadOutlined />}
                  onClick={handleDownload}
                  style={{ borderRadius: '10px' }}
                >
                  保存到本地
                </Button>
                <Button 
                  size="large" 
                  icon={<LinkOutlined />}
                  onClick={handleCopyLink}
                  style={{ borderRadius: '10px' }}
                >
                  复制视频直链
                </Button>
              </Space>
            </div>
          </div>
        </Card>
      )}

      {!loading && !videoInfo && (
        <div style={{ padding: '60px 0' }}>
          <Empty description="暂无解析结果，快去试试吧！" />
        </div>
      )}

      <div style={{ marginTop: '60px', opacity: 0.6 }}>
        <Divider plain>常见问题</Divider>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))', gap: '24px' }}>
           <div>
             <Text strong>Q: 解析失败怎么办？</Text>
             <Paragraph style={{ marginTop: '8px' }}>
               请确认链接是否完整，且该视频在原平台上是否可以正常播放。如果链接已过期，请重新复制。
             </Paragraph>
           </div>
           <div>
             <Text strong>Q: 下载的视频还有水印？</Text>
             <Paragraph style={{ marginTop: '8px' }}>
               部分视频可能由创作者手动添加了不可去除的水印。本工具主要去除平台自动生成的浮动水印。
             </Paragraph>
           </div>
        </div>
      </div>
    </div>
  );
};

export default VideoTools;
