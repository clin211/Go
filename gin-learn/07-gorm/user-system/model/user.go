package model

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint       `gorm:"primaryKey"`
	Username      string     `gorm:"size:20;uniqueIndex;not null"`
	Password      string     `gorm:"size:100;not null"`
	Email         string     `gorm:"size:50;uniqueIndex;not null"`
	Phone         string     `gorm:"size:20;uniqueIndex"`
	IsActive      bool       `gorm:"default:false"`
	LoginAttempts *time.Time `gorm:"default:null"`
	LastLoginAt   *time.Time `gorm:"default:null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

// 表名
func (User) TableName() string {
	return "users"
}

// 密码加密钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password != "" {
		u.Password = HashPassword(u.Password)
	}
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("Password") {
		u.Password = HashPassword(u.Password)
	}
	return nil
}

// 密码加密
func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

// 数据验证
func ValidateUser(user *User) error {
	// 用户名验证
	if user.Username == "" {
		return errors.New("用户名不能为空")
	}
	if len(user.Username) < 3 || len(user.Username) > 20 {
		return errors.New("用户名长度必须在3-20个字符之间")
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !usernameRegex.MatchString(user.Username) {
		return errors.New("用户名只能包含字母、数字和下划线")
	}

	// 密码验证
	if user.Password == "" {
		return errors.New("密码不能为空")
	}
	if len(user.Password) < 6 || len(user.Password) > 20 {
		return errors.New("密码长度必须在6-20个字符之间")
	}
	hasLetter := strings.ContainsAny(user.Password, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	hasNumber := strings.ContainsAny(user.Password, "0123456789")
	if !hasLetter || !hasNumber {
		return errors.New("密码必须同时包含字母和数字")
	}

	// 邮箱验证
	if user.Email == "" {
		return errors.New("邮箱不能为空")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(user.Email) {
		return errors.New("邮箱格式不正确")
	}

	// 手机号验证（可选）
	if user.Phone != "" {
		phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
		if !phoneRegex.MatchString(user.Phone) {
			return errors.New("手机号格式不正确")
		}
	}

	return nil
}
