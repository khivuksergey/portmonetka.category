package error

import (
	"errors"
	"fmt"
)

var (
	CategoryAlreadyExists          = errors.New("category with this name already exists")
	CategoryDoesntExist            = errors.New("category with this id doesn't exists")
	CategoryDoesntBelongToUser     = errors.New("category with this id doesn't belong to user")
	AtLeastOneFieldIsRequired      = errors.New("at least one field for updating category is required")
	CategoryNameLengthError        = errors.New("category name must be from 3 to 128 symbols long")
	CategoryDescriptionLengthError = errors.New("category description must be less than 256 symbols long")
)

const (
	InvalidInputData     = "invalid input data"
	CannotCreateCategory = "cannot create category"
	CannotGetCategories  = "cannot retrieve categories"
	CannotUpdateCategory = "cannot update category"
	CannotDeleteCategory = "cannot delete category"
)

type ErrorMessage string

func (m *ErrorMessage) Append(errMessage string) {
	if *m != "" {
		*m += "; "
	}
	*m += ErrorMessage(errMessage)
}

func (m *ErrorMessage) ToError() error {
	if *m == "" {
		return nil
	}
	return fmt.Errorf(fmt.Sprint(*m))
}
