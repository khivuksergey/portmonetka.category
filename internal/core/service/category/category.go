package category

import (
	serviceerror "github.com/khivuksergey/portmonetka.category/error"
	"github.com/khivuksergey/portmonetka.category/internal/adapter/storage/entity"
	"github.com/khivuksergey/portmonetka.category/internal/core/port/repository"
	"github.com/khivuksergey/portmonetka.category/internal/core/port/service"
	"github.com/khivuksergey/portmonetka.category/internal/model"
)

type category struct {
	categoryRepository repository.CategoryRepository
}

func NewCategoryService(repositoryManager *repository.Manager) service.CategoryService {
	return &category{categoryRepository: repositoryManager.Category}
}

func (c *category) GetCategoriesByUserId(userId uint64) ([]entity.Category, error) {
	return c.categoryRepository.GetCategoriesByUserId(userId)
}

func (c *category) CreateCategory(categoryCreateDTO model.CategoryCreateDTO) (*entity.Category, error) {
	if c.categoryRepository.ExistsWithName(categoryCreateDTO.UserId, categoryCreateDTO.Name) {
		return nil, serviceerror.CategoryAlreadyExists
	}
	return c.categoryRepository.CreateCategory(&entity.Category{
		UserId:      categoryCreateDTO.UserId,
		Name:        categoryCreateDTO.Name,
		Description: categoryCreateDTO.Description,
		Type:        categoryCreateDTO.Type,
	})
}

func (c *category) UpdateCategory(categoryUpdateDTO model.CategoryUpdateDTO) (*entity.Category, error) {
	categoryToUpdate, err := c.categoryRepository.GetCategoryById(categoryUpdateDTO.Id)
	if err != nil {
		return nil, serviceerror.CategoryDoesntExist
	}
	err = c.validateUpdateCategoryAttributes(categoryToUpdate, categoryUpdateDTO)
	if err != nil {
		return nil, err
	}
	return c.categoryRepository.UpdateCategory(categoryToUpdate)
}

func (c *category) DeleteCategory(categoryDeleteDTO model.CategoryDeleteDTO) error {
	if !c.categoryRepository.CategoryBelongsToUser(categoryDeleteDTO.Id, categoryDeleteDTO.UserId) {
		return serviceerror.CategoryDoesntBelongToUser
	}
	return c.categoryRepository.DeleteCategory(categoryDeleteDTO.Id)
}

// TODO move attributes validation to validator
func (c *category) validateUpdateCategoryAttributes(category *entity.Category, categoryUpdateDTO model.CategoryUpdateDTO) error {
	if categoryUpdateDTO.Name == nil && categoryUpdateDTO.Description == nil {
		return serviceerror.AtLeastOneFieldIsRequired
	}
	if categoryUpdateDTO.Name != nil {
		if len(*categoryUpdateDTO.Name) < 3 || len(*categoryUpdateDTO.Name) > 128 {
			return serviceerror.CategoryNameLengthError
		}
		if c.categoryRepository.ExistsWithName(categoryUpdateDTO.UserId, *categoryUpdateDTO.Name) {
			return serviceerror.CategoryAlreadyExists
		}
		category.Name = *categoryUpdateDTO.Name
	}
	if categoryUpdateDTO.Description != nil {
		if len(*categoryUpdateDTO.Description) > 256 {
			return serviceerror.CategoryDescriptionLengthError
		}
		category.Description = *categoryUpdateDTO.Description
	}
	return nil
}
