package model

import (
	"time"

	"gorm.io/gorm"

	"github.com/Forest-211/miniblog/pkg/auth"
)

type UserM struct {
	ID        int64     `gorm:"column:id;primary_key"` //用户ID
	Username  string    `gorm:"column:username"`       //用户名
	Password  string    `gorm:"column:password"`       //密码
	Nickname  string    `gorm:"column:nickname"`       //昵称
	Email     string    `gorm:"column:email"`          //电子邮件地址
	Phone     string    `gorm:"column:phone"`          //手机号码
	CreatedAt time.Time `gorm:"column:createdAt"`      //创建时间
	UpdatedAt time.Time `gorm:"column:updatedAt"`      //更新时间
}

// TableName sets the insert table name for this struct type
func (u *UserM) TableName() string {
	return "user"
}

// BeforeCreate 在创建数据库记录之前加密明文密码.
func (u *UserM) BeforeCreate(tx *gorm.DB) (err error) {
	// Encrypt the user password.
	u.Password, err = auth.Encrypt(u.Password)
	if err != nil {
		return err
	}

	return nil
}
