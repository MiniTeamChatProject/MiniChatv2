package data

import (
	"fmt"
	"room/internal/conf"
	"room/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewRoomRepo, NewRoomMemberRepo)

// Data .
type Data struct {
	db *gorm.DB
}

// NewData .
func NewData(c *conf.Data, l log.Logger) (*Data, func(), error) {
	// Configure GORM
	db, err := gorm.Open(postgres.Open(c.Database.Source), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(int(c.Database.MaxOpenConns))
	sqlDB.SetMaxIdleConns(int(c.Database.MaxIdleConns))
	if c.Database.ConnMaxLifetime != nil {
		sqlDB.SetConnMaxLifetime(c.Database.ConnMaxLifetime.AsDuration())
	}

	// Auto migrate (development only)
	if c.Database.EnableAutoMigrate {
		helper := log.NewHelper(l)
		helper.Info("Enabling auto migration")
		if err := db.AutoMigrate(&model.Room{}, &model.RoomMember{}); err != nil {
			return nil, nil, fmt.Errorf("failed to auto migrate: %w", err)
		}
	}

	cleanup := func() {
		log.Info("closing the database connection")
		sqlDB.Close()
	}

	log.Info("database connection established")
	return &Data{db: db}, cleanup, nil
}

// DB returns the GORM database instance
func (d *Data) DB() *gorm.DB {
	return d.db
}
