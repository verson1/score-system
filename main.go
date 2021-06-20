package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"score-system/service"
)

const DbDns = "root:123456@tcp(101.132.140.237:3306)/test?charset=utf8mb4&parseTime=true"

func main() {

	db, err := gorm.Open(mysql.Open(DbDns), &gorm.Config{
		Logger: logger.Default,
	})
	if err != nil {
		panic(err)
	}

	svc := service.NewService(db)

	engine := gin.Default()

	engine.POST("/upload-file", svc.UpdateFile)
	engine.GET("/get-achievement", svc.GetAchievement)

	if err := engine.Run(":8080"); err != nil {
		panic(err)
	}

}
