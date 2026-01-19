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
  Popconfirm,
  Tag,
  Card,
  Tree,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  KeyOutlined,
} from '@ant-design/icons';
import { getAllRoles, createRole, updateRole, deleteRole, assignMenus } from '../api/role';
import { getMenuTree } from '../api/menu';

const { Option } = Select;
const { TextArea } = Input;

const SystemRoles = () => {
  const [loading, setLoading] = useState(false);
  const [roles, setRoles] = useState([]);
  const [menuTree, setMenuTree] = useState([]);
  const [modalVisible, setModalVisible] = useState(false);
  const [menuModalVisible, setMenuModalVisible] = useState(false);
  const [editingRole, setEditingRole] = useState(null);
  const [selectedRole, setSelectedRole] = useState(null);
  const [selectedMenuKeys, setSelectedMenuKeys] = useState([]);
  const [form] = Form.useForm();

  // 加载角色数据
  const loadRoles = async () => {
    setLoading(true);
    try {
      const response = await getAllRoles();
      setRoles(response || []);
    } catch (error) {
      message.error('加载角色失败: ' + (error.message || '未知错误'));
    } finally {
      setLoading(false);
    }
  };

  // 加载菜单树
  const loadMenuTree = async () => {
    try {
      const response = await getMenuTree();
      const treeData = convertToTreeData(response || []);
      setMenuTree(treeData);
    } catch (error) {
      message.error('加载菜单失败: ' + (error.message || '未知错误'));
    }
  };

  // 转换菜单数据为 Tree 组件需要的格式
  const convertToTreeData = (menus) => {
    return menus.map(menu => ({
      title: menu.name,
      key: menu.id,
      children: menu.children && menu.children.length > 0 
        ? convertToTreeData(menu.children) 
        : undefined,
    }));
  };

  useEffect(() => {
    loadRoles();
    loadMenuTree();
  }, []);

  // 打开新建/编辑对话框
  const handleOpenModal = (role = null) => {
    setEditingRole(role);
    if (role) {
      form.setFieldsValue({
        name: role.name,
        code: role.code,
        description: role.description,
      });
    } else {
      form.resetFields();
    }
    setModalVisible(true);
  };

  // 关闭对话框
  const handleCloseModal = () => {
    setModalVisible(false);
    setEditingRole(null);
    form.resetFields();
  };

  // 提交表单
  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setLoading(true);

      if (editingRole) {
        // 更新角色
        await updateRole(editingRole.id, values);
        message.success('更新角色成功');
      } else {
        // 创建角色
        await createRole(values);
        message.success('创建角色成功');
      }

      handleCloseModal();
      loadRoles();
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

  // 删除角色
  const handleDelete = async (id) => {
    setLoading(true);
    try {
      await deleteRole(id);
      message.success('删除角色成功');
      loadRoles();
    } catch (error) {
      message.error('删除失败: ' + (error.message || '未知错误'));
    } finally {
      setLoading(false);
    }
  };

  // 打开分配菜单对话框
  const handleOpenMenuModal = (role) => {
    setSelectedRole(role);
    // 获取角色已有的菜单ID
    const menuIds = role.menus ? role.menus.map(m => m.id) : [];
    setSelectedMenuKeys(menuIds);
    setMenuModalVisible(true);
  };

  // 关闭分配菜单对话框
  const handleCloseMenuModal = () => {
    setMenuModalVisible(false);
    setSelectedRole(null);
    setSelectedMenuKeys([]);
  };

  // 提交菜单分配
  const handleAssignMenus = async () => {
    if (!selectedRole) return;

    setLoading(true);
    try {
      await assignMenus(selectedRole.id, selectedMenuKeys);
      message.success('分配菜单成功');
      handleCloseMenuModal();
      loadRoles();
    } catch (error) {
      message.error('分配菜单失败: ' + (error.message || '未知错误'));
    } finally {
      setLoading(false);
    }
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
      title: '角色名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '角色编码',
      dataIndex: 'code',
      key: 'code',
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: '超级管理员',
      dataIndex: 'is_super',
      key: 'is_super',
      width: 120,
      render: (isSuper) => (
        <Tag color={isSuper ? 'red' : 'default'}>
          {isSuper ? '是' : '否'}
        </Tag>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status) => (
        <Tag color={status === 1 ? 'success' : 'default'}>
          {status === 1 ? '启用' : '禁用'}
        </Tag>
      ),
    },
    {
      title: '菜单数量',
      dataIndex: 'menus',
      key: 'menus',
      width: 100,
      render: (menus) => menus ? menus.length : 0,
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      render: (_, record) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            icon={<KeyOutlined />}
            onClick={() => handleOpenMenuModal(record)}
          >
            分配菜单
          </Button>
          {!record.is_super && (
            <>
              <Button
                type="link"
                size="small"
                icon={<EditOutlined />}
                onClick={() => handleOpenModal(record)}
              >
                编辑
              </Button>
              <Popconfirm
                title="确定要删除这个角色吗？"
                description="删除后将无法恢复"
                onConfirm={() => handleDelete(record.id)}
                okText="确定"
                cancelText="取消"
              >
                <Button
                  type="link"
                  size="small"
                  danger
                  icon={<DeleteOutlined />}
                >
                  删除
                </Button>
              </Popconfirm>
            </>
          )}
          {record.is_super && (
            <Tag color="red">受保护</Tag>
          )}
        </Space>
      ),
    },
  ];

  return (
    <div>
      <Card
        title="角色管理"
        extra={
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => handleOpenModal()}
          >
            新建角色
          </Button>
        }
      >
        <Table
          columns={columns}
          dataSource={roles}
          rowKey="id"
          loading={loading}
          pagination={{
            pageSize: 10,
            showSizeChanger: true,
            showTotal: (total) => `共 ${total} 条`,
          }}
        />
      </Card>

      {/* 新建/编辑对话框 */}
      <Modal
        title={editingRole ? '编辑角色' : '新建角色'}
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
            label="角色名称"
            name="name"
            rules={[
              { required: true, message: '请输入角色名称' },
              { min: 2, max: 50, message: '角色名称长度为2-50个字符' },
            ]}
          >
            <Input placeholder="请输入角色名称" />
          </Form.Item>

          <Form.Item
            label="角色编码"
            name="code"
            rules={[
              { required: true, message: '请输入角色编码' },
              { min: 2, max: 50, message: '角色编码长度为2-50个字符' },
              { pattern: /^[a-z_]+$/, message: '只能包含小写字母和下划线' },
            ]}
          >
            <Input placeholder="例如: editor" disabled={!!editingRole} />
          </Form.Item>

          <Form.Item
            label="描述"
            name="description"
            rules={[
              { max: 200, message: '描述长度不能超过200个字符' },
            ]}
          >
            <TextArea rows={4} placeholder="请输入角色描述" />
          </Form.Item>
        </Form>
      </Modal>

      {/* 分配菜单对话框 */}
      <Modal
        title={`为 "${selectedRole?.name}" 分配菜单`}
        open={menuModalVisible}
        onOk={handleAssignMenus}
        onCancel={handleCloseMenuModal}
        confirmLoading={loading}
        width={600}
      >
        {selectedRole?.is_super ? (
          <div style={{ padding: '20px', textAlign: 'center' }}>
            <Tag color="red" style={{ fontSize: 16, padding: '8px 16px' }}>
              超级管理员角色的菜单不能修改
            </Tag>
          </div>
        ) : (
          <Tree
            checkable
            checkedKeys={selectedMenuKeys}
            onCheck={(checkedKeys) => setSelectedMenuKeys(checkedKeys)}
            treeData={menuTree}
            defaultExpandAll
          />
        )}
      </Modal>
    </div>
  );
};

export default SystemRoles;
