import { useState, useEffect, useCallback } from 'react';
import {
  Box,
  Container,
  Paper,
  Typography,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TablePagination,
  IconButton,
  Chip,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  MenuItem,
  CircularProgress,
  Tooltip,
  Stack,
  InputAdornment,
} from '@mui/material';
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Refresh as RefreshIcon,
  ShoppingBag as ShoppingBagIcon,
  Search as SearchIcon,
  FilterList as FilterIcon,
} from '@mui/icons-material';
import { getOrdersByPage, createOrder, updateOrder, deleteOrder } from '../api/order';

const Orders = () => {
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState({ data: [], total: 0, page: 1, page_size: 10 });
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  
  // 临时筛选状态（输入框中的值）
  const [filters, setFilters] = useState({
    product_id: '',
    status: ''
  });

  // 实际应用的搜索参数
  const [searchParams, setSearchParams] = useState({
    product_id: '',
    status: ''
  });

  // 对话框状态
  const [open, setOpen] = useState(false);
  const [editingOrder, setEditingOrder] = useState(null);
  const [formData, setFormData] = useState({
    product_id: '',
    quantity: 1,
    total: 0,
    status: 'pending',
    payment_id: 0
  });

  const fetchOrders = useCallback(async () => {
    try {
      setLoading(true);
      const res = await getOrdersByPage(
        page + 1, 
        rowsPerPage, 
        searchParams.product_id || null,
        searchParams.status || null
      );
      setData(res);
    } catch (error) {
      console.error('加载订单失败:', error);
    } finally {
      setLoading(false);
    }
  }, [page, rowsPerPage, searchParams]);

  useEffect(() => {
    fetchOrders();
  }, [fetchOrders]);

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const handleFilterChange = (e) => {
    const { name, value } = e.target;
    setFilters(prev => ({ ...prev, [name]: value }));
  };

  // 触发搜索
  const handleSearch = () => {
    setSearchParams({ ...filters });
    setPage(0); // 重置到第一页
  };

  // 重置筛选
  const handleResetFilters = () => {
    const defaultFilters = { product_id: '', status: '' };
    setFilters(defaultFilters);
    setSearchParams(defaultFilters);
    setPage(0);
  };

  const handleOpenDialog = (order = null) => {
    if (order) {
      setEditingOrder(order);
      setFormData({
        product_id: order.product_id,
        quantity: order.quantity,
        total: order.total,
        status: order.status,
        payment_id: order.payment_id || 0
      });
    } else {
      setEditingOrder(null);
      setFormData({
        product_id: '',
        quantity: 1,
        total: 0,
        status: 'pending',
        payment_id: 0
      });
    }
    setOpen(true);
  };

  const handleCloseDialog = () => {
    setOpen(false);
  };

  const handleSubmit = async () => {
    try {
      if (!formData.product_id) {
        alert('请填写商品ID');
        return;
      }

      const body = {
        ...formData,
        product_id: Number(formData.product_id),
        quantity: Number(formData.quantity),
        total: Number(formData.total),
        payment_id: Number(formData.payment_id),
      };

      if (editingOrder) {
        await updateOrder(editingOrder.id, body);
      } else {
        await createOrder(body);
      }
      handleCloseDialog();
      fetchOrders();
    } catch (error) {
      console.error('保存失败:', error);
      alert(error.message);
    }
  };

  const handleDelete = async (id) => {
    if (window.confirm('确定要删除这个订单吗？')) {
      try {
        await deleteOrder(id);
        fetchOrders();
      } catch (error) {
        console.error('删除失败:', error);
        alert(error.message);
      }
    }
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 'paid': return 'success';
      case 'pending': return 'warning';
      case 'cancelled': return 'error';
      default: return 'default';
    }
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleString();
  };

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      {/* 标题栏 */}
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4 }}>
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
          <ShoppingBagIcon sx={{ mr: 1, color: 'primary.main', fontSize: 32 }} />
          <Typography variant="h4" sx={{ fontWeight: 700 }}>
            订单管理
          </Typography>
        </Box>
        <Box>
          <Tooltip title="刷新">
            <IconButton onClick={fetchOrders} sx={{ mr: 1 }}>
              <RefreshIcon />
            </IconButton>
          </Tooltip>
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={() => handleOpenDialog()}
            sx={{
              borderRadius: 2,
              textTransform: 'none',
              px: 3,
              boxShadow: '0 4px 12px rgba(102, 126, 234, 0.4)',
            }}
          >
            新订单
          </Button>
        </Box>
      </Box>

      {/* 筛选栏 */}
      <Paper elevation={0} sx={{ p: 3, mb: 3, borderRadius: 3, bgcolor: 'background.paper', border: '1px solid', borderColor: 'divider' }}>
        <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} alignItems="center">
          <TextField
            label="商品 ID (模糊)"
            name="product_id"
            value={filters.product_id}
            onChange={handleFilterChange}
            size="small"
            sx={{ minWidth: 200 }}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <SearchIcon fontSize="small" color="action" />
                </InputAdornment>
              ),
            }}
          />
          <TextField
            select
            label="订单状态"
            name="status"
            value={filters.status}
            onChange={handleFilterChange}
            size="small"
            sx={{ minWidth: 200 }}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <FilterIcon fontSize="small" color="action" />
                </InputAdornment>
              ),
            }}
          >
            <MenuItem value="">全部状态</MenuItem>
            <MenuItem value="pending">待支付 (Pending)</MenuItem>
            <MenuItem value="paid">已支付 (Paid)</MenuItem>
            <MenuItem value="cancelled">已取消 (Cancelled)</MenuItem>
          </TextField>
          <Button 
            variant="contained" 
            onClick={handleSearch}
            startIcon={<SearchIcon />}
            sx={{ borderRadius: 2, textTransform: 'none' }}
          >
            搜索
          </Button>
          <Button 
            variant="outlined" 
            onClick={handleResetFilters}
            sx={{ borderRadius: 2, textTransform: 'none' }}
          >
            重置
          </Button>
        </Stack>
      </Paper>

      {/* 订单表格 */}
      <Paper elevation={4} sx={{ borderRadius: 3, overflow: 'hidden' }}>
        <TableContainer>
          <Table>
            <TableHead sx={{ bgcolor: 'rgba(102, 126, 234, 0.05)' }}>
              <TableRow>
                <TableCell sx={{ fontWeight: 700 }}>ID</TableCell>
                <TableCell sx={{ fontWeight: 700 }}>用户 ID</TableCell>
                <TableCell sx={{ fontWeight: 700 }}>商品 ID</TableCell>
                <TableCell sx={{ fontWeight: 700 }}>数量</TableCell>
                <TableCell sx={{ fontWeight: 700 }}>总金额</TableCell>
                <TableCell sx={{ fontWeight: 700 }}>状态</TableCell>
                <TableCell sx={{ fontWeight: 700 }}>创建时间</TableCell>
                <TableCell align="right" sx={{ fontWeight: 700 }}>操作</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={8} align="center" sx={{ py: 8 }}>
                    <CircularProgress />
                  </TableCell>
                </TableRow>
              ) : data.data.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={8} align="center" sx={{ py: 8 }}>
                    <Typography color="text.secondary">无符合条件的订单</Typography>
                  </TableCell>
                </TableRow>
              ) : (
                data.data.map((order) => (
                  <TableRow key={order.id} hover>
                    <TableCell>{order.id}</TableCell>
                    <TableCell>{order.user_id}</TableCell>
                    <TableCell>{order.product_id}</TableCell>
                    <TableCell>{order.quantity}</TableCell>
                    <TableCell>¥{order.total}</TableCell>
                    <TableCell>
                      <Chip
                        label={order.status.toUpperCase()}
                        size="small"
                        color={getStatusColor(order.status)}
                        variant="outlined"
                      />
                    </TableCell>
                    <TableCell>{formatDate(order.created_at)}</TableCell>
                    <TableCell align="right">
                      <IconButton onClick={() => handleOpenDialog(order)} color="primary" size="small">
                        <EditIcon fontSize="small" />
                      </IconButton>
                      <IconButton onClick={() => handleDelete(order.id)} color="error" size="small">
                        <DeleteIcon fontSize="small" />
                      </IconButton>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </TableContainer>
        <TablePagination
          component="div"
          count={data.total}
          page={page}
          onPageChange={handleChangePage}
          rowsPerPage={rowsPerPage}
          onRowsPerPageChange={handleChangeRowsPerPage}
          labelRowsPerPage="每页显示"
          labelDisplayedRows={({ from, to, count }) => `${from}-${to} 共 ${count}`}
        />
      </Paper>

      {/* 编辑/创建对话框 */}
      <Dialog open={open} onClose={handleCloseDialog} fullWidth maxWidth="sm">
        <DialogTitle sx={{ fontWeight: 700 }}>
          {editingOrder ? '编辑订单' : '创建新订单'}
        </DialogTitle>
        <DialogContent dividers>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, pt: 1 }}>
            <TextField
              label="商品 ID"
              type="number"
              fullWidth
              value={formData.product_id}
              onChange={(e) => setFormData({ ...formData, product_id: e.target.value })}
              required
            />
            <Box sx={{ display: 'flex', gap: 2 }}>
              <TextField
                label="数量"
                type="number"
                fullWidth
                value={formData.quantity}
                onChange={(e) => setFormData({ ...formData, quantity: e.target.value })}
              />
              <TextField
                label="总金额"
                type="number"
                fullWidth
                value={formData.total}
                onChange={(e) => setFormData({ ...formData, total: e.target.value })}
              />
            </Box>
            <TextField
              select
              label="订单状态"
              fullWidth
              value={formData.status}
              onChange={(e) => setFormData({ ...formData, status: e.target.value })}
            >
              <MenuItem value="pending">待支付 (Pending)</MenuItem>
              <MenuItem value="paid">已支付 (Paid)</MenuItem>
              <MenuItem value="cancelled">已取消 (Cancelled)</MenuItem>
            </TextField>
            <TextField
              label="支付 ID"
              type="number"
              fullWidth
              value={formData.payment_id}
              onChange={(e) => setFormData({ ...formData, payment_id: e.target.value })}
            />
          </Box>
        </DialogContent>
        <DialogActions sx={{ p: 2, px: 3 }}>
          <Button onClick={handleCloseDialog} color="inherit">取消</Button>
          <Button onClick={handleSubmit} variant="contained">提交</Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
};

export default Orders;
