package model

import "time"

// User 表示系统中的用户
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey;column:id"`
	Username  string    `json:"username" gorm:"column:username;type:varchar(50);not null;uniqueIndex:idx_username"`
	Email     string    `json:"email" gorm:"column:email;type:varchar(100);not null;uniqueIndex:idx_email"`
	Password  string    `json:"-" gorm:"column:password;type:varchar(255);not null"` // 永不暴露密码
	FirstName string    `json:"first_name" gorm:"column:first_name;type:varchar(50);not null"`
	LastName  string    `json:"last_name" gorm:"column:last_name;type:varchar(50);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}

// UserCreate 表示创建新用户所需的数据
type UserCreate struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// UserUpdate 表示可以更新的用户数据
type UserUpdate struct {
	Email     string `json:"email" binding:"omitempty,email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UserResponse 表示返回给客户端的数据
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
