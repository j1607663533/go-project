import { useState, useEffect } from 'react';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import {
  Layout,
  Menu,
  Avatar,
  Dropdown,
  Space,
  Typography,
  theme,
  Badge,
} from 'antd';
import {
  HomeOutlined,
  ShoppingCartOutlined,
  UserOutlined,
  LogoutOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  BellOutlined,
  SettingOutlined,
  TeamOutlined,
  MenuOutlined as MenuIconOutlined,
} from '@ant-design/icons';
import { logout, getCurrentUser, getUserMenus } from '../api/auth';

const { Header, Sider, Content } = Layout;
const { Text } = Typography;

// 图标映射
const iconMap = {
  HomeOutlined: <HomeOutlined />,
  ShoppingCartOutlined: <ShoppingCartOutlined />,
  UserOutlined: <UserOutlined />,
  SettingOutlined: <SettingOutlined />,
  TeamOutlined: <TeamOutlined />,
  MenuOutlined: <MenuIconOutlined />,
};

const AntdLayout = () => {
  const [collapsed, setCollapsed] = useState(false);
  const [menuItems, setMenuItems] = useState([]);
  const navigate = useNavigate();
  const location = useLocation();
  const user = getCurrentUser();
  
  const {
    token: { colorBgContainer, borderRadiusLG },
  } = theme.useToken();

  // 加载用户菜单
  useEffect(() => {
    const menus = getUserMenus();
    const items = buildMenuItems(menus);
    setMenuItems(items);
  }, []);

  // 构建菜单项
  const buildMenuItems = (menus) => {
    return menus.map(menu => {
      const item = {
        key: menu.path,
        icon: iconMap[menu.icon] || <MenuIconOutlined />,
        label: menu.name,
      };

      if (menu.children && menu.children.length > 0) {
        item.children = buildMenuItems(menu.children);
      }

      return item;
    });
  };

  // 退出登录
  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  // 用户下拉菜单
  const userMenuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人中心',
      onClick: () => navigate('/'),
    },
    {
      key: 'settings',
      icon: <SettingOutlined />,
      label: '设置',
    },
    {
      type: 'divider',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: handleLogout,
    },
  ];

  const handleMenuClick = ({ key, keyPath }) => {
    // 获取点击的菜单项
    const findMenuItem = (items, targetKey) => {
      for (const item of items) {
        if (item.key === targetKey) {
          return item;
        }
        if (item.children) {
          const found = findMenuItem(item.children, targetKey);
          if (found) return found;
        }
      }
      return null;
    };

    const clickedItem = findMenuItem(menuItems, key);
    
    // 只有没有子菜单的项才导航
    if (clickedItem && (!clickedItem.children || clickedItem.children.length === 0)) {
      navigate(key);
    }
  };

  return (
    <Layout style={{ minHeight: '100vh' }}>
      {/* 侧边栏 */}
      <Sider 
        trigger={null} 
        collapsible 
        collapsed={collapsed}
        style={{
          overflow: 'auto',
          height: '100vh',
          position: 'fixed',
          left: 0,
          top: 0,
          bottom: 0,
        }}
        theme="dark"
      >
        {/* Logo 区域 */}
        <div
          style={{
            height: 64,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            background: 'rgba(255, 255, 255, 0.1)',
            margin: '16px',
            borderRadius: '8px',
            transition: 'all 0.2s',
          }}
        >
          <Text
            style={{
              color: '#fff',
              fontSize: collapsed ? 16 : 20,
              fontWeight: 'bold',
              background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
              WebkitBackgroundClip: 'text',
              WebkitTextFillColor: 'transparent',
              transition: 'all 0.2s',
            }}
          >
            {collapsed ? 'RA' : 'React Admin'}
          </Text>
        </div>

        {/* 菜单 */}
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={handleMenuClick}
          style={{
            borderRight: 0,
          }}
        />
      </Sider>

      {/* 主布局 */}
      <Layout style={{ marginLeft: collapsed ? 80 : 200, transition: 'all 0.2s' }}>
        {/* 顶部导航栏 */}
        <Header
          style={{
            padding: '0 24px',
            background: colorBgContainer,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            boxShadow: '0 1px 4px rgba(0,21,41,.08)',
            position: 'sticky',
            top: 0,
            zIndex: 1,
          }}
        >
          {/* 左侧：折叠按钮和标题 */}
          <Space size="large">
            <div
              onClick={() => setCollapsed(!collapsed)}
              style={{
                fontSize: 18,
                cursor: 'pointer',
                transition: 'color 0.3s',
              }}
              onMouseEnter={(e) => (e.target.style.color = '#1890ff')}
              onMouseLeave={(e) => (e.target.style.color = 'inherit')}
            >
              {collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            </div>
            <Text strong style={{ fontSize: 16 }}>
              {menuItems.find(item => item.key === location.pathname)?.label || '后台管理系统'}
            </Text>
          </Space>

          {/* 右侧：通知和用户信息 */}
          <Space size="large">
            {/* 通知图标 */}
            <Badge count={5} size="small">
              <BellOutlined
                style={{
                  fontSize: 18,
                  cursor: 'pointer',
                  transition: 'color 0.3s',
                }}
                onMouseEnter={(e) => (e.target.style.color = '#1890ff')}
                onMouseLeave={(e) => (e.target.style.color = 'inherit')}
              />
            </Badge>

            {/* 用户信息下拉菜单 */}
            <Dropdown
              menu={{ items: userMenuItems }}
              placement="bottomRight"
              arrow
            >
              <Space style={{ cursor: 'pointer' }}>
                <Avatar
                  style={{
                    background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                  }}
                  size="default"
                >
                  {user?.username?.charAt(0).toUpperCase()}
                </Avatar>
                <Text>{user?.username}</Text>
              </Space>
            </Dropdown>
          </Space>
        </Header>

        {/* 内容区域 */}
        <Content
          style={{
            margin: '24px 16px',
            padding: 24,
            minHeight: 280,
            background: colorBgContainer,
            borderRadius: borderRadiusLG,
          }}
        >
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
};

export default AntdLayout;
