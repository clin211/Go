package model

import (
	"blog-service/global"
	"blog-service/pkg/setting"
	"fmt"

	"github.com/jinzhu/gorm"
)

type Model struct {
	ID         uint32 `gorm:"primary_key" json:"id"` // 自增长id
	CreatedBy  string `json:"created_by"`            // 创建人
	ModifiedBy string `json:"modified_by"`           // 修改人
	CreatedOn  uint32 `json:"created_on"`            // 创建时间
	ModifiedOn uint32 `json:"modified_on"`           // 修改时间
	DeletedOn  uint32 `json:"deleted_on"`            // 删除时间
	IsDel      uint8  `json:"is_del"`
}

func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	db, err := gorm.Open(databaseSetting.DBType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime,
	))

	if err != nil {
		return nil, err
	}

	if global.ServerSetting.RunMode == "debug" {
		db.LogMode(true)
	}

	db.SingularTable(true)
	db.DB().SetMaxIdleConns(databaseSetting.MaxIdleConns)
	db.DB().SetMaxOpenConns(databaseSetting.MaxOpenConns)

	return db, nil
}
