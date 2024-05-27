package gorm

import (
	"fmt"
	"github.com/khivuksergey/portmonetka.category/config"
	"github.com/khivuksergey/portmonetka.category/internal/adapter/storage/entity"
	"github.com/khivuksergey/portmonetka.category/internal/adapter/storage/gorm/repo"
	"github.com/khivuksergey/portmonetka.category/internal/core/port/repository"
	"github.com/khivuksergey/portmonetka.category/internal/core/port/storage"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type dbManager struct {
	db  *gorm.DB
	cfg *config.DBConfig
}

func NewDbManager(config config.DBConfig) storage.IDB {
	dbm := dbManager{}
	err := dbm.InitDB(config)
	if err != nil {
		panic(err)
	}
	return &dbm
}

func (m *dbManager) InitDB(config config.DBConfig) (err error) {
	dsn := fmt.Sprintf(config.ConnectionString,
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_NAME"),
		viper.GetString("DB_HOST"),
	)

	m.db, err = gorm.Open(
		postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true,
		}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)

	if err != nil {
		return err
	}

	err = m.db.AutoMigrate(&entity.Category{})

	return err
}

func (m *dbManager) InitRepositoryManager() *repository.Manager {
	return &repository.Manager{
		Category: repo.NewCategoryRepository(m.db),
	}
}

func (m *dbManager) Close() (err error) {
	db, err := m.db.DB()
	if err != nil {
		return
	}
	err = db.Close()
	return
}
