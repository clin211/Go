package main

import (
	"fmt"

	"github.com/clin211/08-gorm-gen/dal/query"
	"github.com/clin211/08-gorm-gen/internal/api"
	"github.com/clin211/08-gorm-gen/internal/config"
	"github.com/clin211/08-gorm-gen/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	if err := config.InitConfig(); err != nil {
		fmt.Println("failed to init config:", err)
		panic(err)
	}

	_, err := initDB()
	if err != nil {
		fmt.Println("failed to connect database:", err)
		panic(err)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	userService := service.NewUserService()
	userAPI := api.NewUserAPI(userService)

	// 定义 API 路由
	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.GET("", userAPI.GetUsers)
			users.GET("/:id", userAPI.GetUserByID)
			users.POST("", userAPI.CreateUser)
			users.PUT("/:id", userAPI.UpdateUser)
			users.DELETE("/:id", userAPI.DeleteUser)
		}
	}

	r.Run(config.App.Port)
}

func initDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.MySQL.User,
		config.MySQL.Password,
		config.MySQL.Host,
		config.MySQL.Port,
		config.MySQL.DBName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	// 设置 Gorm Gen 使用的默认数据库
	query.SetDefault(db)

	return db, err
}
