package services

import (
	"gin-backend/models"
	"gin-backend/repositories"
)

type FileService interface {
	UploadFile(file *models.File) error
	GetFileList(page, pageSize int, name string) ([]models.File, int64, error)
	GetFileByID(id uint) (*models.File, error)
	DeleteFile(id uint) error
}

type fileService struct {
	fileRepo repositories.FileRepository
}

func NewFileService(fileRepo repositories.FileRepository) FileService {
	return &fileService{fileRepo: fileRepo}
}

func (s *fileService) UploadFile(file *models.File) error {
	return s.fileRepo.Create(file)
}

func (s *fileService) GetFileList(page, pageSize int, name string) ([]models.File, int64, error) {
	offset := (page - 1) * pageSize
	return s.fileRepo.FindAll(offset, pageSize, name)
}

func (s *fileService) GetFileByID(id uint) (*models.File, error) {
	return s.fileRepo.FindByID(id)
}

func (s *fileService) DeleteFile(id uint) error {
	return s.fileRepo.Delete(id)
}
