import React, { useState, useEffect } from 'react';
import { Card, Form, Input, InputNumber, Button, Switch, DatePicker, message, Space, Divider, Select, Table, Modal, Popconfirm, Tabs } from 'antd';
import { MinusCircleOutlined, PlusOutlined, EditOutlined, DeleteOutlined, GiftOutlined, UnorderedListOutlined } from '@ant-design/icons';
import { getLotteryActivities, saveLotteryActivity, toggleLotteryActivityStatus, deleteLotteryActivity } from '../api/lottery';
import LotteryRecordsAdmin from '../components/LotteryRecordsAdmin';
import dayjs from 'dayjs';

const { Option } = Select;

const LotteryAdmin = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [activities, setActivities] = useState([]);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [editingId, setEditingId] = useState(null);

  useEffect(() => {
    fetchActivities();
  }, []);

  const fetchActivities = async () => {
    try {
      setLoading(true);
      const res = await getLotteryActivities();
      if (res) {
        setActivities(res);
      }
    } catch (error) {
      message.error("获取活动列表失败");
    } finally {
      setLoading(false);
    }
  };

  const handleToggleStatus = async (id, currentStatus) => {
    try {
      setLoading(true);
      const newStatus = currentStatus === 1 ? 0 : 1;
      await toggleLotteryActivityStatus(id, newStatus);
      message.success("状态更新成功");
      fetchActivities(); // 刷新列表
    } catch (e) {
      message.error("状态更新失败");
      setLoading(false);
    }
  };

  const handleDelete = async (id) => {
    try {
      setLoading(true);
      await deleteLotteryActivity(id);
      message.success("删除成功");
      fetchActivities();
    } catch (e) {
      message.error("删除失败");
      setLoading(false);
    }
  };

  const showEditModal = (record) => {
    setEditingId(record ? record.activity.id : null);
    if (record) {
      const { activity, prizes } = record;
      form.setFieldsValue({
        title: activity.title,
        timeRange: [dayjs(activity.start_time), dayjs(activity.end_time)],
        daily_limit: activity.daily_limit,
        status: activity.status === 1,
        prizes: prizes.map(p => ({
          name: p.name,
          total_stock: p.total_stock,
          weight: p.weight,
          type: p.type,
          sort: p.sort
        }))
      });
    } else {
      // 清空为无默认值状态
      form.setFieldsValue({
        title: "",
        timeRange: null,
        daily_limit: 1,
        status: false,
        prizes: []
      });
    }
    setIsModalVisible(true);
  };

  const handleModalCancel = () => {
    setIsModalVisible(false);
    form.resetFields();
  };

  const onFinish = async (values) => {
    try {
      setLoading(true);
      const payload = {
        id: editingId || 0,
        title: values.title,
        start_time: values.timeRange[0].toISOString(),
        end_time: values.timeRange[1].toISOString(),
        daily_limit: values.daily_limit,
        status: values.status ? 1 : 0,
        prizes: values.prizes
      };

      await saveLotteryActivity(payload);
      message.success("保存成功");
      setIsModalVisible(false);
      form.resetFields();
      fetchActivities();
    } catch (e) {
      if (e.message) {
        message.error("保存失败：" + e.message);
      } else {
        message.error("保存失败");
      }
      setLoading(false);
    }
  };

  const columns = [
    {
      title: '活动名称',
      dataIndex: ['activity', 'title'],
      key: 'title',
    },
    {
      title: '活动时间',
      key: 'time',
      render: (_, record) => {
        const start = dayjs(record.activity.start_time).format('YYYY-MM-DD HH:mm');
        const end = dayjs(record.activity.end_time).format('YYYY-MM-DD HH:mm');
        return `${start} 至 ${end}`;
      }
    },
    {
      title: '每日限次',
      dataIndex: ['activity', 'daily_limit'],
      key: 'daily_limit',
      width: 100,
    },
    {
      title: '状态',
      key: 'status',
      width: 120,
      render: (_, record) => {
        const isActive = record.activity.status === 1;
        return (
          <Switch 
            checkedChildren="启用中" 
            unCheckedChildren="已禁用" 
            checked={isActive} 
            onChange={() => handleToggleStatus(record.activity.id, record.activity.status)} 
          />
        );
      }
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      render: (_, record) => (
        <Space size="middle">
          <Button type="link" icon={<EditOutlined />} onClick={() => showEditModal(record)}>编辑</Button>
          <Popconfirm title="确定要删除这个抽奖活动吗?" onConfirm={() => handleDelete(record.activity.id)}>
            <Button type="link" danger icon={<DeleteOutlined />}>删除</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div style={{ padding: 24, maxWidth: 1200, margin: '0 auto' }}>
      <Card title="抽奖系统管理">
        <Tabs defaultActiveKey="1">
          <Tabs.TabPane tab={<span><GiftOutlined />活动配置</span>} key="1">
            <div style={{ marginBottom: 16, textAlign: 'right' }}>
              <Button type="primary" icon={<PlusOutlined />} onClick={() => showEditModal(null)}>
                新建抽奖活动
              </Button>
            </div>
            <Table
                dataSource={activities}
                columns={columns}
                rowKey={(record) => record.activity.id}
                loading={loading}
                pagination={false}
            />
          </Tabs.TabPane>
          <Tabs.TabPane tab={<span><UnorderedListOutlined />全站抽奖统计</span>} key="2">
             <LotteryRecordsAdmin />
          </Tabs.TabPane>
        </Tabs>
      </Card>

        <Modal
            title={editingId ? "编辑抽奖活动" : "新建抽奖活动"}
            open={isModalVisible}
            onCancel={handleModalCancel}
            onOk={() => form.submit()}
            confirmLoading={loading}
            width={800}
            destroyOnClose
        >
            <Form form={form} layout="vertical" onFinish={onFinish}>
                <Space size="large" style={{ display: 'flex', width: '100%' }}>
                    <Form.Item name="title" label="活动名称" rules={[{ required: true }]}>
                        <Input placeholder="输入主题，例如: 周年庆抽奖" />
                    </Form.Item>
                    
                    <Form.Item name="status" label="开启状态" valuePropName="checked" help="开启此活动将自动禁用其他活动">
                        <Switch checkedChildren="开启" unCheckedChildren="关闭" />
                    </Form.Item>
                    
                    <Form.Item name="daily_limit" label="每人每日次数">
                        <InputNumber min={1} max={100} />
                    </Form.Item>
                </Space>
                
                <Form.Item name="timeRange" label="活动开放时间" rules={[{ required: true }]}>
                    <DatePicker.RangePicker showTime />
                </Form.Item>

                <Divider orientation="left">奖项配置 (至少配置2个以上，推荐偶数个以确保大转盘美观)</Divider>
                <p style={{color: 'gray', fontSize: 13}}>库存填 -1 表示无限库存。转盘中奖几率 = 该奖品权重 / 剩余所有实物奖品权重总和</p>
                
                <Form.List name="prizes">
                {(fields, { add, remove }) => (
                    <>
                    {fields.map(({ key, name, ...restField }) => (
                        <Space key={key} style={{ display: 'flex', marginBottom: 8 }} align="baseline">
                        <Form.Item
                            {...restField}
                            name={[name, 'name']}
                            rules={[{ required: true, message: '请输入奖项名' }]}
                            label="奖品名称"
                        >
                            <Input placeholder="例如: 谢谢惠顾" />
                        </Form.Item>

                        <Form.Item
                            {...restField}
                            name={[name, 'type']}
                            label="类型"
                        >
                            <Select style={{ width: 120 }}>
                            <Option value={1}>实物/虚拟奖</Option>
                            <Option value={2}>积分</Option>
                            <Option value={3}>谢谢(兜底)</Option>
                            </Select>
                        </Form.Item>

                        <Form.Item
                            {...restField}
                            name={[name, 'total_stock']}
                            label="总库存"
                        >
                            <InputNumber placeholder="发完即止" />
                        </Form.Item>

                        <Form.Item
                            {...restField}
                            name={[name, 'weight']}
                            label="权重"
                        >
                            <InputNumber min={0} />
                        </Form.Item>

                        <Form.Item
                            {...restField}
                            name={[name, 'sort']}
                            label="排序"
                        >
                            <InputNumber min={1} />
                        </Form.Item>

                        <MinusCircleOutlined onClick={() => remove(name)} style={{ color: 'red', cursor: 'pointer' }} />
                        </Space>
                    ))}
                    <Form.Item>
                        <Button type="dashed" onClick={() => add()} block icon={<PlusOutlined />}>
                            添加转盘奖项格子
                        </Button>
                    </Form.Item>
                    </>
                )}
                </Form.List>
            </Form>
        </Modal>
    </div>
  );
};

export default LotteryAdmin;
