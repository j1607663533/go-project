import { useState, useEffect } from 'react';
import { getUserMenus } from '../api/auth';
import { Card, Typography, Tree, Button, message } from 'antd';
import { ReloadOutlined } from '@ant-design/icons';

const { Title, Paragraph, Text } = Typography;
const { DirectoryTree } = Tree;

const MenuDebug = () => {
  const [menuData, setMenuData] = useState([]);
  const [treeData, setTreeData] = useState([]);

  const loadMenus = () => {
    const menus = getUserMenus();
    setMenuData(menus);
    
    // 转换为树形数据
    const convertToTreeData = (menus) => {
      return menus.map(menu => ({
        title: `${menu.name} (${menu.path})`,
        key: menu.id,
        children: menu.children && menu.children.length > 0 
          ? convertToTreeData(menu.children) 
          : undefined,
      }));
    };
    
    setTreeData(convertToTreeData(menus));
  };

  useEffect(() => {
    loadMenus();
  }, []);

  const handleRefresh = () => {
    loadMenus();
    message.success('菜单数据已刷新');
  };

  return (
    <div style={{ padding: 24 }}>
      <Card
        title={
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Title level={3} style={{ margin: 0 }}>菜单数据调试</Title>
            <Button icon={<ReloadOutlined />} onClick={handleRefresh}>
              刷新
            </Button>
          </div>
        }
      >
        <Paragraph>
          <Text strong>菜单数量：</Text>
          <Text>{menuData.length}</Text>
        </Paragraph>

        <Paragraph>
          <Text strong>原始数据：</Text>
        </Paragraph>
        <pre style={{ 
          background: '#f5f5f5', 
          padding: 16, 
          borderRadius: 4,
          overflow: 'auto',
          maxHeight: 400
        }}>
          {JSON.stringify(menuData, null, 2)}
        </pre>

        <Paragraph style={{ marginTop: 24 }}>
          <Text strong>树形结构：</Text>
        </Paragraph>
        {treeData.length > 0 ? (
          <DirectoryTree
            treeData={treeData}
            defaultExpandAll
          />
        ) : (
          <Text type="secondary">没有菜单数据</Text>
        )}

        <Paragraph style={{ marginTop: 24 }}>
          <Text strong type="warning">提示：</Text>
          <br />
          <Text>如果看不到子菜单，请尝试：</Text>
          <br />
          <Text>1. 退出登录</Text>
          <br />
          <Text>2. 重新登录</Text>
          <br />
          <Text>3. 检查后端是否已重启</Text>
        </Paragraph>
      </Card>
    </div>
  );
};

export default MenuDebug;
