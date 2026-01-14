//path: ./registration/internal/data/user.go 


package data

import (
	"registration/internal/conf"
	"os"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"fmt"
)

// ProviderSet 告诉 Wire：Data 层提供了这些东西（数据库连接、用户仓库实现）
var ProviderSet = wire.NewSet(NewData, NewDB, NewUserRepo)

type Data struct {
	db *gorm.DB
}

// 1. 初始化数据库连接 (NewDB)
// 这里我们为了演示方便，先把配置写死。
// 正规的 Kratos 是通过 c *conf.Data 读取配置文件的，以后可以改。
func NewDB(c *conf.Data) *gorm.DB {

	// 1. 尝试从环境变量获取数据库地址 (适配 Docker)
    dbHost := os.Getenv("DB_HOST")
    if dbHost == "" {
        dbHost = "localhost" // 默认回退到本地开发
    }
    
    // 2. 拼接 DSN
    // 注意把 host=%s 替换进去
    // 使用你之前验证过的连接字符串
     dsn := fmt.Sprintf("host=%s user=postgres password=07210721 dbname=benutzer_db port=5432 sslmode=disable TimeZone=Asia/Shanghai", dbHost)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// 自动迁移表结构 (这一步就像之前的 DB.AutoMigrate)
	// 注意：这里引用的是下面 user.go 里定义的 User 结构体，此时还没写，暂时先注释掉，或者等会儿写完再回来开
	 if err := db.AutoMigrate(&User{}); err != nil {
	  	panic(err)
	 }
	
	return db
}

// 2. 初始化 Data 结构体
func NewData(c *conf.Data, logger log.Logger, db *gorm.DB) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		// 如果需要关闭数据库连接池，可以在这里写
	}
	return &Data{db: db}, cleanup, nil
}
