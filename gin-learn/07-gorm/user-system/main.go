package main

import (
	"fmt"
	"gorm-demo/user-system/model"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// 连接数据库
	dsn := "gorm:gorm123456@tcp(127.0.0.1:3306)/gorm_demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 设置日志级别
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Successfully connected to database!")
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 设置空闲连接池中的最大连接数
	sqlDB.SetMaxOpenConns(100)          // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置连接可复用的最大时间

	// 自动迁移
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	err = createUser(db)
	if err != nil {
		log.Fatal("Failed to create user:", err)
	}

	err = batchCreateUsers(db)
	if err != nil {
		log.Fatal("Failed to batch create users:", err)
	}

	user, err := GetUserByID(db, 13)
	if err != nil {
		log.Fatal("Failed to get user:", err)
	}
	fmt.Printf("User: %+v\n", user)

	user, err = GetUserByUsername(db, "john_doe")
	if err != nil {
		log.Fatal("Failed to get user:", err)
	}
	fmt.Printf("User: %+v\n", user)

	user, err = GetUserByEmail(db, "david-miller@example.com")
	if err != nil {
		log.Fatal("Failed to get user:", err)
	}
	fmt.Printf("User: %+v\n", user)

	// 使用示例
	users, total, err := ListUsers(db, 1, 5, map[string]interface{}{
		"is_active": true,
	})
	if err != nil {
		log.Fatal("Failed to list users:", err)
	}
	fmt.Printf("Users: %+v\n", users)
	fmt.Printf("Total: %d\n", total)

	err = UpdateUserStatus(db, 1, true)
	if err != nil {
		log.Fatal("Failed to update user status:", err)
	}
	fmt.Println("User status updated successfully")

	err = UpdateUser(db, 1, map[string]interface{}{
		"email": "new_email@example.com",
		"phone": "13900139011",
	})
	if err != nil {
		log.Fatal("Failed to update user:", err)
	}
	fmt.Println("User updated successfully")

	err = SoftDeleteUser(db, 1)
	if err != nil {
		log.Fatal("Failed to soft delete user:", err)
	}
	fmt.Println("User soft deleted successfully")

	err = HardDeleteUser(db, 12)
	if err != nil {
		log.Fatal("Failed to hard delete user:", err)
	}
	fmt.Println("User hard deleted successfully")
}

func createUser(db *gorm.DB) error {
	// 创建用户
	user := &model.User{
		Username: "john_doe",
		Password: "secure123",
		Email:    "john@example.com",
		Phone:    "13800138000",
		IsActive: false,
	}

	// 数据验证
	if err := model.ValidateUser(user); err != nil {
		log.Fatal("Validation failed:", err)
		return err
	}

	// 创建用户
	result := db.Create(user)
	if result.Error != nil {
		log.Fatal("Failed to create user:", result.Error)
	}

	fmt.Printf("Successfully created user with ID: %d\n", user.ID)
	return nil
}

// 批量创建用户
func batchCreateUsers(db *gorm.DB) error {
	users := []*model.User{
		{
			Username: "john_done",
			Password: "secure123",
			Email:    "john-doe@example.com",
			Phone:    "13810138100",
			IsActive: false,
		},
		{
			Username: "alice_smith",
			Password: "alice456",
			Email:    "alice-smith@example.com",
			Phone:    "13900139000",
			IsActive: true,
		},
		{
			Username: "bob_wilson",
			Password: "bob789",
			Email:    "bob-wilson@example.com",
			Phone:    "13700137000",
			IsActive: true,
		},
		{
			Username: "emma_davis",
			Password: "emma123",
			Email:    "emma-davis@example.com",
			Phone:    "13600136000",
			IsActive: false,
		},
		{
			Username: "michael_brown",
			Password: "mike456",
			Email:    "michael-brown@example.com",
			Phone:    "13500135000",
			IsActive: true,
		},
		{
			Username: "sarah_jones",
			Password: "sarah789",
			Email:    "sarah-jones@example.com",
			Phone:    "13400134000",
			IsActive: false,
		},
		{
			Username: "david_miller",
			Password: "dave123",
			Email:    "david-miller@example.com",
			Phone:    "13300133000",
			IsActive: true,
		},
		{
			Username: "lisa_wang",
			Password: "lisa456",
			Email:    "lisa-wang@example.com",
			Phone:    "13200132000",
			IsActive: true,
		},
		{
			Username: "james_zhang",
			Password: "james789",
			Email:    "james-zhang@example.com",
			Phone:    "13100131000",
			IsActive: false,
		},
		{
			Username: "sophie_li",
			Password: "sophie123",
			Email:    "sophie-li@example.com",
			Phone:    "13000130000",
			IsActive: true,
		},
	}

	// 数据验证
	for _, user := range users {
		if err := model.ValidateUser(user); err != nil {
			return err
		}
	}

	// 批量创建
	result := db.CreateInBatches(users, 100) // 每批100条
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// 根据 ID 查询用户
func GetUserByID(db *gorm.DB, id uint) (*model.User, error) {
	var user model.User
	result := db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// 根据用户名查询
func GetUserByUsername(db *gorm.DB, username string) (*model.User, error) {
	var user model.User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// 根据邮箱查询
func GetUserByEmail(db *gorm.DB, email string) (*model.User, error) {
	var user model.User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// 分页查询用户列表
func ListUsers(db *gorm.DB, page, pageSize int, conditions map[string]interface{}) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	// 构建查询
	query := db.Model(&model.User{})

	// 添加查询条件
	for key, value := range conditions {
		query = query.Where(key+" = ?", value)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	result := query.Offset(offset).Limit(pageSize).Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return users, total, nil
}

// 更新用户状态
func UpdateUserStatus(db *gorm.DB, id uint, isActive bool) error {
	result := db.Model(&model.User{}).Where("id = ?", id).Update("is_active", isActive)
	return result.Error
}

// 更新用户信息
func UpdateUser(db *gorm.DB, id uint, updates map[string]interface{}) error {
	result := db.Model(&model.User{}).Where("id = ?", id).Updates(updates)
	return result.Error
}

// 软删除用户
func SoftDeleteUser(db *gorm.DB, id uint) error {
	result := db.Delete(&model.User{}, id)
	return result.Error
}

// 硬删除用户
func HardDeleteUser(db *gorm.DB, id uint) error {
	result := db.Unscoped().Delete(&model.User{}, id)
	return result.Error
}
