package data

import (
	"admin/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData)

// Data .
type Data struct {
	db *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}

	db, err := NewDB(c)
	if err != nil {
		// 数据库加载失败，拒绝启动项目（若无强制需要数据源可抛异常处理）
		panic(err)
	}

	return &Data{db: db}, cleanup, nil
}

func NewDB(c *conf.Data) (*gorm.DB, error) {
	source := c.Database.GetSource()

	db := &gorm.DB{}

	db, err := gorm.Open(mysql.Open(source), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
