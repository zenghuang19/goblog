package model

import (
	"goblog/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// DB gorm.DB 对象
var DB *gorm.DB

func ConnectDB()*gorm.DB  {
	var err error
	
	config := mysql.New(mysql.Config{
		DSN: "root:root@tcp(mysql:3306)/cesi?charset=utf8&parseTime=True&loc=Local",
	})

	// 准备数据库连接池
	DB,err =gorm.Open(config,&gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})

	logger.LogError(err)

	return DB
}