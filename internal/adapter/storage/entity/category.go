package entity

import (
	"gorm.io/gorm"
	"time"
)

type Category struct {
	Id          uint64         `json:"id" gorm:"primarykey"`
	UserId      uint64         `json:"userId" gorm:"not null;uniqueIndex:idx_userid_name_deletedat" validate:"required"`
	Name        string         `json:"name" gorm:"not null;uniqueIndex:idx_userid_name_deletedat" validate:"required,min=3,max=128"`
	Description string         `json:"description" gorm:"null" validate:"max=256"`
	Type        CategoryType   `json:"type" gorm:"not null" validate:"required,oneof=INCOME EXPENSE"`
	CreatedAt   time.Time      `json:"createdAt" gorm:"<-:create"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index;uniqueIndex:idx_userid_name_deletedat"`
}

func (Category) TableName() string { return "portmonetka.categories" }

type CategoryType string

const (
	Income  CategoryType = "INCOME"
	Expense CategoryType = "EXPENSE"
)
