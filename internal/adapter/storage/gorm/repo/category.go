package repo

import (
	"github.com/khivuksergey/portmonetka.category/internal/adapter/storage/entity"
	"github.com/khivuksergey/portmonetka.category/internal/core/port/repository"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db        *gorm.DB
	tableName string
}

func NewCategoryRepository(db *gorm.DB) repository.CategoryRepository {
	return &categoryRepository{db: db, tableName: entity.Category{}.TableName()}
}

func (w *categoryRepository) ExistsWithName(userId uint64, name string) bool {
	var count int64
	w.db.Model(&entity.Category{}).Where("user_id = ? AND name = ?", userId, name).Count(&count)
	return count == 1
}

func (w *categoryRepository) GetCategoryById(id uint64) (*entity.Category, error) {
	category := &entity.Category{}
	result := w.db.First(category, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return category, nil
}

func (w *categoryRepository) CategoryBelongsToUser(id, userId uint64) bool {
	category, err := w.GetCategoryById(id)
	if err != nil || category == nil {
		return false
	}
	return category.UserId == userId
}

func (w *categoryRepository) GetCategoriesByUserId(userId uint64) ([]entity.Category, error) {
	var categories []entity.Category
	result := w.db.
		Where("user_id = ?", userId).
		Order("updated_at desc").
		Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}
	return categories, nil
}

func (w *categoryRepository) CreateCategory(category *entity.Category) (*entity.Category, error) {
	if err := w.db.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (w *categoryRepository) UpdateCategory(category *entity.Category) (*entity.Category, error) {
	err := w.db.Save(category).Error
	return category, err
}

func (w *categoryRepository) DeleteCategory(id uint64) error {
	return w.db.Delete(&entity.Category{}, id).Error
}
