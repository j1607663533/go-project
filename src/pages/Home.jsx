import { useEffect, useState, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Container,
  Paper,
  Typography,
  Button,
  Avatar,
  Card,
  CardContent,
  Grid,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  CircularProgress,
} from '@mui/material';
import {
  Person,
  Email,
  CalendarToday,
  ExitToApp,
  People,
} from '@mui/icons-material';
import { getCurrentUser, logout, isAuthenticated, getAllUsers } from '../api/auth';

const Home = () => {
  const [user, setUser] = useState(null);
  const [users, setUsers] = useState([]);
  const [loadingUsers, setLoadingUsers] = useState(true);
  const usersLoaded = useRef(false);

  useEffect(() => {
    // 获取用户信息
    const currentUser = getCurrentUser();
    setUser(currentUser);

    // 只加载一次用户列表
    if (!usersLoaded.current) {
      usersLoaded.current = true;
      loadUsers();
    }
  }, []);

  const loadUsers = async () => {
    try {
      setLoadingUsers(true);
      const userList = await getAllUsers();
      setUsers(userList);
    } catch (err) {
      console.error('获取用户列表失败：', err);
    } finally {
      setLoadingUsers(false);
    }
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  if (!user) {
    return null;
  }

  return (
    <Container maxWidth="lg">
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" sx={{ fontWeight: 700, mb: 1 }}>
          欢迎回来，{user.nickname || user.username}！
        </Typography>
        <Typography variant="body1" color="text.secondary">
          这是您的个人概况 📊
        </Typography>
      </Box>

      {/* 用户信息卡片 */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} md={3}>
          <Card elevation={2} sx={{ borderRadius: 2 }}>
            <CardContent sx={{ textAlign: 'center' }}>
              <Avatar
                sx={{
                  width: 80,
                  height: 80,
                  bgcolor: 'primary.main',
                  fontSize: '2rem',
                  mx: 'auto',
                  mb: 2,
                }}
              >
                {user.username.charAt(0).toUpperCase()}
              </Avatar>
              <Typography variant="h6">{user.nickname || user.username}</Typography>
              <Typography variant="body2" color="text.secondary">ID: {user.id}</Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} md={9}>
          <Grid container spacing={2}>
            <Grid item xs={12} sm={6}>
              <Card elevation={2} sx={{ borderRadius: 2 }}>
                <CardContent>
                  <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                    <Person sx={{ mr: 1, color: 'primary.main' }} fontSize="small" />
                    <Typography variant="subtitle2" color="text.secondary">用户名</Typography>
                  </Box>
                  <Typography variant="body1">{user.username}</Typography>
                </CardContent>
              </Card>
            </Grid>
            <Grid item xs={12} sm={6}>
              <Card elevation={2} sx={{ borderRadius: 2 }}>
                <CardContent>
                  <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                    <Email sx={{ mr: 1, color: 'primary.main' }} fontSize="small" />
                    <Typography variant="subtitle2" color="text.secondary">电子邮箱</Typography>
                  </Box>
                  <Typography variant="body1">{user.email}</Typography>
                </CardContent>
              </Card>
            </Grid>
            <Grid item xs={12} sm={6}>
              <Card elevation={2} sx={{ borderRadius: 2 }}>
                <CardContent>
                  <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                    <CalendarToday sx={{ mr: 1, color: 'primary.main' }} fontSize="small" />
                    <Typography variant="subtitle2" color="text.secondary">注册时间</Typography>
                  </Box>
                  <Typography variant="body1">{formatDate(user.created_at)}</Typography>
                </CardContent>
              </Card>
            </Grid>
          </Grid>
        </Grid>
      </Grid>

      {/* 用户列表 */}
      <Paper
        elevation={2}
        sx={{
          p: 3,
          borderRadius: 2,
        }}
      >
        <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
          <People sx={{ mr: 1, color: 'primary.main' }} />
          <Typography variant="h6" sx={{ fontWeight: 700 }}>
            系统用户列表
          </Typography>
          <Chip
            label={`共 ${users.length} 人`}
            color="primary"
            size="small"
            variant="outlined"
            sx={{ ml: 2 }}
          />
        </Box>

        {loadingUsers ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
            <CircularProgress size={30} />
          </Box>
        ) : (
          <TableContainer>
            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>ID</TableCell>
                  <TableCell>用户名</TableCell>
                  <TableCell>邮箱</TableCell>
                  <TableCell>注册时间</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {users.map((u) => (
                  <TableRow
                    key={u.id}
                    sx={{
                      backgroundColor: u.id === user.id ? 'rgba(102, 126, 234, 0.05)' : 'inherit',
                    }}
                  >
                    <TableCell>{u.id}</TableCell>
                    <TableCell>{u.username}{u.id === user.id && ' (我)'}</TableCell>
                    <TableCell>{u.email}</TableCell>
                    <TableCell>{formatDate(u.created_at)}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        )}
      </Paper>
    </Container>
  );
};

export default Home;

