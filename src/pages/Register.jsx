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
  Stepper,
  Step,
  StepLabel,
} from '@mui/material';
import {
  Visibility,
  VisibilityOff,
  PersonAdd,
  Refresh,
} from '@mui/icons-material';
import { getCaptcha, register, isAuthenticated } from '../api/auth';

const Register = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
    nickname: '',
    captcha: '',
  });
  const [captchaData, setCaptchaData] = useState(null);
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
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

  const validateForm = () => {
    if (formData.username.length < 3 || formData.username.length > 20) {
      setError('用户名长度必须在3-20个字符之间');
      return false;
    }

    if (!/^[a-zA-Z0-9]+$/.test(formData.username)) {
      setError('用户名只能包含字母和数字');
      return false;
    }

    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      setError('请输入有效的邮箱地址');
      return false;
    }

    if (formData.password.length < 6) {
      setError('密码长度至少为6个字符');
      return false;
    }

    if (formData.password !== formData.confirmPassword) {
      setError('两次输入的密码不一致');
      return false;
    }

    if (formData.captcha.length !== 6) {
      setError('请输入6位验证码');
      return false;
    }

    return true;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    if (!validateForm()) {
      return;
    }

    setLoading(true);

    try {
      await register(
        formData.username,
        formData.email,
        formData.password,
        formData.nickname,
        captchaData.captcha_id,
        formData.captcha
      );

      setSuccess(true);
      
      // 3秒后跳转到登录页面
      setTimeout(() => {
        navigate('/login');
      }, 3000);
    } catch (err) {
      setError(err.message);
      // 刷新验证码
      refreshCaptcha();
    } finally {
      setLoading(false);
    }
  };

  if (success) {
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
              p: 6,
              borderRadius: 3,
              textAlign: 'center',
              background: 'rgba(255, 255, 255, 0.95)',
            }}
          >
            <Avatar
              sx={{
                width: 80,
                height: 80,
                bgcolor: 'success.main',
                margin: '0 auto 24px',
              }}
            >
              <PersonAdd sx={{ fontSize: 40 }} />
            </Avatar>
            
            <Typography variant="h4" gutterBottom sx={{ fontWeight: 700, color: 'success.main' }}>
              注册成功！
            </Typography>
            
            <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
              欢迎加入我们！即将跳转到登录页面...
            </Typography>
            
            <CircularProgress />
          </Paper>
        </Container>
      </Box>
    );
  }

  return (
    <Box
      sx={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        py: 4,
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
              <PersonAdd sx={{ fontSize: 32 }} />
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
            创建账号
          </Typography>

          <Typography
            variant="body2"
            align="center"
            color="text.secondary"
            sx={{ mb: 4 }}
          >
            填写以下信息完成注册
          </Typography>

          {/* 错误提示 */}
          {error && (
            <Alert severity="error" sx={{ mb: 3 }}>
              {error}
            </Alert>
          )}

          {/* 注册表单 */}
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
              helperText="3-20个字符，只能包含字母和数字"
              sx={{ mb: 2 }}
            />

            <TextField
              fullWidth
              label="邮箱"
              name="email"
              type="email"
              value={formData.email}
              onChange={handleChange}
              margin="normal"
              required
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
              helperText="至少6个字符"
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

            <TextField
              fullWidth
              label="确认密码"
              name="confirmPassword"
              type={showConfirmPassword ? 'text' : 'password'}
              value={formData.confirmPassword}
              onChange={handleChange}
              margin="normal"
              required
              disabled={loading}
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton
                      onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                      edge="end"
                    >
                      {showConfirmPassword ? <VisibilityOff /> : <Visibility />}
                    </IconButton>
                  </InputAdornment>
                ),
              }}
              sx={{ mb: 2 }}
            />

            <TextField
              fullWidth
              label="昵称（可选）"
              name="nickname"
              value={formData.nickname}
              onChange={handleChange}
              margin="normal"
              disabled={loading}
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
                      objectFit: 'cover',
                    }}
                  />
                )}
                <IconButton
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
                </IconButton>
              </Box>
            </Box>

            {/* 注册按钮 */}
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
                '注 册'
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
            已有账号？
            <Button
              size="small"
              sx={{ ml: 1 }}
              onClick={() => navigate('/login')}
            >
              立即登录
            </Button>
          </Typography>
        </Paper>
      </Container>
    </Box>
  );
};

export default Register;
