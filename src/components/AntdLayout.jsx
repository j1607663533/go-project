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
  RobotOutlined,
} from '@ant-design/icons';
import { logout, getCurrentUser, getUserMenus } from '../api/auth';

const { Header, Sider, Content } = Layout;
const { Text, Title } = Typography;

// 图标映射
const iconMap = {
  HomeOutlined: <HomeOutlined />,
  ShoppingCartOutlined: <ShoppingCartOutlined />,
  UserOutlined: <UserOutlined />,
  SettingOutlined: <SettingOutlined />,
  TeamOutlined: <TeamOutlined />,
  MenuOutlined: <MenuIconOutlined />,
  RobotOutlined: <RobotOutlined />,
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
    
    // 强制加入 AI 对话菜单项
    const hasAiChat = items.some(item => item.key === '/ai-chat');
    if (!hasAiChat) {
      items.push({
        key: '/ai-chat',
        icon: <RobotOutlined />,
        label: 'AI 智能对话',
      });
    }
    
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
            background: 'rgba(255, 255, 255, 0.05)',
            margin: '20px 16px 32px',
            borderRadius: '12px',
            transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
            border: '1px solid rgba(255, 255, 255, 0.1)',
          }}
        >
          <Text
            style={{
              color: '#fff',
              fontSize: collapsed ? 16 : 22,
              fontFamily: 'var(--font-heading)',
              fontWeight: 700,
              background: 'linear-gradient(135deg, #6366f1 0%, #a855f7 100%)',
              WebkitBackgroundClip: 'text',
              WebkitTextFillColor: 'transparent',
              transition: 'all 0.3s',
              letterSpacing: collapsed ? 0 : 1,
            }}
          >
            {collapsed ? 'RA' : 'REACT ADMIN'}
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
            padding: '0 8px',
          }}
        />
      </Sider>

      {/* 主布局 */}
      <Layout style={{ marginLeft: collapsed ? 80 : 200, transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)' }}>
        {/* 顶部导航栏 */}
        <Header
          className="glass-morphism"
          style={{
            padding: '0 24px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            position: 'sticky',
            top: 0,
            zIndex: 10,
            height: 64,
            width: '100%',
          }}
        >
          {/* 左侧：折叠按钮和标题 */}
          <Space size="large">
            <div
              onClick={() => setCollapsed(!collapsed)}
              style={{
                fontSize: 20,
                cursor: 'pointer',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                width: 36,
                height: 36,
                borderRadius: '8px',
                background: 'rgba(99, 102, 241, 0.1)',
                color: '#6366f1',
                transition: 'all 0.2s',
              }}
            >
              {collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            </div>
            <Title level={4} style={{ margin: 0, fontSize: 18, fontFamily: 'var(--font-heading)' }}>
              {menuItems.find(item => item.key === location.pathname)?.label || 'Dashboard Overview'}
            </Title>
          </Space>

          {/* 右侧：通知和用户信息 */}
          <Space size="large">
            {/* 通知图标 */}
            <Badge count={5} size="small" offset={[-2, 6]} color="#f43f5e">
              <div
                style={{
                  width: 36,
                  height: 36,
                  borderRadius: '10px',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  background: '#f8fafc',
                  border: '1px solid #f1f5f9',
                  cursor: 'pointer',
                  transition: 'all 0.2s',
                }}
              >
                <BellOutlined style={{ fontSize: 18, color: '#64748b' }} />
              </div>
            </Badge>

            {/* 用户信息下拉菜单 */}
            <Dropdown
              menu={{ items: userMenuItems }}
              placement="bottomRight"
              arrow={{ pointAtCenter: true }}
            >
              <Space style={{ cursor: 'pointer', padding: '4px 8px', borderRadius: '8px', transition: 'all 0.2s' }} className="hover-bg">
                <Avatar
                  style={{
                    background: 'linear-gradient(135deg, #6366f1 0%, #a855f7 100%)',
                    boxShadow: '0 2px 8px rgba(99, 102, 241, 0.3)',
                  }}
                  size="default"
                >
                  {user?.username?.charAt(0).toUpperCase()}
                </Avatar>
                <Text strong style={{ fontSize: 14 }}>{user?.username}</Text>
              </Space>
            </Dropdown>
          </Space>
        </Header>

        {/* 内容区域 */}
        <Content
          style={{
            margin: '32px 24px',
            padding: 0,
            minHeight: 280,
          }}
        >
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
};

export default AntdLayout;
