import React, { useState, useEffect } from 'react';
import { Form, Input, Button, Select, Table, Space, Tag, Typography } from 'antd';
import { SearchOutlined, ReloadOutlined } from '@ant-design/icons';
import { getLotteryAdminRecords } from '../api/lottery';
import dayjs from 'dayjs';

const { Option } = Select;
const { Text } = Typography;

const LotteryRecordsAdmin = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [records, setRecords] = useState([]);
  const [pagination, setPagination] = useState({ current: 1, pageSize: 10, total: 0 });

  const fetchRecords = async (params = {}) => {
    try {
      setLoading(true);
      const res = await getLotteryAdminRecords({
        page: params.current || pagination.current,
        pageSize: params.pageSize || pagination.pageSize,
        ...form.getFieldsValue()
      });
      if (res && res.records) {
        setRecords(res.records);
        setPagination({
          ...pagination,
          current: params.current || pagination.current,
          total: res.total,
        });
      }
    } catch (error) {
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRecords();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleTableChange = (newPagination) => {
    fetchRecords(newPagination);
  };

  const handleSearch = () => {
    fetchRecords({ current: 1 });
  };

  const handleReset = () => {
    form.resetFields();
    fetchRecords({ current: 1 });
  };

  const columns = [
    {
      title: '流水ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
      render: (text) => text || '未知用户',
    },
    {
      title: '所属抽奖活动',
      dataIndex: 'activity_title',
      key: 'activity_title',
      render: (text, record) => (
        <Space>
          {text}
          {record.activity_status === 1 ? (
            <Tag color="success">启用中</Tag>
          ) : (
            <Tag color="default">已禁用</Tag>
          )}
        </Space>
      ),
    },
    {
      title: '获得奖品',
      dataIndex: 'prize_name',
      key: 'prize_name',
      render: (text, record) => (
        <Text strong style={{ color: record.is_hit ? '#ff4d4f' : '#8c8c8c' }}>
          {text}
        </Text>
      ),
    },
    {
      title: '抽奖类型',
      key: 'is_hit',
      width: 100,
      render: (_, record) => (
        record.is_hit ? <Tag color="error">中大奖</Tag> : <Tag color="default">兜底/谢谢惠顾</Tag>
      ),
    },
    {
      title: '抽奖时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (text) => dayjs(text).format('YYYY-MM-DD HH:mm:ss'),
    },
  ];



  return (
    <div>
      <Form form={form} layout="inline" style={{ marginBottom: 20 }}>
        <Form.Item name="title" label="活动名称">
          <Input placeholder="输入主题关键字检索" allowClear />
        </Form.Item>
        <Form.Item name="status" label="活动状态">
          <Select placeholder="不限" style={{ width: 120 }} allowClear>
            <Option value="1">启用中</Option>
            <Option value="0">已禁用</Option>
          </Select>
        </Form.Item>
        <Form.Item>
          <Space>
            <Button type="primary" icon={<SearchOutlined />} onClick={handleSearch}>
              检索流水
            </Button>
            <Button icon={<ReloadOutlined />} onClick={handleReset}>
              重置
            </Button>
          </Space>
        </Form.Item>
      </Form>

      <Table
        dataSource={records}
        columns={columns}
        rowKey="id"
        loading={loading}
        pagination={{
            ...pagination,
            showSizeChanger: true,
            showTotal: (total) => `共 ${total} 条记录`
        }}
        onChange={handleTableChange}
      />
    </div>
  );
};

export default LotteryRecordsAdmin;
