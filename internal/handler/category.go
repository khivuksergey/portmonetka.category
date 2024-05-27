package handler

import (
	"github.com/go-playground/validator/v10"
	serviceerror "github.com/khivuksergey/portmonetka.category/error"
	"github.com/khivuksergey/portmonetka.category/internal/core/port/service"
	"github.com/khivuksergey/portmonetka.category/internal/model"
	"github.com/khivuksergey/portmonetka.common"
	"github.com/khivuksergey/webserver/logger"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type CategoryHandler struct {
	categoryService service.CategoryService
	logger          logger.Logger
	validate        *validator.Validate
}

func NewCategoryHandler(services *service.Manager, logger logger.Logger) *CategoryHandler {
	return &CategoryHandler{
		categoryService: services.Category,
		logger:          logger,
		validate:        model.GetCategoryValidator(),
	}
}

// GetCategories retrieves user's categories.
//
// @Tags Category
// @Summary Get user's categories
// @Description Gets user's categories
// @ID get-categories
// @Accept json
// @Produce json
// @Param userId path uint64 true "Authorized user ID"
// @Success 200 {object} model.Response "Categories retrieved"
// @Failure 422 {object} model.Response "Unprocessable entity"
// @Router /users/{userId}/categories [get]
func (w CategoryHandler) GetCategories(c echo.Context) error {
	requestUuid := c.Get(common.RequestUuidKey).(string)
	userId := c.Get("userId").(uint64)

	categories, err := w.categoryService.GetCategoriesByUserId(userId)
	if err != nil {
		return common.NewUnprocessableEntityError(serviceerror.CannotGetCategories, err)
	}

	w.logger.Info(logger.LogMessage{
		Action:      "GetCategories",
		Message:     "Categories retrieved",
		UserId:      &userId,
		RequestUuid: requestUuid,
	})

	return c.JSON(http.StatusOK, model.Response{
		Message:     "Categories retrieved",
		Data:        categories,
		RequestUuid: requestUuid,
	})
}

// CreateCategory creates a new category for user.
//
// @Tags Category
// @Summary Create a new category
// @Description Creates a new category with the provided information
// @ID create-category
// @Accept json
// @Produce json
// @Param userId path uint64 true "Authorized user ID"
// @Param category body model.CategoryCreateDTO true "Category object to be created"
// @Success 201 {object} model.Response "Category created"
// @Failure 400 {object} model.Response "Bad request"
// @Failure 422 {object} model.Response "Unprocessable entity"
// @Router /users/{userId}/categories [post]
func (w CategoryHandler) CreateCategory(c echo.Context) error {
	requestUuid := c.Get(common.RequestUuidKey).(string)
	userId := c.Get("userId").(uint64)
	categoryCreateDTO := &model.CategoryCreateDTO{}

	err := bindDtoValidate[model.CategoryCreateDTO](c, w.validate, categoryCreateDTO)
	if err != nil {
		return common.NewValidationError(serviceerror.InvalidInputData, err)
	}

	categoryCreateDTO.UserId = userId

	category, err := w.categoryService.CreateCategory(*categoryCreateDTO)
	if err != nil {
		return common.NewUnprocessableEntityError(serviceerror.CannotCreateCategory, err)
	}

	w.logger.Info(logger.LogMessage{
		Action:      "CreateCategory",
		Message:     "Category created",
		UserId:      &category.UserId,
		Data:        map[string]uint64{"id": category.Id},
		RequestUuid: requestUuid,
	})

	return c.JSON(http.StatusCreated, model.Response{
		Message:     "Category created",
		Data:        category,
		RequestUuid: requestUuid,
	})
}

// UpdateCategory updates the category.
//
// @Tags Category
// @Summary Update category
// @Description Updates category's properties
// @ID update-category
// @Accept json
// @Produce json
// @Param userId path uint64 true "Authorized user ID"
// @Param category body model.CategoryUpdateDTO true "Category update attributes"
// @Success 200 {object} model.Response "Category updated"
// @Failure 400 {object} model.Response "Bad request"
// @Failure 422 {object} model.Response "Unprocessable entity"
// @Router /users/{userId}/categories/{categoryId} [patch]
func (w CategoryHandler) UpdateCategory(c echo.Context) error {
	requestUuid := c.Get(common.RequestUuidKey).(string)
	userId := c.Get("userId").(uint64)
	categoryId, _ := strconv.ParseUint(c.Param("categoryId"), 10, 64)
	categoryUpdateDTO := &model.CategoryUpdateDTO{}

	err := bindDtoValidate[model.CategoryUpdateDTO](c, w.validate, categoryUpdateDTO)
	if err != nil {
		return common.NewValidationError(serviceerror.InvalidInputData, err)
	}

	categoryUpdateDTO.Id = categoryId
	categoryUpdateDTO.UserId = userId

	category, err := w.categoryService.UpdateCategory(*categoryUpdateDTO)
	if err != nil {
		return common.NewUnprocessableEntityError(serviceerror.CannotUpdateCategory, err)
	}

	w.logger.Info(logger.LogMessage{
		Action:      "UpdateCategory",
		Message:     "Category updated",
		UserId:      &userId,
		Data:        map[string]uint64{"id": category.Id},
		RequestUuid: requestUuid,
	})

	return c.JSON(http.StatusOK, model.Response{
		Message:     "Category updated",
		Data:        category,
		RequestUuid: requestUuid,
	})
}

// DeleteCategory deletes the category by ID.
//
// @Tags Category
// @Summary Delete category
// @Description Deletes category by the provided category ID
// @ID delete-category
// @Accept json
// @Produce json
// @Param userId path uint64 true "Authorized user ID"
// @Param category body model.CategoryDeleteDTO true "Category delete request"
// @Success 204 {string} string "No content"
// @Failure 400 {object} model.Response "Bad request"
// @Failure 422 {object} model.Response "Unprocessable entity"
// @Router /users/{userId}/categories/{categoryId} [delete]
func (w CategoryHandler) DeleteCategory(c echo.Context) error {
	requestUuid := c.Get(common.RequestUuidKey).(string)
	userId := c.Get("userId").(uint64)
	categoryId, _ := strconv.ParseUint(c.Param("categoryId"), 10, 64)
	categoryDeleteDTO := &model.CategoryDeleteDTO{}

	err := bindDtoValidate[model.CategoryDeleteDTO](c, w.validate, categoryDeleteDTO)
	if err != nil {
		return common.NewValidationError(serviceerror.InvalidInputData, err)
	}

	categoryDeleteDTO.Id = categoryId
	categoryDeleteDTO.UserId = userId

	if err := w.categoryService.DeleteCategory(*categoryDeleteDTO); err != nil {
		return common.NewUnprocessableEntityError(serviceerror.CannotDeleteCategory, err)
	}

	w.logger.Info(logger.LogMessage{
		Action:      "DeleteCategory",
		Message:     "Category deleted",
		UserId:      &userId,
		Data:        map[string]uint64{"id": categoryDeleteDTO.Id},
		RequestUuid: requestUuid,
	})

	return c.NoContent(http.StatusNoContent)
}

func bindDtoValidate[T any](c echo.Context, validate *validator.Validate, dto *T) error {
	if err := c.Bind(dto); err != nil {
		return err
	}
	if err := validate.Struct(dto); err != nil {
		return err
	}
	return nil
}
