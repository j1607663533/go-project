package controllers

import (
	"fmt"
	"gin-backend/models"
	"gin-backend/services"
	"gin-backend/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FileController struct {
	fileService services.FileService
}

func NewFileController(fileService services.FileService) *FileController {
	return &FileController{fileService: fileService}
}

// Upload 上传文件
func (ctrl *FileController) Upload(c *gin.Context) {
	fmt.Println("收到上传请求")
	file, err := c.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "获取上传文件失败")
		return
	}

	// 创建上传目录
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, os.ModePerm)
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(uploadDir, filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, dst); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "保存文件失败")
		return
	}

	// 获取当前用户信息
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")

	// 保存到数据库
	fileModel := &models.File{
		Name:     file.Filename,
		Path:     dst,
		Size:     file.Size,
		Ext:      ext,
		Type:     c.PostForm("type"),
		UserID:   userID.(uint),
		Username: username.(string),
	}

	if err := ctrl.fileService.UploadFile(fileModel); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "保存文件信息到数据库失败")
		return
	}

	utils.SuccessResponse(c, fileModel)
}

// List 获取文件列表
func (ctrl *FileController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	name := c.Query("name")

	files, total, err := ctrl.fileService.GetFileList(page, pageSize, name)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取文件列表失败")
		return
	}

	utils.SuccessPaginationResponse(c, files, total, page, pageSize)
}

// Delete 删除文件
func (ctrl *FileController) Delete(c *gin.Context) {
	idStr := c.Param("id")

	fmt.Println("收到删除请求, ID: ", idStr, c)
	id, _ := strconv.ParseUint(idStr, 10, 32)

	file, err := ctrl.fileService.GetFileByID(uint(id))
	if err != nil {
		fmt.Printf("删除失败: 文件ID %d 不存在, 错误: %v\n", id, err)
		utils.ErrorResponse(c, http.StatusNotFound, "文件不存在")
		return
	}

	// 删除物理文件
	if err := os.Remove(file.Path); err != nil {
		// 如果文件不存在，可以忽略错误继续删除数据库记录
		if !os.IsNotExist(err) {
			utils.ErrorResponse(c, http.StatusInternalServerError, "删除物理文件失败")
			return
		}
	}

	// 删除数据库记录
	if err := ctrl.fileService.DeleteFile(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除数据库记录失败")
		return
	}

	utils.SuccessResponse(c, gin.H{
		"code": 200,
		"message": "success",
		"data": nil,
	})
}

// Download 下载文件
func (ctrl *FileController) Download(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	file, err := ctrl.fileService.GetFileByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "文件不存在")
		return
	}

	c.FileAttachment(file.Path, file.Name)
}
