import { useState, useEffect } from 'react';
import {
  Table,
  Button,
  Space,
  Modal,
  Form,
  Input,
  InputNumber,
  Select,
  Switch,
  message,
  Popconfirm,
  Tag,
  Card,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  HomeOutlined,
  ShoppingCartOutlined,
  UserOutlined,
  SettingOutlined,
  TeamOutlined,
  MenuOutlined,
  FileAddOutlined, 
} from '@ant-design/icons';
import { getMenuTree, createMenu, updateMenu, deleteMenu } from '../api/menu';

const { Option } = Select;

// 图标选项
const iconOptions = [
  { label: '首页', value: 'HomeOutlined', icon: <HomeOutlined /> },
  { label: '购物车', value: 'ShoppingCartOutlined', icon: <ShoppingCartOutlined /> },
  { label: '用户', value: 'UserOutlined', icon: <UserOutlined /> },
  { label: '设置', value: 'SettingOutlined', icon: <SettingOutlined /> },
  { label: '团队', value: 'TeamOutlined', icon: <TeamOutlined /> },
  { label: '菜单', value: 'MenuOutlined', icon: <MenuOutlined /> },
  { label: '文件', value: 'FileAddOutlined', icon: <FileAddOutlined /> },
];

const SystemMenus = () => {
  const [loading, setLoading] = useState(false);
  const [menuData, setMenuData] = useState([]);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingMenu, setEditingMenu] = useState(null);
  const [form] = Form.useForm();

  // 加载菜单数据
  const loadMenus = async () => {
    setLoading(true);
    try {
      const response = await getMenuTree();
      // 将树形数据转换为平铺数据以便在表格中显示
      const flatData = flattenMenuTree(response || []);
      setMenuData(flatData);
    } catch (error) {
      message.error('加载菜单失败: ' + (error.message || '未知错误'));
    } finally {
      setLoading(false);
    }
  };

  // 将树形菜单数据转换为平铺数据
  const flattenMenuTree = (menus, level = 0, parentName = '') => {
    let result = [];
    menus.forEach(menu => {
      result.push({
        ...menu,
        level,
        parentName,
      });
      // if (menu.children && menu.children.length > 0) {
      //   result = result.concat(flattenMenuTree(menu.children, level + 1, menu.name));
      // }
    });
    return result;
  };

  useEffect(() => {
    loadMenus();
  }, []);

  // 打开新建/编辑对话框
  const handleOpenModal = (menu = null) => {
    setEditingMenu(menu);
    if (menu) {
      form.setFieldsValue({
        parent_id: menu.parent_id || 0,
        name: menu.name,
        path: menu.path,
        component: menu.component,
        icon: menu.icon,
        sort: menu.sort,
        type: menu.type,
        hidden: menu.hidden,
      });
    } else {
      form.resetFields();
      form.setFieldsValue({
        parent_id: 0,
        type: 1,
        sort: 0,
        hidden: false,
      });
    }
    setModalVisible(true);
  };

  // 关闭对话框
  const handleCloseModal = () => {
    setModalVisible(false);
    setEditingMenu(null);
    form.resetFields();
  };

  // 提交表单
  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setLoading(true);

      if (editingMenu) {
        // 更新菜单
        await updateMenu(editingMenu.id, values);
        message.success('更新菜单成功');
      } else {
        // 创建菜单
        await createMenu(values);
        message.success('创建菜单成功');
      }

      handleCloseModal();
      loadMenus();
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

  // 删除菜单
  const handleDelete = async (id) => {
    setLoading(true);
    try {
      await deleteMenu(id);
      message.success('删除菜单成功');
      loadMenus();
    } catch (error) {
      message.error('删除失败: ' + (error.message || '未知错误'));
    } finally {
      setLoading(false);
    }
  };

  // 获取父菜单选项
  const getParentMenuOptions = () => {
    const topLevelMenus = menuData.filter(m => m.parent_id === 0);
    return [
      { label: '顶级菜单', value: 0 },
      ...topLevelMenus.map(m => ({ label: m.name, value: m.id })),
    ];
  };

  // 表格列定义
  const columns = [
    {
      title: '菜单名称',
      dataIndex: 'name',
      key: 'name',
      render: (text, record) => (
        <Space>
          {record.level > 0 && <span style={{ marginLeft: record.level * 20 }}>└─</span>}
          {text}
        </Space>
      ),
    },
    {
      title: '图标',
      dataIndex: 'icon',
      key: 'icon',
      width: 80,
      render: (icon) => {
        const iconOption = iconOptions.find(opt => opt.value === icon);
        return iconOption ? iconOption.icon : <MenuOutlined />;
      },
    },
    {
      title: '路径',
      dataIndex: 'path',
      key: 'path',
    },
    {
      title: '组件',
      dataIndex: 'component',
      key: 'component',
    },
    {
      title: '排序',
      dataIndex: 'sort',
      key: 'sort',
      width: 80,
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      width: 100,
      render: (type) => (
        <Tag color={type === 1 ? 'blue' : 'green'}>
          {type === 1 ? '菜单' : '按钮'}
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
      title: '隐藏',
      dataIndex: 'hidden',
      key: 'hidden',
      width: 80,
      render: (hidden) => (
        <Tag color={hidden ? 'warning' : 'default'}>
          {hidden ? '是' : '否'}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
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
          <Popconfirm
            title="确定要删除这个菜单吗？"
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
        </Space>
      ),
    },
  ];

  return (
    <div>
      <Card
        title="菜单管理"
        extra={
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => handleOpenModal()}
          >
            新建菜单
          </Button>
        }
      >
        <Table
          columns={columns}
          dataSource={menuData}
          rowKey="id"
          loading={loading}
          pagination={{
            pageSize: 20,
            showSizeChanger: true,
            showTotal: (total) => `共 ${total} 条`,
          }}
        />
      </Card>

      {/* 新建/编辑对话框 */}
      <Modal
        title={editingMenu ? '编辑菜单' : '新建菜单'}
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
            label="父菜单"
            name="parent_id"
            rules={[{ required: true, message: '请选择父菜单' }]}
          >
            <Select placeholder="请选择父菜单">
              {getParentMenuOptions().map(opt => (
                <Option key={opt.value} value={opt.value}>
                  {opt.label}
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            label="菜单名称"
            name="name"
            rules={[
              { required: true, message: '请输入菜单名称' },
              { min: 2, max: 50, message: '菜单名称长度为2-50个字符' },
            ]}
          >
            <Input placeholder="请输入菜单名称" />
          </Form.Item>

          <Form.Item
            label="路由路径"
            name="path"
            rules={[
              { required: true, message: '请输入路由路径' },
              { max: 200, message: '路径长度不能超过200个字符' },
            ]}
          >
            <Input placeholder="例如: /system/users" />
          </Form.Item>

          <Form.Item
            label="组件路径"
            name="component"
            rules={[
              { max: 200, message: '组件路径长度不能超过200个字符' },
            ]}
          >
            <Input placeholder="例如: SystemUsers" />
          </Form.Item>

          <Form.Item
            label="图标"
            name="icon"
          >
            <Select placeholder="请选择图标">
              {iconOptions.map(opt => (
                <Option key={opt.value} value={opt.value}>
                  <Space>
                    {opt.icon}
                    {opt.label}
                  </Space>
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            label="排序"
            name="sort"
            rules={[{ required: true, message: '请输入排序值' }]}
          >
            <InputNumber min={0} style={{ width: '100%' }} placeholder="数字越小越靠前" />
          </Form.Item>

          <Form.Item
            label="类型"
            name="type"
            rules={[{ required: true, message: '请选择类型' }]}
          >
            <Select placeholder="请选择类型">
              <Option value={1}>菜单</Option>
              <Option value={2}>按钮</Option>
            </Select>
          </Form.Item>

          <Form.Item
            label="是否隐藏"
            name="hidden"
            valuePropName="checked"
          >
            <Switch checkedChildren="是" unCheckedChildren="否" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default SystemMenus;
