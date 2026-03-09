import { Row, Col, Card, Statistic, Typography, Space, Progress, Table, Tag } from 'antd';
import {
  ArrowUpOutlined,
  ArrowDownOutlined,
  UserOutlined,
  ShoppingCartOutlined,
  DollarOutlined,
  EyeOutlined,
} from '@ant-design/icons';
import { motion } from 'framer-motion';

const { Title, Text } = Typography;

const AntdHome = () => {
  // 统计数据
  const statisticsData = [
    {
      title: '总用户数',
      value: 11893,
      prefix: <UserOutlined />,
      suffix: '人',
      trend: 'up',
      trendValue: 12.5,
      color: '#6366f1',
    },
    {
      title: '总订单数',
      value: 8846,
      prefix: <ShoppingCartOutlined />,
      suffix: '单',
      trend: 'up',
      trendValue: 8.2,
      color: '#10b981',
    },
    {
      title: '总销售额',
      value: 298456,
      prefix: <DollarOutlined />,
      suffix: '元',
      trend: 'down',
      trendValue: 3.1,
      color: '#f59e0b',
    },
    {
      title: '页面浏览量',
      value: 156789,
      prefix: <EyeOutlined />,
      suffix: '次',
      trend: 'up',
      trendValue: 15.8,
      color: '#a855f7',
    },
  ];

  // 最近订单数据
  const recentOrders = [
    {
      key: '1',
      orderId: 'ORD-2024-001',
      customer: '张三',
      product: 'MacBook Pro 16"',
      amount: 18999,
      status: 'completed',
      date: '2024-01-06',
    },
    {
      key: '2',
      orderId: 'ORD-2024-002',
      customer: '李四',
      product: 'iPhone 15 Pro',
      amount: 8999,
      status: 'pending',
      date: '2024-01-06',
    },
    {
      key: '3',
      orderId: 'ORD-2024-003',
      customer: '王五',
      product: 'AirPods Pro',
      amount: 1999,
      status: 'processing',
      date: '2024-01-05',
    },
    {
      key: '4',
      orderId: 'ORD-2024-004',
      customer: '赵六',
      product: 'iPad Air',
      amount: 4799,
      status: 'completed',
      date: '2024-01-05',
    },
  ];

  const orderColumns = [
    {
      title: '订单号',
      dataIndex: 'orderId',
      key: 'orderId',
      render: (text) => <Text strong style={{ color: '#6366f1' }}>{text}</Text>
    },
    {
      title: '客户',
      dataIndex: 'customer',
      key: 'customer',
    },
    {
      title: '产品',
      dataIndex: 'product',
      key: 'product',
    },
    {
      title: '金额',
      dataIndex: 'amount',
      key: 'amount',
      render: (amount) => (
        <Text strong>
          ¥{amount.toLocaleString()}
        </Text>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status) => {
        const statusMap = {
          completed: { color: '#10b981', text: '已完成' },
          pending: { color: '#f59e0b', text: '待处理' },
          processing: { color: '#6366f1', text: '处理中' },
        };
        return (
          <Tag color={statusMap[status].color} style={{ borderRadius: 6, border: 'none' }}>
            {statusMap[status].text}
          </Tag>
        );
      },
    },
    {
      title: '日期',
      dataIndex: 'date',
      key: 'date',
      render: (text) => <Text type="secondary">{text}</Text>
    },
  ];

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1
      }
    }
  };

  const itemVariants = {
    hidden: { y: 20, opacity: 0 },
    visible: {
      y: 0,
      opacity: 1,
      transition: { 
        type: 'spring',
        stiffness: 100
      }
    }
  };

  return (
    <motion.div
      initial="hidden"
      animate="visible"
      variants={containerVariants}
    >
      {/* 欢迎标题 */}
      <motion.div variants={itemVariants} style={{ marginBottom: 32 }}>
        <Title level={1} style={{ margin: 0, fontFamily: 'var(--font-heading)' }}>
          欢迎回来，<span className="gradient-text">管理员</span> 👋
        </Title>
        <Text type="secondary" style={{ fontSize: 16 }}>这是您的实时业务数据概览</Text>
      </motion.div>

      {/* 统计卡片 */}
      <Row gutter={[20, 20]} style={{ marginBottom: 32 }}>
        {statisticsData.map((stat, index) => (
          <Col xs={24} sm={12} lg={6} key={index}>
            <motion.div variants={itemVariants}>
              <Card
                bordered={false}
                style={{
                  background: `linear-gradient(135deg, ${stat.color}08 0%, #ffffff 100%)`,
                  borderTop: `4px solid ${stat.color}`,
                }}
              >
                <Space direction="vertical" size="small" style={{ width: '100%' }}>
                  <Text type="secondary" style={{ fontSize: 13, textTransform: 'uppercase', letterSpacing: 0.5 }}>
                    {stat.title}
                  </Text>
                  <Statistic
                    value={stat.value}
                    prefix={stat.prefix}
                    suffix={<span style={{ fontSize: 14 }}>{stat.suffix}</span>}
                    valueStyle={{ 
                      fontSize: 28, 
                      fontFamily: 'var(--font-heading)',
                      fontWeight: 700,
                      color: '#1e293b',
                    }}
                  />
                  <div style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                    <Tag 
                      color={stat.trend === 'up' ? 'success' : 'error'}
                      icon={stat.trend === 'up' ? <ArrowUpOutlined /> : <ArrowDownOutlined />}
                      style={{ border: 'none', borderRadius: 4, margin: 0 }}
                    >
                      {stat.trendValue}%
                    </Tag>
                    <Text type="secondary" style={{ fontSize: 12 }}>
                      vs 上月
                    </Text>
                  </div>
                </Space>
              </Card>
            </motion.div>
          </Col>
        ))}
      </Row>

      {/* 数据图表和最近订单 */}
      <Row gutter={[20, 20]}>
        {/* 销售趋势 */}
        <Col xs={24} lg={12}>
          <motion.div variants={itemVariants}>
            <Card
              title={<span style={{ fontFamily: 'var(--font-heading)' }}>销售趋势分析</span>}
              bordered={false}
              extra={<a href="#" style={{ fontSize: 13 }}>详细统计 &rarr;</a>}
              style={{ height: '100%' }}
            >
              <Space direction="vertical" size="large" style={{ width: '100%' }}>
                <div>
                  <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 12 }}>
                    <Text strong>本月目标完成度</Text>
                    <Text style={{ color: '#6366f1' }}>75%</Text>
                  </div>
                  <Progress percent={75} strokeColor={{ '0%': '#6366f1', '100%': '#a855f7' }} showInfo={false} strokeWidth={10} />
                </div>
                <div>
                  <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 12 }}>
                    <Text strong>客户满意度</Text>
                    <Text style={{ color: '#10b981' }}>92%</Text>
                  </div>
                  <Progress percent={92} strokeColor="#10b981" showInfo={false} strokeWidth={10} />
                </div>
                <div>
                  <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 12 }}>
                    <Text strong>订单完成率</Text>
                    <Text style={{ color: '#f43f5e' }}>88%</Text>
                  </div>
                  <Progress percent={88} strokeColor="#f43f5e" showInfo={false} strokeWidth={10} />
                </div>
              </Space>
            </Card>
          </motion.div>
        </Col>

        {/* 快速操作 */}
        <Col xs={24} lg={12}>
          <motion.div variants={itemVariants}>
            <Card title={<span style={{ fontFamily: 'var(--font-heading)' }}>快捷功能</span>} bordered={false}>
              <Row gutter={[16, 16]}>
                {[
                  { label: '新建订单', icon: <ShoppingCartOutlined />, gradient: 'linear-gradient(135deg, #6366f1, #818cf8)' },
                  { label: '添加用户', icon: <UserOutlined />, gradient: 'linear-gradient(135deg, #f43f5e, #fb7185)' },
                  { label: '财务报表', icon: <DollarOutlined />, gradient: 'linear-gradient(135deg, #10b981, #34d399)' },
                  { label: '数据分析', icon: <EyeOutlined />, gradient: 'linear-gradient(135deg, #a855f7, #c084fc)' },
                ].map((action, i) => (
                  <Col span={12} key={i}>
                    <Card
                      hoverable
                      style={{
                        textAlign: 'center',
                        background: action.gradient,
                        border: 'none',
                      }}
                      bodyStyle={{ padding: '24px 16px' }}
                    >
                      <div style={{ color: '#fff' }}>
                        <div style={{ fontSize: 24, marginBottom: 8 }}>{action.icon}</div>
                        <div style={{ fontWeight: 600, fontSize: 13, letterSpacing: 0.5 }}>{action.label}</div>
                      </div>
                    </Card>
                  </Col>
                ))}
              </Row>
            </Card>
          </motion.div>
        </Col>
      </Row>

      {/* 最近订单 */}
      <motion.div variants={itemVariants} style={{ marginTop: 24 }}>
        <Card
          title={<span style={{ fontFamily: 'var(--font-heading)' }}>最近交易记录</span>}
          bordered={false}
          extra={<a href="/orders" style={{ fontSize: 13 }}>管理所有订单 &rarr;</a>}
        >
          <Table
            columns={orderColumns}
            dataSource={recentOrders}
            pagination={false}
            size="middle"
          />
        </Card>
      </motion.div>
    </motion.div>
  );
};

export default AntdHome;
