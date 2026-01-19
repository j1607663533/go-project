import { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Container,
  Paper,
  TextField,
  Button,
  Typography,
  InputAdornment,
  IconButton,
  Alert,
  CircularProgress,
  Avatar,
} from '@mui/material';
import {
  Visibility,
  VisibilityOff,
  LockOutlined,
  Refresh,
} from '@mui/icons-material';
import { getCaptcha, login, isAuthenticated } from '../api/auth';

const Login = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    username: '',
    password: '',
    captcha: '',
  });
  const [captchaData, setCaptchaData] = useState(null);
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const captchaLoaded = useRef(false); // 防止重复加载

  // 检查是否已登录并获取验证码
  useEffect(() => {
    if (isAuthenticated()) {
      navigate('/');
      return;
    }

    // 只加载一次验证码
    if (!captchaLoaded.current) {
      captchaLoaded.current = true;
      loadCaptcha();
    }
  }, [navigate]);

  const loadCaptcha = async () => {
    try {
      const data = await getCaptcha();
      setCaptchaData(data);
    } catch (err) {
      setError('获取验证码失败：' + err.message);
    }
  };

  const refreshCaptcha = async () => {
    try {
      if (captchaData) {
        const data = await getCaptcha(captchaData.captcha_id);
        setCaptchaData(data);
        setFormData({ ...formData, captcha: '' });
      }
    } catch (err) {
      setError('刷新验证码失败：' + err.message);
    }
  };

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
    setError('');
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await login(
        formData.username,
        formData.password,
        captchaData.captcha_id,
        formData.captcha
      );
      
      // 登录成功，跳转到首页
      navigate('/');
    } catch (err) {
      setError(err.message);
      // 刷新验证码
      refreshCaptcha();
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box
      sx={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
      }}
    >
      <Container maxWidth="sm">
        <Paper
          elevation={24}
          sx={{
            p: 4,
            borderRadius: 3,
            backdropFilter: 'blur(10px)',
            backgroundColor: 'rgba(255, 255, 255, 0.95)',
          }}
        >
          {/* Logo */}
          <Box sx={{ display: 'flex', justifyContent: 'center', mb: 2 }}>
            <Avatar
              sx={{
                width: 64,
                height: 64,
                bgcolor: 'primary.main',
              }}
            >
              <LockOutlined sx={{ fontSize: 32 }} />
            </Avatar>
          </Box>

          {/* 标题 */}
          <Typography
            variant="h4"
            align="center"
            gutterBottom
            sx={{
              fontWeight: 700,
              background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
              WebkitBackgroundClip: 'text',
              WebkitTextFillColor: 'transparent',
            }}
          >
            欢迎登录
          </Typography>

          <Typography
            variant="body2"
            align="center"
            color="text.secondary"
            sx={{ mb: 4 }}
          >
            请输入您的账号和密码
          </Typography>

          {/* 错误提示 */}
          {error && (
            <Alert severity="error" sx={{ mb: 3 }}>
              {error}
            </Alert>
          )}

          {/* 登录表单 */}
          <form onSubmit={handleSubmit}>
            <TextField
              fullWidth
              label="用户名"
              name="username"
              value={formData.username}
              onChange={handleChange}
              margin="normal"
              required
              autoFocus
              disabled={loading}
              sx={{ mb: 2 }}
            />

            <TextField
              fullWidth
              label="密码"
              name="password"
              type={showPassword ? 'text' : 'password'}
              value={formData.password}
              onChange={handleChange}
              margin="normal"
              required
              disabled={loading}
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton
                      onClick={() => setShowPassword(!showPassword)}
                      edge="end"
                    >
                      {showPassword ? <VisibilityOff /> : <Visibility />}
                    </IconButton>
                  </InputAdornment>
                ),
              }}
              sx={{ mb: 2 }}
            />

            {/* 验证码 */}
            <Box sx={{ display: 'flex', gap: 2, mb: 3 }}>
              <TextField
                label="验证码"
                name="captcha"
                value={formData.captcha}
                onChange={handleChange}
                required
                disabled={loading}
                inputProps={{ maxLength: 6 }}
                sx={{ flex: 1 }}
              />
              
              <Box
                sx={{
                  position: 'relative',
                  width: 120,
                  height: 56,
                  cursor: 'pointer',
                  borderRadius: 1,
                  overflow: 'hidden',
                  border: '1px solid',
                  borderColor: 'divider',
                  '&:hover': {
                    borderColor: 'primary.main',
                  },
                }}
                onClick={refreshCaptcha}
              >
                {captchaData && (
                  <img
                    src={captchaData.captcha_image}
                    alt="验证码"
                    style={{
                      width: '100%',
                      height: '100%',
                      objectFit: 'fill',
                    }}
                  />
                )}
                {/* <IconButton
                  size="small"
                  sx={{
                    position: 'absolute',
                    top: 2,
                    right: 2,
                    bgcolor: 'rgba(255, 255, 255, 0.8)',
                    '&:hover': {
                      bgcolor: 'rgba(255, 255, 255, 0.9)',
                    },
                  }}
                >
                  <Refresh fontSize="small" />
                </IconButton> */}
              </Box>
            </Box>

            {/* 登录按钮 */}
            <Button
              type="submit"
              fullWidth
              variant="contained"
              size="large"
              disabled={loading}
              sx={{
                py: 1.5,
                fontSize: '1.1rem',
                fontWeight: 600,
                background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                '&:hover': {
                  background: 'linear-gradient(135deg, #5568d3 0%, #6a3f8f 100%)',
                },
              }}
            >
              {loading ? (
                <CircularProgress size={24} color="inherit" />
              ) : (
                '登 录'
              )}
            </Button>
          </form>

          {/* 底部提示 */}
          <Typography
            variant="body2"
            align="center"
            color="text.secondary"
            sx={{ mt: 3 }}
          >
            还没有账号？
            <Button
              size="small"
              sx={{ ml: 1 }}
              onClick={() => navigate('/register')}
            >
              立即注册
            </Button>
          </Typography>
        </Paper>
      </Container>
    </Box>
  );
};

export default Login;
