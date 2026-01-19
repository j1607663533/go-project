import { Row, Col, Card, Statistic, Typography, Space, Progress, Table, Tag } from 'antd';
import {
  ArrowUpOutlined,
  ArrowDownOutlined,
  UserOutlined,
  ShoppingCartOutlined,
  DollarOutlined,
  EyeOutlined,
} from '@ant-design/icons';

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
      color: '#1890ff',
    },
    {
      title: '总订单数',
      value: 8846,
      prefix: <ShoppingCartOutlined />,
      suffix: '单',
      trend: 'up',
      trendValue: 8.2,
      color: '#52c41a',
    },
    {
      title: '总销售额',
      value: 298456,
      prefix: <DollarOutlined />,
      suffix: '元',
      trend: 'down',
      trendValue: 3.1,
      color: '#faad14',
    },
    {
      title: '页面浏览量',
      value: 156789,
      prefix: <EyeOutlined />,
      suffix: '次',
      trend: 'up',
      trendValue: 15.8,
      color: '#722ed1',
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
      render: (amount) => `¥${amount.toLocaleString()}`,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status) => {
        const statusMap = {
          completed: { color: 'success', text: '已完成' },
          pending: { color: 'warning', text: '待处理' },
          processing: { color: 'processing', text: '处理中' },
        };
        return <Tag color={statusMap[status].color}>{statusMap[status].text}</Tag>;
      },
    },
    {
      title: '日期',
      dataIndex: 'date',
      key: 'date',
    },
  ];

  return (
    <div>
      {/* 欢迎标题 */}
      <div style={{ marginBottom: 24 }}>
        <Title level={2}>欢迎回来！👋</Title>
        <Text type="secondary">这是您的数据概览</Text>
      </div>

      {/* 统计卡片 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        {statisticsData.map((stat, index) => (
          <Col xs={24} sm={12} lg={6} key={index}>
            <Card
              bordered={false}
              style={{
                background: `linear-gradient(135deg, ${stat.color}15 0%, ${stat.color}05 100%)`,
                borderLeft: `4px solid ${stat.color}`,
              }}
            >
              <Space direction="vertical" size="small" style={{ width: '100%' }}>
                <Text type="secondary" style={{ fontSize: 14 }}>
                  {stat.title}
                </Text>
                <Statistic
                  value={stat.value}
                  prefix={stat.prefix}
                  suffix={stat.suffix}
                  valueStyle={{ 
                    fontSize: 24, 
                    fontWeight: 'bold',
                    color: stat.color,
                  }}
                />
                <div style={{ display: 'flex', alignItems: 'center', gap: 4 }}>
                  {stat.trend === 'up' ? (
                    <ArrowUpOutlined style={{ color: '#52c41a', fontSize: 12 }} />
                  ) : (
                    <ArrowDownOutlined style={{ color: '#ff4d4f', fontSize: 12 }} />
                  )}
                  <Text
                    style={{
                      fontSize: 12,
                      color: stat.trend === 'up' ? '#52c41a' : '#ff4d4f',
                    }}
                  >
                    {stat.trendValue}%
                  </Text>
                  <Text type="secondary" style={{ fontSize: 12 }}>
                    vs 上月
                  </Text>
                </div>
              </Space>
            </Card>
          </Col>
        ))}
      </Row>

      {/* 数据图表和最近订单 */}
      <Row gutter={[16, 16]}>
        {/* 销售趋势 */}
        <Col xs={24} lg={12}>
          <Card
            title="销售趋势"
            bordered={false}
            extra={<a href="#">查看详情</a>}
          >
            <Space direction="vertical" size="large" style={{ width: '100%' }}>
              <div>
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 8 }}>
                  <Text>本月目标完成度</Text>
                  <Text strong>75%</Text>
                </div>
                <Progress percent={75} strokeColor="#1890ff" />
              </div>
              <div>
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 8 }}>
                  <Text>客户满意度</Text>
                  <Text strong>92%</Text>
                </div>
                <Progress percent={92} strokeColor="#52c41a" />
              </div>
              <div>
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 8 }}>
                  <Text>订单完成率</Text>
                  <Text strong>88%</Text>
                </div>
                <Progress percent={88} strokeColor="#722ed1" />
              </div>
            </Space>
          </Card>
        </Col>

        {/* 快速操作 */}
        <Col xs={24} lg={12}>
          <Card title="快速操作" bordered={false}>
            <Row gutter={[16, 16]}>
              <Col span={12}>
                <Card
                  hoverable
                  style={{
                    textAlign: 'center',
                    background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                    color: '#fff',
                  }}
                >
                  <ShoppingCartOutlined style={{ fontSize: 32, marginBottom: 8 }} />
                  <div>新建订单</div>
                </Card>
              </Col>
              <Col span={12}>
                <Card
                  hoverable
                  style={{
                    textAlign: 'center',
                    background: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
                    color: '#fff',
                  }}
                >
                  <UserOutlined style={{ fontSize: 32, marginBottom: 8 }} />
                  <div>添加用户</div>
                </Card>
              </Col>
              <Col span={12}>
                <Card
                  hoverable
                  style={{
                    textAlign: 'center',
                    background: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
                    color: '#fff',
                  }}
                >
                  <DollarOutlined style={{ fontSize: 32, marginBottom: 8 }} />
                  <div>财务报表</div>
                </Card>
              </Col>
              <Col span={12}>
                <Card
                  hoverable
                  style={{
                    textAlign: 'center',
                    background: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
                    color: '#fff',
                  }}
                >
                  <EyeOutlined style={{ fontSize: 32, marginBottom: 8 }} />
                  <div>数据分析</div>
                </Card>
              </Col>
            </Row>
          </Card>
        </Col>
      </Row>

      {/* 最近订单 */}
      <Card
        title="最近订单"
        bordered={false}
        style={{ marginTop: 16 }}
        extra={<a href="/orders">查看全部</a>}
      >
        <Table
          columns={orderColumns}
          dataSource={recentOrders}
          pagination={false}
        />
      </Card>
    </div>
  );
};

export default AntdHome;
