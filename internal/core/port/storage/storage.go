package storage

import "github.com/khivuksergey/portmonetka.category/internal/core/port/repository"

type IDB interface {
	InitRepositoryManager() *repository.Manager
	Close() error
}
