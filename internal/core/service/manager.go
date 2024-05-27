package service

import (
	"github.com/khivuksergey/portmonetka.category/internal/core/port/repository"
	"github.com/khivuksergey/portmonetka.category/internal/core/port/service"
	"github.com/khivuksergey/portmonetka.category/internal/core/service/category"
)

func NewServiceManager(repositoryManager *repository.Manager) *service.Manager {
	return &service.Manager{
		Category: category.NewCategoryService(repositoryManager),
	}
}
