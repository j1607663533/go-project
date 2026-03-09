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
  QrCode,
  Laptop,
} from '@mui/icons-material';
import WeChatIcon from '@mui/icons-material/WhatsApp'; // Using a similar icon or we can use custom
import { getCaptcha, login, isAuthenticated, getWeChatQrCode, checkWeChatLoginStatus } from '../api/auth';
import Welcome3D from '../components/Welcome3D';

import loginBg from '../assets/login-bg.png';

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
  const [loginMethod, setLoginMethod] = useState('account'); // 'account' or 'wechat'
  const [wechatData, setWechatData] = useState({ qr_url: '', scene_id: '', status: 'IDLE' }); // IDLE, SCANNING, SUCCESS, EXPIRED
  const pollingRef = useRef(null);
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

    return () => {
      if (pollingRef.current) clearInterval(pollingRef.current);
    };
  }, [navigate]);

  // 微信扫码轮询
  useEffect(() => {
    if (loginMethod === 'wechat' && wechatData.status === 'IDLE') {
      fetchWeChatQr();
    } else if (loginMethod !== 'wechat') {
      stopPolling();
    }
  }, [loginMethod, wechatData.status]);

  const fetchWeChatQr = async () => {
    try {
      const data = await getWeChatQrCode();
      setWechatData({ ...data, status: 'SCANNING' });
      startPolling(data.scene_id);
    } catch (err) {
      setError('获取二维码失败');
    }
  };

  const startPolling = (sceneId) => {
    if (pollingRef.current) clearInterval(pollingRef.current);
    
    pollingRef.current = setInterval(async () => {
      try {
        const result = await checkWeChatLoginStatus(sceneId);
        if (result.status === 'SUCCESS') {
          stopPolling();
          setWechatData(prev => ({ ...prev, status: 'SUCCESS' }));
          
          // 模拟登录成功逻辑
          localStorage.setItem("token", result.token);
          localStorage.setItem("user", JSON.stringify(result.user));
          localStorage.setItem("menus", JSON.stringify(result.menus));
          
          setTimeout(() => navigate('/'), 1000);
        } else if (result.status === 'EXPIRED') {
          stopPolling();
          setWechatData(prev => ({ ...prev, status: 'EXPIRED' }));
        }
      } catch (err) {
        console.error('轮询出错', err);
      }
    }, 2000); // 2秒轮询一次
  };

  const stopPolling = () => {
    if (pollingRef.current) {
      clearInterval(pollingRef.current);
      pollingRef.current = null;
    }
  };

  const refreshWeChatQr = () => {
    setWechatData({ qr_url: '', scene_id: '', status: 'IDLE' });
  };

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
        justifyContent: 'center',
        backgroundImage: `linear-gradient(rgba(15, 23, 42, 0.6), rgba(15, 23, 42, 0.6)), url(${loginBg})`,
        backgroundSize: 'cover',
        backgroundPosition: 'center',
        backgroundAttachment: 'fixed',
      }}
    >
      <Container maxWidth="sm">
        <Paper
          elevation={0}
          className="glass-morphism"
          sx={{
            p: 5,
            borderRadius: '24px',
            border: '1px solid rgba(255, 255, 255, 0.1) !important',
            boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.5) !important',
          }}
        >
          {/* Logo */}
          <Box sx={{ display: 'flex', justifyContent: 'center', mb: 4 }}>
            <Avatar
              sx={{
                width: 80,
                height: 80,
                background: 'linear-gradient(135deg, #6366f1 0%, #a855f7 100%)',
                boxShadow: '0 8px 16px rgba(99, 102, 241, 0.4)',
              }}
            >
              <LockOutlined sx={{ fontSize: 40 }} />
            </Avatar>
          </Box>

          {/* 标题 */}
          <Welcome3D />

          <Typography
            variant="body1"
            align="center"
            sx={{ 
              mb: 4, 
              color: 'rgba(255, 255, 255, 0.7)',
              fontFamily: 'var(--font-heading)',
              letterSpacing: 1,
              fontWeight: 300
            }}
          >
            DIGITAL THRESHOLD ALIGNMENT
          </Typography>

          {/* 账户/微信 切换 */}
          <Box sx={{ display: 'flex', justifyContent: 'center', mb: 4, gap: 2 }}>
            <Button
              startIcon={<Laptop />}
              variant={loginMethod === 'account' ? 'contained' : 'outlined'}
              onClick={() => setLoginMethod('account')}
              sx={{ 
                borderRadius: '12px',
                borderColor: 'rgba(255,255,255,0.2)',
                color: '#fff',
                background: loginMethod === 'account' ? 'linear-gradient(135deg, #6366f1 0%, #a855f7 100%)' : 'transparent',
                '&:hover': {
                  borderColor: '#6366f1',
                  background: loginMethod === 'account' ? 'linear-gradient(135deg, #4f46e5 0%, #9333ea 100%)' : 'rgba(255,255,255,0.05)',
                }
              }}
            >
              账号登录
            </Button>
            <Button
              startIcon={<QrCode />}
              variant={loginMethod === 'wechat' ? 'contained' : 'outlined'}
              onClick={() => setLoginMethod('wechat')}
              sx={{ 
                borderRadius: '12px',
                borderColor: 'rgba(255,255,255,0.2)',
                color: '#fff',
                background: loginMethod === 'wechat' ? 'linear-gradient(135deg, #6366f1 0%, #a855f7 100%)' : 'transparent',
                '&:hover': {
                  borderColor: '#6366f1',
                  background: loginMethod === 'wechat' ? 'linear-gradient(135deg, #4f46e5 0%, #9333ea 100%)' : 'rgba(255,255,255,0.05)',
                }
              }}
            >
              微信扫码
            </Button>
          </Box>

          {/* 错误提示 */}
          {error && (
            <Alert 
              severity="error" 
              sx={{ 
                mb: 3, 
                borderRadius: '12px',
                background: 'rgba(244, 63, 94, 0.1)',
                color: '#f43f5e',
                border: '1px solid rgba(244, 63, 94, 0.2)'
              }}
            >
              {error}
            </Alert>
          )}

          {/* 登录表单 */}
          {loginMethod === 'account' ? (
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
                InputProps={{
                  sx: { color: '#fff' }
                }}
                InputLabelProps={{
                  sx: { color: 'rgba(255,255,255,0.5)' }
                }}
                sx={{ 
                  mb: 2,
                  '& .MuiOutlinedInput-root': {
                    '& fieldset': { borderColor: 'rgba(255,255,255,0.2)' },
                    '&:hover fieldset': { borderColor: '#6366f1' },
                    '&.Mui-focused fieldset': { borderColor: '#6366f1' },
                  }
                }}
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
                  sx: { color: '#fff' },
                  endAdornment: (
                    <InputAdornment position="end">
                      <IconButton
                        onClick={() => setShowPassword(!showPassword)}
                        edge="end"
                        sx={{ color: 'rgba(255,255,255,0.5)' }}
                      >
                        {showPassword ? <VisibilityOff /> : <Visibility />}
                      </IconButton>
                    </InputAdornment>
                  ),
                }}
                InputLabelProps={{
                  sx: { color: 'rgba(255,255,255,0.5)' }
                }}
                sx={{ 
                  mb: 2,
                  '& .MuiOutlinedInput-root': {
                    '& fieldset': { borderColor: 'rgba(255,255,255,0.2)' },
                    '&:hover fieldset': { borderColor: '#6366f1' },
                    '&.Mui-focused fieldset': { borderColor: '#6366f1' },
                  }
                }}
              />

              {/* 验证码 */}
              <Box sx={{ display: 'flex', gap: 2, mb: 4 }}>
                <TextField
                  label="验证码"
                  name="captcha"
                  value={formData.captcha}
                  onChange={handleChange}
                  required
                  disabled={loading}
                  inputProps={{ maxLength: 6, sx: { color: '#fff' } }}
                  InputLabelProps={{
                    sx: { color: 'rgba(255,255,255,0.5)' }
                  }}
                  sx={{ 
                    flex: 1,
                    '& .MuiOutlinedInput-root': {
                      '& fieldset': { borderColor: 'rgba(255,255,255,0.2)' },
                      '&:hover fieldset': { borderColor: '#6366f1' },
                      '&.Mui-focused fieldset': { borderColor: '#6366f1' },
                    }
                  }}
                />
                
                <Box
                  sx={{
                    position: 'relative',
                    width: 120,
                    height: 56,
                    cursor: 'pointer',
                    borderRadius: '12px',
                    overflow: 'hidden',
                    border: '1px solid rgba(255,255,255,0.2)',
                    transition: 'border-color 0.2s',
                    '&:hover': {
                      borderColor: '#6366f1',
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
                        filter: 'contrast(1.2) brightness(1.1)',
                      }}
                    />
                  )}
                </Box>
              </Box>

              <Button
                type="submit"
                fullWidth
                variant="contained"
                size="large"
                disabled={loading}
                sx={{
                  py: 1.8,
                  fontSize: '1rem',
                  fontWeight: 700,
                  letterSpacing: 1,
                  borderRadius: '14px',
                  background: 'linear-gradient(135deg, #6366f1 0%, #a855f7 100%)',
                  boxShadow: '0 10px 20px -5px rgba(99, 102, 241, 0.5)',
                  '&:hover': {
                    background: 'linear-gradient(135deg, #4f46e5 0%, #9333ea 100%)',
                    transform: 'translateY(-2px)',
                    boxShadow: '0 15px 25px -5px rgba(99, 102, 241, 0.6)',
                  },
                  transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
                }}
              >
                {loading ? (
                  <CircularProgress size={24} color="inherit" />
                ) : (
                  'ACCESS SYSTEM'
                )}
              </Button>
            </form>
          ) : (
            /* 微信扫码登录逻辑 */
            <Box sx={{ textAlign: 'center', py: 2 }}>
              <Box
                sx={{
                  position: 'relative',
                  width: 220,
                  height: 220,
                  margin: '0 auto',
                  p: 1,
                  border: '1px solid rgba(255,255,255,0.1)',
                  borderRadius: 3,
                  bgcolor: 'rgba(255,255,255,0.05)',
                  backdropFilter: 'blur(5px)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                }}
              >
                {wechatData.qr_url ? (
                  <img
                    src={wechatData.qr_url}
                    alt="WeChat QR Code"
                    style={{ 
                      width: '100%', 
                      height: '100%',
                      borderRadius: '8px',
                      opacity: wechatData.status === 'EXPIRED' ? 0.3 : 1 
                    }}
                  />
                ) : (
                  <CircularProgress sx={{ color: '#6366f1' }} />
                )}
                
                {/* 状态重叠层 */}
                {wechatData.status === 'EXPIRED' && (
                  <Box sx={{ 
                    position: 'absolute', 
                    top: 0, left: 0, right: 0, bottom: 0, 
                    bgcolor: 'rgba(15, 23, 42, 0.8)', 
                    display: 'flex', 
                    flexDirection: 'column', 
                    alignItems: 'center', 
                    justifyContent: 'center',
                    zIndex: 2,
                    borderRadius: 3
                  }}>
                    <Typography variant="body2" sx={{ mb: 2, fontWeight: 'bold', color: '#fff' }}>二维码已失效</Typography>
                    <Button 
                      size="small" 
                      variant="contained" 
                      startIcon={<Refresh />}
                      onClick={refreshWeChatQr}
                      sx={{ background: '#6366f1' }}
                    >
                      点击刷新
                    </Button>
                  </Box>
                )}

                {wechatData.status === 'SUCCESS' && (
                  <Box sx={{ 
                    position: 'absolute', 
                    top: 0, left: 0, right: 0, bottom: 0, 
                    bgcolor: 'rgba(15, 23, 42, 0.9)', 
                    display: 'flex', 
                    flexDirection: 'column', 
                    alignItems: 'center', 
                    justifyContent: 'center',
                    zIndex: 2,
                    borderRadius: 3
                  }}>
                    <Avatar sx={{ bgcolor: '#10b981', mb: 2 }}>
                      <Refresh />
                    </Avatar>
                    <Typography variant="body2" sx={{ fontWeight: 'bold', color: '#10b981' }}>
                      验证通过
                    </Typography>
                  </Box>
                )}
              </Box>
              
              <Typography variant="body2" sx={{ mt: 4, display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 1, color: 'rgba(255,255,255,0.6)' }}>
                请使用 <span style={{ color: '#10b981', fontWeight: 'bold' }}>微信</span> 扫码对齐数字身份
              </Typography>
            </Box>
          )}

          {/* 底部提示 */}
          <Typography
            variant="body2"
            align="center"
            sx={{ mt: 4, color: 'rgba(255,255,255,0.4)' }}
          >
            新成员？
            <Button
              size="small"
              sx={{ ml: 1, color: '#6366f1', fontWeight: 600 }}
              onClick={() => navigate('/register')}
            >
              启动注册序列
            </Button>
          </Typography>
        </Paper>
      </Container>
    </Box>
  );
};

export default Login;
