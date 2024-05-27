package category

import (
	serviceerror "github.com/khivuksergey/portmonetka.category/error"
	"github.com/khivuksergey/portmonetka.category/internal/adapter/storage/entity"
	"github.com/khivuksergey/portmonetka.category/internal/adapter/storage/gorm/repo/mock"
	"github.com/khivuksergey/portmonetka.category/internal/core/port/repository"
	"github.com/khivuksergey/portmonetka.category/internal/core/service/category"
	"github.com/khivuksergey/portmonetka.category/internal/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

//TODO add tests to check struct fields validation

func TestGetCategoriesByUserId_Success(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockCategoryRepository := mock.NewMockCategoryRepository(ctl)
	mockManager := &repository.Manager{
		Category: mockCategoryRepository,
	}

	categoryService := category.NewCategoryService(mockManager)

	userId := uint64(1)
	expectedCategories := []entity.Category{
		{
			Id:          1,
			UserId:      userId,
			Name:        "Test Category 1",
			Description: "Description 1",
			Type:        "INCOME",
		},
		{
			Id:          2,
			UserId:      userId,
			Name:        "Test Category 2",
			Description: "Description 2",
			Type:        "EXPENSE",
		},
	}

	mockCategoryRepository.
		EXPECT().
		GetCategoriesByUserId(userId).
		Times(1).
		Return(expectedCategories, nil)

	actualCategories, err := categoryService.GetCategoriesByUserId(userId)

	assert.NoError(t, err)
	assert.NotNil(t, actualCategories)
	assert.Equal(t, expectedCategories, actualCategories)
}

func TestCreateCategory_Success(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockCategoryRepository := mock.NewMockCategoryRepository(ctl)
	mockManager := &repository.Manager{
		Category: mockCategoryRepository,
	}

	categoryService := category.NewCategoryService(mockManager)

	categoryCreateDTO := &model.CategoryCreateDTO{
		UserId:      1,
		Name:        "Test category",
		Description: "A test category",
		Type:        "INCOME", //this is validated on http handler level
	}

	expectedCategory := &entity.Category{
		UserId:      categoryCreateDTO.UserId,
		Name:        categoryCreateDTO.Name,
		Description: categoryCreateDTO.Description,
		Type:        categoryCreateDTO.Type,
	}

	mockCategoryRepository.
		EXPECT().
		ExistsWithName(categoryCreateDTO.UserId, categoryCreateDTO.Name).
		Times(1).
		Return(false)

	mockCategoryRepository.
		EXPECT().
		CreateCategory(expectedCategory).
		Times(1).
		Return(expectedCategory, nil)

	createdCategory, err := categoryService.CreateCategory(*categoryCreateDTO)

	assert.NoError(t, err)
	assert.NotNil(t, createdCategory)
	assert.Equal(t, createdCategory, expectedCategory)
}

func TestCreateCategory_DuplicateName_Error(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockCategoryRepository := mock.NewMockCategoryRepository(ctl)
	mockManager := &repository.Manager{
		Category: mockCategoryRepository,
	}

	categoryService := category.NewCategoryService(mockManager)

	categoryCreateDTO := &model.CategoryCreateDTO{
		UserId:      1,
		Name:        "Duplicate category",
		Description: "A test category with duplicate name",
		Type:        "expense",
	}

	mockCategoryRepository.
		EXPECT().
		ExistsWithName(categoryCreateDTO.UserId, categoryCreateDTO.Name).
		Times(1).
		Return(true)

	createdCategory, err := categoryService.CreateCategory(*categoryCreateDTO)

	assert.Error(t, err)
	assert.Nil(t, createdCategory)
	assert.Equal(t, serviceerror.CategoryAlreadyExists, err)
}

func TestUpdateCategory_Success(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockCategoryRepository := mock.NewMockCategoryRepository(ctl)
	mockManager := &repository.Manager{
		Category: mockCategoryRepository,
	}

	categoryService := category.NewCategoryService(mockManager)

	categoryUpdateDTO := &model.CategoryUpdateDTO{
		Id:          1,
		UserId:      1,
		Name:        ptr[string]("Updated category name"),
		Description: ptr[string]("Updated description"),
	}

	existingCategory := &entity.Category{
		Id:          1,
		UserId:      1,
		Name:        "Old category name",
		Description: "Old description",
		Type:        "INCOME",
	}

	updatedCategory := &entity.Category{
		Id:          1,
		UserId:      1,
		Name:        "Updated category name",
		Description: "Updated description",
		Type:        "INCOME",
	}

	mockCategoryRepository.
		EXPECT().
		GetCategoryById(categoryUpdateDTO.Id).
		Times(1).
		Return(existingCategory, nil)

	mockCategoryRepository.
		EXPECT().
		ExistsWithName(existingCategory.UserId, *categoryUpdateDTO.Name).
		Times(1).
		Return(false)

	mockCategoryRepository.
		EXPECT().
		UpdateCategory(existingCategory).
		Times(1).
		DoAndReturn(func(category *entity.Category) (*entity.Category, error) {
			category.Name = *categoryUpdateDTO.Name
			category.Description = *categoryUpdateDTO.Description
			return category, nil
		})

	updatedCategoryFromService, err := categoryService.UpdateCategory(*categoryUpdateDTO)

	assert.NoError(t, err)
	assert.NotNil(t, updatedCategoryFromService)
	assert.Equal(t, updatedCategory, updatedCategoryFromService)
}

func TestUpdateCategory_CategoryNotFound_Error(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockCategoryRepository := mock.NewMockCategoryRepository(ctl)
	mockManager := &repository.Manager{
		Category: mockCategoryRepository,
	}

	categoryService := category.NewCategoryService(mockManager)

	categoryUpdateDTO := &model.CategoryUpdateDTO{
		Id:          1,
		Name:        ptr[string]("Non-existent category"),
		Description: ptr[string]("This category does not exist"),
	}

	mockCategoryRepository.
		EXPECT().
		GetCategoryById(categoryUpdateDTO.Id).
		Times(1).
		Return(nil, serviceerror.CategoryDoesntExist)

	updatedCategory, err := categoryService.UpdateCategory(*categoryUpdateDTO)

	assert.Error(t, err)
	assert.Nil(t, updatedCategory)
	assert.Equal(t, serviceerror.CategoryDoesntExist, err)
}

func TestDeleteCategory_Success(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockCategoryRepository := mock.NewMockCategoryRepository(ctl)
	mockManager := &repository.Manager{
		Category: mockCategoryRepository,
	}

	categoryService := category.NewCategoryService(mockManager)

	categoryDeleteDTO := &model.CategoryDeleteDTO{
		Id: 1,
	}

	mockCategoryRepository.
		EXPECT().
		CategoryBelongsToUser(categoryDeleteDTO.Id, categoryDeleteDTO.UserId).
		Times(1).
		Return(true)

	mockCategoryRepository.
		EXPECT().
		DeleteCategory(categoryDeleteDTO.Id).
		Times(1).
		Return(nil)

	err := categoryService.DeleteCategory(*categoryDeleteDTO)

	assert.NoError(t, err)
}

func TestDeleteCategory_CategoryDoesntBelongToUser_Error(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockCategoryRepository := mock.NewMockCategoryRepository(ctl)
	mockManager := &repository.Manager{
		Category: mockCategoryRepository,
	}

	categoryService := category.NewCategoryService(mockManager)

	categoryDeleteDTO := &model.CategoryDeleteDTO{
		Id:     1,
		UserId: 1,
	}

	mockCategoryRepository.
		EXPECT().
		CategoryBelongsToUser(categoryDeleteDTO.Id, categoryDeleteDTO.UserId).
		Times(1).
		Return(false)

	err := categoryService.DeleteCategory(*categoryDeleteDTO)

	assert.Error(t, err)
	assert.Equal(t, serviceerror.CategoryDoesntBelongToUser, err)
}
