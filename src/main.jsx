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
        path: "debug",
        element: <MenuDebug />,
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
          colorPrimary: '#667eea',
          borderRadius: 8,
          fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif',
        },
        algorithm: theme.defaultAlgorithm,
      }}
    >
      <RouterProvider router={router} />
    </ConfigProvider>
  </React.StrictMode>,
)
