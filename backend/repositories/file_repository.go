package repositories

import (
	"fmt"
	"gin-backend/models"

	"gorm.io/gorm"
)

type FileRepository interface {
	Create(file *models.File) error
	FindAll(offset, limit int, name string) ([]models.File, int64, error)
	FindByID(id uint) (*models.File, error)
	Delete(id uint) error
}

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) Create(file *models.File) error {
	return r.db.Create(file).Error
}

func (r *fileRepository) FindAll(offset, limit int, name string) ([]models.File, int64, error) {
	var files []models.File
	var total int64

	query := r.db.Model(&models.File{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&files).Error
	return files, total, err
}

func (r *fileRepository) FindByID(id uint) (*models.File, error) {

	fmt.Println("收到查询请求, ID: ", id)
	var file models.File
	err := r.db.First(&file, id).Error
	return &file, err
}

func (r *fileRepository) Delete(id uint) error {
	return r.db.Delete(&models.File{}, id).Error
}
