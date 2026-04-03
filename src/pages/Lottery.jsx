import React, { useState, useEffect } from 'react';
import { Card, Button, Spin, message, Typography, List, Empty, Badge, Avatar, Tabs } from 'antd';
import { GiftOutlined, NotificationOutlined, UserOutlined, HistoryOutlined } from '@ant-design/icons';
import { getLotteryInfo, drawLottery, getLotteryPublicRecords, getMyLotteryRecords } from '../api/lottery';

const { Title, Text } = Typography;

const Lottery = () => {
  const [loading, setLoading] = useState(true);
  const [info, setInfo] = useState({ active: false, prizes: [], remain: 0 });
  const [publicRecords, setPublicRecords] = useState([]);
  const [myRecords, setMyRecords] = useState([]);
  const [isDrawing, setIsDrawing] = useState(false);
  const [rotateDegree, setRotateDegree] = useState(0);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      setLoading(true);
      const res = await getLotteryInfo();
      if (res) {
        setInfo(res);
        // 基于当前有效的活动 ID 获取战报与个人记录
        if (res.activity_id) {
          const publicRes = await getLotteryPublicRecords({ activity_id: res.activity_id });
          if (publicRes) setPublicRecords(publicRes);

          const myRes = await getMyLotteryRecords({ activity_id: res.activity_id });
          if (myRes) setMyRecords(myRes);
        }
      }
    } catch (error) {
      console.error('Failed to fetch lottery info:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleDraw = async () => {
    if (!info.active) {
      message.warning('活动未开启或已结束');
      return;
    }
    if (info.remain <= 0) {
      message.warning('您今日的抽奖次数已用尽');
      return;
    }
    if (isDrawing) return;

    try {
      setIsDrawing(true);
      const res = await drawLottery();
      
      if (res && res.prize_id) {
        const prizeId = res.prize_id;
        const prizeIndex = info.prizes.findIndex(p => p.id === prizeId);
        
        if (prizeIndex === -1) {
          message.error('抽奖异常');
          setIsDrawing(false);
          return;
        }

        const slicesCount = info.prizes.length;
        const sliceAngle = 360 / slicesCount;
        const targetAngle = 360 - (prizeIndex * sliceAngle + sliceAngle / 2);
        
        // 追加8圈以产生更强烈的旋转动感
        const currentRotation = rotateDegree + 360 * 8 + (targetAngle - (rotateDegree % 360));
        
        setRotateDegree(currentRotation);
        setInfo(prev => ({ ...prev, remain: prev.remain - 1 }));

        setTimeout(() => {
          if (res.type === 3) {
            message.info(`很遗憾，${res.name}！`);
          } else {
            message.success(`🎉 恭喜您中奖啦！获得了：${res.name}`);
          }
          setIsDrawing(false);
          fetchData();
        }, 5000); // 5秒转盘动画
      } else {
        message.error(res.message || '抽奖失败');
        setIsDrawing(false);
      }
    } catch (error) {
      message.error(error.message || '系统异常');
      setIsDrawing(false);
    }
  };

  if (loading) return <Spin size="large" style={{ display: 'flex', justifyContent: 'center', marginTop: '20vh' }} />;

  // 为奖品生成丰富色彩
  const generateColor = (idx, total) => {
    const colors = [
      '#ffeaa7', '#fab1a0', '#81ecec', '#74b9ff', 
      '#a29bfe', '#55efc4', '#ff7675', '#ffeaa7'
    ];
    return colors[idx % colors.length];
  };

  return (
    <div style={{ 
      padding: '40px 20px', 
      minHeight: '85vh',
      background: 'linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%)',
      borderRadius: 16
    }}>
      <div style={{ textAlign: 'center', marginBottom: 40 }}>
        <Title level={1} style={{ 
            color: '#2d3436', 
            textShadow: '2px 2px 4px rgba(0,0,0,0.1)',
            fontFamily: "'Fredoka One', 'Microsoft YaHei', cursive",
            margin: 0
        }}>
          💎 幸 运 大 转 盘 💎
        </Title>
        <Text type="secondary" style={{ fontSize: 16, marginTop: 10, display: 'block' }}>
          参与试试手气，海量好礼等你来拿！
        </Text>
      </div>
      
      {info.active ? (
        <div style={{ display: 'flex', gap: 40, flexWrap: 'wrap', justifyContent: 'center' }}>
          
          {/* 转盘卡片 */}
          <Card 
            title={<><GiftOutlined /> 立即抽奖</>}
            bordered={false}
            style={{ 
              flex: '1 1 450px', 
              maxWidth: 550, 
              textAlign: 'center',
              boxShadow: '0 20px 40px rgba(0,0,0,0.08)',
              borderRadius: 24,
              overflow: 'hidden'
            }}
            headStyle={{ borderBottom: 'none', fontSize: 18, paddingTop: 20 }}
          >
            <div style={{ marginBottom: 20 }}>
              <Badge count={`今日剩余: ${info.remain} 次`} style={{ backgroundColor: '#52c41a', fontSize: 16, padding: '0 12px', height: 28, lineHeight: '28px' }} />
            </div>

            {info.prizes && info.prizes.length > 0 ? (
               <div style={{ 
                 position: 'relative', 
                 width: 360, 
                 height: 360, 
                 margin: '0 auto', 
                 borderRadius: '50%', 
                 border: '12px solid rgba(255, 255, 255, 0.5)',
                 boxShadow: '0 0 30px rgba(255, 118, 117, 0.4), inset 0 0 20px rgba(0,0,0,0.1)',
                 background: '#fff',
                 padding: 5
               }}>
                 <div style={{ 
                   width: '100%', 
                   height: '100%', 
                   borderRadius: '50%',
                   overflow: 'hidden',
                   position: 'relative',
                 }}>
                   <div 
                     style={{ 
                       width: '100%', 
                       height: '100%', 
                       transition: 'transform 5s cubic-bezier(0.2, 0.8, 0.1, 1)', 
                       transform: `rotate(${rotateDegree}deg)`,
                       position: 'relative'
                     }}
                   >
                     {/* 渲染扇形背景 */}
                     {info.prizes.map((prize, idx) => {
                       const sliceAngle = 360 / info.prizes.length;
                       const transform = `rotate(${idx * sliceAngle}deg) skewY(${-(90 - sliceAngle)}deg)`;
                       const bgColor = generateColor(idx, info.prizes.length);
                       
                       return (
                         <div 
                           key={`bg-${prize.id}`} 
                           style={{
                             position: 'absolute',
                             top: 0,
                             right: 0,
                             width: '50%',
                             height: '50%',
                             transformOrigin: '0% 100%',
                             transform: transform,
                             backgroundColor: bgColor,
                             border: '1px solid rgba(255,255,255,0.7)',
                             boxShadow: 'inset 0 0 10px rgba(0,0,0,0.05)'
                           }} 
                         />
                       );
                     })}

                     {/* 渲染居中文字（独立图层，无歪曲变形） */}
                     {info.prizes.map((prize, idx) => {
                        const sliceAngle = 360 / info.prizes.length;
                        // 计算扇形中轴线角度：由于背景起始 12 点钟，而 DOM 的 0 度是 3 点钟，所以计算出的相对角度需要顺时针减去 90度
                        const centerAngle = idx * sliceAngle + sliceAngle / 2;
                        return (
                            <div
                               key={`text-${prize.id}`}
                               style={{
                                   position: 'absolute',
                                   top: '50%',
                                   left: '50%',
                                   width: 140,       // 控制字体的区域长度（转盘半径是 180）
                                   height: 40,
                                   marginTop: -20,   // Y 轴中心回归
                                   transformOrigin: '0 50%',
                                   transform: `rotate(${centerAngle - 90}deg)`, 
                                   display: 'flex',
                                   alignItems: 'center',
                                   justifyContent: 'center', // 水平居中
                                   fontSize: 14, 
                                   fontWeight: 800,
                                   color: '#34495e',
                                   textShadow: '0 1px 1px rgba(255,255,255,0.9)'
                               }}
                            >
                                <span style={{ paddingLeft: 30, textAlign: 'center', lineHeight: '1.2' }}>
                                    {prize.name}
                                </span>
                            </div>
                        )
                     })}
                   </div>
                 </div>
                 
                 {/* 抽奖指针 */}
                 <div style={{ 
                     position: 'absolute', 
                     top: '50%', 
                     left: '50%', 
                     transform: 'translate(-50%, -50%)', 
                     zIndex: 10,
                     filter: 'drop-shadow(0 4px 6px rgba(0,0,0,0.2))'
                 }}>
                   {/* 顶部红色三角形指针 */}
                   <div style={{ 
                       width: 0, 
                       height: 0, 
                       borderLeft: '16px solid transparent', 
                       borderRight: '16px solid transparent', 
                       borderBottom: '35px solid #ff4757', 
                       margin: '0 auto', 
                       marginBottom: -15,
                       position: 'relative',
                       zIndex: 11
                   }} />
                   
                   <Button 
                     type="primary" 
                     shape="circle" 
                     style={{ 
                         width: 80, 
                         height: 80, 
                         fontSize: 22,
                         fontWeight: 900, 
                         background: isDrawing || info.remain <= 0 ? '#b2bec3' : 'linear-gradient(135deg, #ff4757 0%, #ff6b81 100%)',
                         border: '4px solid #fff',
                         boxShadow: isDrawing ? 'none' : '0 4px 15px rgba(255, 71, 87, 0.4)',
                         color: '#ffffff',
                         display: 'flex',
                         alignItems: 'center',
                         justifyContent: 'center'
                     }}
                     onClick={handleDraw}
                     disabled={isDrawing || info.remain <= 0}
                   >
                     {isDrawing ? 'GO!' : '抽奖'}
                   </Button>
                 </div>
             </div>
            ) : (
              <Empty description="奖池尚未配置完成" style={{ padding: '60px 0' }} />
            )}
          </Card>

          {/* 记录卡片 */}
          <Card 
            bordered={false}
            style={{ 
                flex: '1 1 350px', 
                maxWidth: 450,
                boxShadow: '0 20px 40px rgba(0,0,0,0.05)',
                borderRadius: 24,
            }}
            bodyStyle={{ padding: '5px 20px 20px' }}
          >
           <Tabs defaultActiveKey="1" centered size="large">
            <Tabs.TabPane tab={<span><NotificationOutlined />全服动态</span>} key="1">
              {publicRecords.length > 0 ? (
                <div style={{ maxHeight: 380, overflowY: 'auto', paddingRight: 10 }}>
                    <List
                      itemLayout="horizontal"
                      dataSource={publicRecords}
                      renderItem={item => (
                        <List.Item style={{ padding: '12px 16px', borderBottom: '1px solid #f1f2f6' }}>
                          <List.Item.Meta
                            avatar={<Avatar icon={<UserOutlined />} style={{ backgroundColor: '#f56a00' }} />}
                            title={
                              <div>
                                  <Text strong style={{ marginRight: 8, color: '#2d3436' }}>{item.username}</Text>
                                  <Text type="secondary" style={{ fontSize: 12 }}>刚刚抽中了</Text>
                              </div>
                            }
                            description={<Text strong style={{ color: '#ff4757', fontSize: 14 }}>{item.prize_name}</Text>}
                          />
                        </List.Item>
                      )}
                    />
                </div>
              ) : (
                  <Empty description="全服大奖虚位以待，快来试试手气！" image={Empty.PRESENTED_IMAGE_SIMPLE} style={{ padding: '40px 0' }} />
              )}
            </Tabs.TabPane>
            
            <Tabs.TabPane tab={<span><HistoryOutlined />我的记录</span>} key="2">
              {myRecords.length > 0 ? (
                <div style={{ maxHeight: 380, overflowY: 'auto', paddingRight: 10 }}>
                    <List
                      itemLayout="horizontal"
                      dataSource={myRecords}
                      renderItem={item => (
                        <List.Item style={{ padding: '12px 16px', borderBottom: '1px solid #f1f2f6', transition: 'all 0.3s' }} className="hover-bg">
                          <List.Item.Meta
                            title={<Text strong style={{ color: item.is_hit ? '#ff4757' : '#636e72', fontSize: 15 }}>{item.prize_name}</Text>}
                            description={<Text type="secondary" style={{ fontSize: 13 }}>{new Date(item.created_at).toLocaleString()}</Text>}
                          />
                          <div>
                              {item.is_hit ? <Badge status="success" text="中大奖" /> : <Badge status="default" text="未中奖" />}
                          </div>
                        </List.Item>
                      )}
                    />
                </div>
              ) : (
                  <Empty description="您还没有在这场活动中抽奖哦" image={Empty.PRESENTED_IMAGE_SIMPLE} style={{ padding: '40px 0' }} />
              )}
            </Tabs.TabPane>
           </Tabs>
          </Card>

        </div>
      ) : (
        <Card style={{ maxWidth: 600, margin: '40px auto', borderRadius: 20, textAlign: 'center', padding: '40px 0', boxShadow: '0 10px 30px rgba(0,0,0,0.05)' }}>
          <Empty 
            description={<span style={{ fontSize: 18, color: '#636e72' }}>本期活动暂未开启哦！</span>} 
            image="https://gw.alipayobjects.com/zos/antfincdn/ZHrcdLPrvN/empty.svg"
            imageStyle={{ height: 120 }}
          >
            <Button type="primary" size="large" onClick={() => window.location.reload()} style={{ marginTop: 20 }}>
                刷新重试
            </Button>
          </Empty>
        </Card>
      )}
    </div>
  );
};

export default Lottery;
