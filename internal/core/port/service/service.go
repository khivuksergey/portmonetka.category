package service

import (
	"github.com/khivuksergey/portmonetka.category/internal/adapter/storage/entity"
	"github.com/khivuksergey/portmonetka.category/internal/model"
)

type Manager struct {
	Category CategoryService
}

type CategoryService interface {
	GetCategoriesByUserId(userId uint64) ([]entity.Category, error)
	CreateCategory(categoryCreateDTO model.CategoryCreateDTO) (*entity.Category, error)
	UpdateCategory(categoryUpdateDTO model.CategoryUpdateDTO) (*entity.Category, error)
	DeleteCategory(categoryDeleteDTO model.CategoryDeleteDTO) error
}
