//path: ./registration/internal/data/user.go


package data

import (
	"registration/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

// ProviderSet 告诉 Wire：Data 层提供了这些东西（数据库连接、用户仓库实现）
var ProviderSet = wire.NewSet(NewData, NewDB, NewUserRepo)

type Data struct {
	db *gorm.DB
}

// 1. 初始化数据库连接 (NewDB)
// 优先使用环境变量，适配 Docker 环境
func NewDB(c *conf.Data) *gorm.DB {
	// 优先从环境变量获取 (适配 Docker)
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost" // 默认回退到本地开发
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" {
		dbPass = "07210721"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "benutzer_db"
	}
	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Error: failed to connect to the data-base: " + err.Error())
	}

	// 自动迁移表结构
	 if err := db.AutoMigrate(&User{}); err != nil {
	  	panic(err)
	 }

	return db
}

// 2. 初始化 Data 结构体
func NewData(c *conf.Data, logger log.Logger, db *gorm.DB) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("Closing the connexion to data resources")
		// 如果需要关闭数据库连接池，可以在这里写
	}
	return &Data{db: db}, cleanup, nil
}
