import React, { useState, useEffect } from 'react';
import { 
  Table, 
  Button, 
  Input, 
  Space, 
  message, 
  Modal, 
  Upload, 
  Popconfirm,
  Tag,
  Card,
  Breadcrumb
} from 'antd';
import { 
  UploadOutlined, 
  SearchOutlined, 
  DeleteOutlined, 
  DownloadOutlined,
  FileOutlined,
  FileImageOutlined,
  FilePdfOutlined,
  FileWordOutlined,
  FileExcelOutlined,
  InboxOutlined
} from '@ant-design/icons';
import { getFileList, uploadFile, deleteFile } from '../api/file';
import moment from 'moment';

const { Dragger } = Upload;

const FileManager = () => {
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState([]);
  const [total, setTotal] = useState(0);
  const [queryParams, setQueryParams] = useState({
    page: 1,
    pageSize: 10,
    name: ''
  });
  const [isUploadModalVisible, setIsUploadModalVisible] = useState(false);

  useEffect(() => {
    fetchData();
  }, [queryParams]);

  const fetchData = async () => {
    setLoading(true);
    try {
      const res = await getFileList(queryParams);
      if (res) {
        setData(res.list || []);
        setTotal(res.total || 0);
      }
      
    } catch (error) {
      message.error('获取文件列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (value) => {
    setQueryParams({
      ...queryParams,
      page: 1,
      name: value
    });
  };

  const handleDelete = async (id) => {
    try {
      const res = await deleteFile(id);
      if (res.code === 0) {
        message.success('删除成功');
        fetchData();
      }
    } catch (error) {
      message.error('删除失败');
    }
  };

  const handleDownload = (id) => {
    // 获取 token
    const token = localStorage.getItem('token');
    // 创建隐藏的 a 标签下载
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';
    const downloadUrl = `${apiUrl}/api/v1/files/download/${id}?token=${token}`;
    window.open(downloadUrl, '_blank');
  };

  const getFileIcon = (ext) => {
    const e = ext.toLowerCase();
    if (['.jpg', '.jpeg', '.png', '.gif', '.webp'].includes(e)) return <FileImageOutlined style={{ color: '#52c41a' }} />;
    if (['.pdf'].includes(e)) return <FilePdfOutlined style={{ color: '#f5222d' }} />;
    if (['.doc', '.docx'].includes(e)) return <FileWordOutlined style={{ color: '#1890ff' }} />;
    if (['.xls', '.xlsx'].includes(e)) return <FileExcelOutlined style={{ color: '#52c41a' }} />;
    return <FileOutlined />;
  };

  const formatSize = (bytes) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const columns = [
    {
      title: '文件名',
      dataIndex: 'name',
      key: 'name',
      render: (text, record) => (
        <Space>
          {getFileIcon(record.ext)}
          <span>{text}</span>
        </Space>
      )
    },
    {
      title: '大小',
      dataIndex: 'size',
      key: 'size',
      render: (size) => formatSize(size)
    },
    {
      title: '上传者',
      dataIndex: 'username',
      key: 'username',
    },
    {
      title: '上传时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (time) => moment(time).format('YYYY-MM-DD HH:mm:ss')
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space size="middle">
          <Button 
            type="link" 
            icon={<DownloadOutlined />} 
            onClick={() => handleDownload(record.id)}
          >
            下载
          </Button>
          <Popconfirm
            title="确定要删除这个文件吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const uploadProps = {
    name: 'file',
    multiple: true,
    customRequest: async (options) => {
      const { onSuccess, onError, file, onProgress } = options;
      const formData = new FormData();
      formData.append('file', file);
      formData.append('type', 'file'); // 可以根据需要设置

      try {
        const res = await uploadFile(formData);
        // 如果 http.post 已经处理了成功响应脱壳，res 就是 res.data
        // 根据 src/utils/request.js 的实现，res 将是后端返回的 Data 字段
        onSuccess(res, file);
        message.success(`${file.name} 上传成功`);
        fetchData();
      } catch (error) {
        onError(error);
        message.error(`${file.name} 上传失败: ${error.message}`);
      }
    },
    onDrop(e) {
      console.log('Dropped files', e.dataTransfer.files);
    },
  };

  return (
    <div style={{ padding: '24px' }}>
      <Breadcrumb style={{ marginBottom: '16px' }}>
        <Breadcrumb.Item>首页</Breadcrumb.Item>
        <Breadcrumb.Item>文件管理</Breadcrumb.Item>
      </Breadcrumb>

      <Card>
        <div style={{ marginBottom: '16px', display: 'flex', justifyContent: 'space-between' }}>
          <Space>
            <Input.Search
              placeholder="搜索文件名"
              onSearch={handleSearch}
              style={{ width: 250 }}
              allowClear
            />
          </Space>
          <Button 
            type="primary" 
            icon={<UploadOutlined />} 
            onClick={() => setIsUploadModalVisible(true)}
          >
            上传文件
          </Button>
        </div>  
        <Table
          columns={columns}
          dataSource={data}
          loading={loading}
          rowKey="id"
          pagination={{
            current: queryParams.page,
            pageSize: queryParams.pageSize,
            total: total,
            showSizeChanger: true,
            showTotal: (total) => `共 ${total} 条`,
            onChange: (page, pageSize) => {
              setQueryParams({ ...queryParams, page, pageSize });
            }
          }}
        />
      </Card>

      <Modal
        title="上传文件"
        open={isUploadModalVisible}
        onCancel={() => setIsUploadModalVisible(false)}
        footer={null}
        width={600}
      >
        <Dragger {...uploadProps}>
          <p className="ant-upload-drag-icon">
            <InboxOutlined />
          </p>
          <p className="ant-upload-text">点击或将文件拖拽到此区域上传</p>
          <p className="ant-upload-hint">
            支持单个或批量上传，严禁上传含有违法违规内容的文件。
          </p>
        </Dragger>
      </Modal>
    </div>
  );
};

export default FileManager;
