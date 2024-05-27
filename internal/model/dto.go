package model

import (
	"github.com/khivuksergey/portmonetka.category/internal/adapter/storage/entity"
)

type CategoryCreateDTO struct {
	UserId      uint64              `json:"userId"`
	Name        string              `json:"name" validate:"required"`
	Description string              `json:"description"`
	Type        entity.CategoryType `json:"type" validate:"required,oneof=INCOME EXPENSE"`
}

type CategoryUpdateDTO struct {
	Id          uint64  `json:"id"`
	UserId      uint64  `json:"userId"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type CategoryDeleteDTO struct {
	Id     uint64 `json:"id"`
	UserId uint64 `json:"userId"`
}
