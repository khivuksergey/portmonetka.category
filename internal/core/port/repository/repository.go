package repository

import (
	"github.com/khivuksergey/portmonetka.category/internal/adapter/storage/entity"
)

type Manager struct {
	Category CategoryRepository
}

//go:generate mockgen -source=repository.go -destination=../../../adapter/storage/gorm/repo/mock/mock_repository.go -package=mock
type CategoryRepository interface {
	ExistsWithName(userId uint64, name string) bool
	CategoryBelongsToUser(id, userId uint64) bool
	GetCategoryById(id uint64) (*entity.Category, error)
	GetCategoriesByUserId(userId uint64) ([]entity.Category, error)
	CreateCategory(category *entity.Category) (*entity.Category, error)
	UpdateCategory(category *entity.Category) (*entity.Category, error)
	DeleteCategory(id uint64) error
}
