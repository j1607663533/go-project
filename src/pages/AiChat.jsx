import React, { useState, useRef, useEffect } from 'react';
import { Layout, Input, Button, Avatar, List, Spin, Typography, Card, message } from 'antd';
import { SendOutlined, UserOutlined, RobotOutlined, ClearOutlined } from '@ant-design/icons';
import { sendChatMessage } from '../api/ai';
import { saveMessage, getAllMessages, clearUserMessages } from '../utils/db';
import { getCurrentUser } from '../api/auth';

const { Content } = Layout;
const { Text } = Typography;

const AiChat = () => {
    const [messages, setMessages] = useState([]);
    const [inputValue, setInputValue] = useState('');
    const [loading, setLoading] = useState(false);
    const scrollRef = useRef(null);
    const currentUser = getCurrentUser();
    const userId = currentUser?.id || currentUser?.username || 'guest';

    // 初始化加载历史记录
    useEffect(() => {
        const loadHistory = async () => {
            try {
                const history = await getAllMessages(userId);
                if (history && history.length > 0) {
                    setMessages(history.sort((a, b) => a.id - b.id));
                } else {
                    const welcomeMsg = { 
                        id: Date.now(), 
                        role: 'ai', 
                        content: `您好 ${currentUser?.nickname || currentUser?.username || ''}！我是您的 AI 助手。这里的记录仅您可见。` 
                    };
                    setMessages([welcomeMsg]);
                    await saveMessage(userId, welcomeMsg);
                }
            } catch (error) {
                console.error('Failed to load history:', error);
            }
        };
        loadHistory();
    }, [userId]);

    // 自动滚动到最新消息
    useEffect(() => {
        if (scrollRef.current) {
            scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
        }
    }, [messages, loading]);

    const handleSend = async () => {
        if (!inputValue.trim()) return;

        const userMsg = {
            id: Date.now(),
            role: 'user',
            content: inputValue,
            timestamp: new Date().toISOString()
        };

        // 更新 UI 并持久化到 IndexedDB
        setMessages(prev => [...prev, userMsg]);
        await saveMessage(userId, userMsg);
        
        setInputValue('');
        setLoading(true);

        try {
            const res = await sendChatMessage(userMsg.content);
            const aiMsg = {
                id: Date.now() + 1,
                role: 'ai',
                content: res.reply,
                timestamp: new Date().toISOString()
            };
            
            setMessages(prev => [...prev, aiMsg]);
            await saveMessage(userId, aiMsg);
        } catch (error) {
            message.error('发送失败，请稍后重试');
        } finally {
            setLoading(false);
        }
    };

    const clearChat = async () => {
        try {
            await clearUserMessages(userId);
            const welcomeMsg = { id: Date.now(), role: 'ai', content: '会话已清空。有什么我可以帮您的吗？' };
            setMessages([welcomeMsg]);
            await saveMessage(userId, welcomeMsg);
            message.success('您的聊天记录已清除');
        } catch (error) {
            message.error('清除失败');
        }
    };

    return (
        <Layout className="chat-layout" style={{ height: 'calc(100vh - 64px)', background: 'transparent' }}>
            <Content style={{ padding: '24px', display: 'flex', flexDirection: 'column', gap: '20px' }}>
                <Card 
                    title={
                        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                            <span><RobotOutlined /> AI 智能助手 ({currentUser?.username})</span>
                            <Button size="small" icon={<ClearOutlined />} onClick={clearChat} danger>清空我的对话</Button>
                        </div>
                    }
                    className="chat-card"
                    style={{ 
                        flex: 1, 
                        display: 'flex', 
                        flexDirection: 'column',
                        borderRadius: '16px',
                        overflow: 'hidden',
                        boxShadow: '0 8px 32px rgba(0,0,0,0.1)',
                        border: 'none',
                        background: 'rgba(255, 255, 255, 0.8)',
                        backdropFilter: 'blur(8px)'
                    }}
                    bodyStyle={{ 
                        flex: 1, 
                        overflow: 'hidden', 
                        display: 'flex', 
                        flexDirection: 'column',
                        padding: '0'
                    }}
                >
                    <div 
                        ref={scrollRef}
                        style={{ 
                            flex: 1, 
                            overflowY: 'auto', 
                            padding: '20px',
                            background: '#f8fafc'
                        }}
                    >
                        <List
                            itemLayout="horizontal"
                            dataSource={messages}
                            renderItem={(item) => (
                                <List.Item style={{ 
                                    border: 'none', 
                                    padding: '12px 0',
                                    justifyContent: item.role === 'user' ? 'flex-end' : 'flex-start'
                                }}>
                                    <div style={{ 
                                        display: 'flex', 
                                        flexDirection: item.role === 'user' ? 'row-reverse' : 'row',
                                        maxWidth: '80%',
                                        gap: '12px',
                                        alignItems: 'flex-start'
                                    }}>
                                        <Avatar 
                                            icon={item.role === 'user' ? <UserOutlined /> : <RobotOutlined />} 
                                            style={{ 
                                                backgroundColor: item.role === 'user' ? '#1890ff' : '#52c41a',
                                                flexShrink: 0
                                            }} 
                                        />
                                        <div style={{
                                            padding: '12px 16px',
                                            borderRadius: '12px',
                                            backgroundColor: item.role === 'user' ? '#1890ff' : '#fff',
                                            color: item.role === 'user' ? '#fff' : '#333',
                                            boxShadow: '0 2px 8px rgba(0,0,0,0.05)',
                                            position: 'relative',
                                            whiteSpace: 'pre-wrap',
                                            borderTopLeftRadius: item.role === 'ai' ? '2px' : '12px',
                                            borderTopRightRadius: item.role === 'user' ? '2px' : '12px'
                                        }}>
                                            <Text style={{ color: 'inherit' }}>{item.content}</Text>
                                        </div>
                                    </div>
                                </List.Item>
                            )}
                        />
                        {loading && (
                            <div style={{ margin: '12px 0', display: 'flex', gap: '12px' }}>
                                <Avatar icon={<RobotOutlined />} style={{ backgroundColor: '#52c41a' }} />
                                <div style={{ 
                                    padding: '12px 16px', 
                                    borderRadius: '12px', 
                                    backgroundColor: '#fff',
                                    borderTopLeftRadius: '2px'
                                }}>
                                    <Spin size="small" tip="AI 正在思考..." />
                                </div>
                            </div>
                        )}
                    </div>

                    <div style={{ padding: '20px', background: '#fff', borderTop: '1px solid #f0f0f0' }}>
                        <div style={{ display: 'flex', gap: '12px' }}>
                            <Input.TextArea
                                value={inputValue}
                                onChange={(e) => setInputValue(e.target.value)}
                                onPressEnter={(e) => {
                                    if (!e.shiftKey) {
                                        e.preventDefault();
                                        handleSend();
                                    }
                                }}
                                placeholder="输入您的问题..."
                                autoSize={{ minRows: 1, maxRows: 4 }}
                                style={{ borderRadius: '8px' }}
                            />
                            <Button 
                                type="primary" 
                                icon={<SendOutlined />} 
                                onClick={handleSend}
                                loading={loading}
                                style={{ height: 'auto', borderRadius: '8px' }}
                            >
                                发送
                            </Button>
                        </div>
                    </div>
                </Card>
            </Content>
            
            <style jsx>{`
                .chat-layout {
                    animation: fadeIn 0.5s ease-out;
                }
                @keyframes fadeIn {
                    from { opacity: 0; transform: translateY(10px); }
                    to { opacity: 1; transform: translateY(0); }
                }
                ::-webkit-scrollbar {
                    width: 6px;
                }
                ::-webkit-scrollbar-thumb {
                    background: #e2e8f0;
                    border-radius: 3px;
                }
                ::-webkit-scrollbar-thumb:hover {
                    background: #cbd5e1;
                }
            `}</style>
        </Layout>
    );
};

export default AiChat;
