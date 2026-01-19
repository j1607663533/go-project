import { useState, useEffect } from 'react';
import {
  Table,
  Button,
  Space,
  Modal,
  Form,
  Input,
  Select,
  message,
  Tag,
  Card,
  Avatar,
} from 'antd';
import {
  UserOutlined,
  EditOutlined,
} from '@ant-design/icons';
import { getAllUsers, updateUser } from '../api/auth';
import { getAllRoles } from '../api/role';

const { Option } = Select;

const SystemUsers = () => {
  const [loading, setLoading] = useState(false);
  const [users, setUsers] = useState([]);
  const [roles, setRoles] = useState([]);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingUser, setEditingUser] = useState(null);
  const [form] = Form.useForm();

  // 加载用户数据
  const loadUsers = async () => {
    setLoading(true);
    try {
      const response = await getAllUsers();
      setUsers(response || []);
    } catch (error) {
      message.error('加载用户失败: ' + (error.message || '未知错误'));
    } finally {
      setLoading(false);
    }
  };

  // 加载角色数据
  const loadRoles = async () => {
    try {
      const response = await getAllRoles();
      setRoles(response || []);
    } catch (error) {
      message.error('加载角色失败: ' + (error.message || '未知错误'));
    }
  };

  useEffect(() => {
    loadUsers();
    loadRoles();
  }, []);

  // 打开编辑对话框
  const handleOpenModal = (user) => {
    setEditingUser(user);
    form.setFieldsValue({
      role_id: user.role_id,
      nickname: user.nickname,
      email: user.email,
    });
    setModalVisible(true);
  };

  // 关闭对话框
  const handleCloseModal = () => {
    setModalVisible(false);
    setEditingUser(null);
    form.resetFields();
  };

  // 提交表单
  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setLoading(true);

      // 更新用户
      await updateUser(editingUser.id, values);
      message.success('更新用户成功');

      handleCloseModal();
      loadUsers();
    } catch (error) {
      if (error.errorFields) {
        message.error('请检查表单输入');
      } else {
        message.error('操作失败: ' + (error.message || '未知错误'));
      }
    } finally {
      setLoading(false);
    }
  };

  // 获取角色名称
  const getRoleName = (roleId) => {
    const role = roles.find(r => r.id === roleId);
    return role ? role.name : '未知';
  };

  // 表格列定义
  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '头像',
      dataIndex: 'avatar',
      key: 'avatar',
      width: 80,
      render: (avatar, record) => (
        <Avatar
          src={avatar}
          icon={<UserOutlined />}
          style={{
            background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
          }}
        >
          {record.username?.charAt(0).toUpperCase()}
        </Avatar>
      ),
    },
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
    },
    {
      title: '昵称',
      dataIndex: 'nickname',
      key: 'nickname',
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      key: 'email',
    },
    {
      title: '角色',
      dataIndex: 'role_id',
      key: 'role_id',
      render: (roleId) => {
        const role = roles.find(r => r.id === roleId);
        return (
          <Tag color={role?.is_super ? 'red' : 'blue'}>
            {getRoleName(roleId)}
          </Tag>
        );
      },
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (time) => new Date(time).toLocaleString('zh-CN'),
    },
    {
      title: '操作',
      key: 'action',
      width: 120,
      render: (_, record) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleOpenModal(record)}
          >
            编辑
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <Card title="用户管理">
        <Table
          columns={columns}
          dataSource={users}
          rowKey="id"
          loading={loading}
          pagination={{
            pageSize: 10,
            showSizeChanger: true,
            showTotal: (total) => `共 ${total} 条`,
          }}
        />
      </Card>

      {/* 编辑对话框 */}
      <Modal
        title="编辑用户"
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={handleCloseModal}
        confirmLoading={loading}
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          autoComplete="off"
        >
          <Form.Item
            label="角色"
            name="role_id"
            rules={[{ required: true, message: '请选择角色' }]}
          >
            <Select placeholder="请选择角色">
              {roles.map(role => (
                <Option key={role.id} value={role.id}>
                  {role.name}
                  {role.is_super && <Tag color="red" style={{ marginLeft: 8 }}>超级管理员</Tag>}
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            label="昵称"
            name="nickname"
            rules={[
              { max: 50, message: '昵称长度不能超过50个字符' },
            ]}
          >
            <Input placeholder="请输入昵称" />
          </Form.Item>

          <Form.Item
            label="邮箱"
            name="email"
            rules={[
              { required: true, message: '请输入邮箱' },
              { type: 'email', message: '请输入有效的邮箱地址' },
              { max: 100, message: '邮箱长度不能超过100个字符' },
            ]}
          >
            <Input placeholder="请输入邮箱" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default SystemUsers;
