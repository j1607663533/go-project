import React from 'react'
import ReactDOM from 'react-dom/client'
import {
  createBrowserRouter,
  RouterProvider,
  Navigate,
} from "react-router-dom";
import { ConfigProvider, theme } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import App from './App.jsx'
import Home from './pages/Home.jsx';
import AntdHome from './pages/AntdHome.jsx';
import Orders from './pages/Orders.jsx';
import Login from './pages/Login.jsx';
import Register from './pages/Register.jsx';
import SystemUsers from './pages/SystemUsers.jsx';
import SystemRoles from './pages/SystemRoles.jsx';
import SystemMenus from './pages/SystemMenus.jsx';
import MenuDebug from './pages/MenuDebug.jsx';
import Layout from './components/Layout.jsx';
import AntdLayout from './components/AntdLayout.jsx';
import AiChat from './pages/AiChat.jsx';
import FileManager from './pages/FileManager.jsx';
import VideoTools from './pages/VideoTools.jsx';
import Lottery from './pages/Lottery.jsx';
import LotteryAdmin from './pages/LotteryAdmin.jsx';
import { isAuthenticated } from './api/auth';
import './index.css'
import 'antd/dist/reset.css';

// 受保护的路由组件
const ProtectedRoute = ({ children }) => {
  return isAuthenticated() ? children : <Navigate to="/login" replace />;
};

// 路由配置
const router = createBrowserRouter([
  {
    path: "/login",
    element: <Login />,
  },
  {
    path: "/register",
    element: <Register />,
  },
  {
    path: "/",
    element: (
      <ProtectedRoute>
        <AntdLayout />
      </ProtectedRoute>
    ),
    children: [
      {
        index: true,
        element: <AntdHome />,
      },
      {
        path: "orders",
        element: <Orders />,
      },
      {
        path: "system/users",
        element: <SystemUsers />,
      },
      {
        path: "system/roles",
        element: <SystemRoles />,
      },
      {
        path: "system/menus",
        element: <SystemMenus />,
      },
      {
        path: "ai-chat",
        element: <AiChat />,
      },
      {
        path: "files",
        element: <ProtectedRoute><FileManager /></ProtectedRoute>,
      },
      {
        path: "video-tools",
        element: <ProtectedRoute><VideoTools /></ProtectedRoute>,
      },
      {
        path: "debug",
        element: <MenuDebug />,
      },
      {
        path: "lottery",
        element: <ProtectedRoute><Lottery /></ProtectedRoute>,
      },
      {
        path: "lottery-admin",
        element: <ProtectedRoute><LotteryAdmin /></ProtectedRoute>,
      },
    ],
  },
  {
    path: "*",
    element: <Navigate to="/" replace />,
  },
]);


ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <ConfigProvider
      locale={zhCN}
      theme={{
        token: {
          colorPrimary: '#6366f1', // Electric Indigo
          colorSuccess: '#10b981', // Vivid Mint
          colorWarning: '#f59e0b',
          colorError: '#f43f5e', // Sunset Rose
          borderRadius: 12,
          fontFamily: "'Plus Jakarta Sans', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif",
          fontSize: 14,
          wireframe: false,
        },
        components: {
          Card: {
            boxShadowTertiary: '0 1px 2px 0 rgb(0 0 0 / 0.05)',
          },
          Button: {
            controlHeight: 38,
            fontWeight: 500,
          },
          Layout: {
            headerBg: 'rgba(255, 255, 255, 0.8)',
            siderBg: '#0f172a',
          },
          Menu: {
            itemBg: 'transparent',
            itemColor: 'rgba(255, 255, 255, 0.65)',
            itemSelectedBg: 'rgba(99, 102, 241, 0.15)',
            itemSelectedColor: '#fff',
          },
        },
        algorithm: theme.defaultAlgorithm,
      }}
    >
      <RouterProvider router={router} />
    </ConfigProvider>
  </React.StrictMode>,
)
