package model

import (
	"blog-service/pkg/app"

	"github.com/jinzhu/gorm"
)

type Tag struct {
	*Model
	Name  string `json:"name"`  // 标签名称
	State uint8  `json:"state"` // 状态: 0为禁用、1为启用
}

func (a Tag) TableName() string {
	return "tag"
}

type TagSwagger struct {
	List  []*Tag
	Pager *app.Pager
}

func (t Tag) Count(db *gorm.DB) (int, error) {
	var count int
	if t.Name != "" {
		db = db.Where("state = ?", t.Name)
	}

	db = db.Where("state = ?", t.State)

	if err := db.Model(&t).Where("is_del = ?", 0).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (t Tag) List(db *gorm.DB, pageOffset, pageSize int) ([]*Tag, error) {
	var tags []*Tag
	var err error
	if pageOffset >= 0 && pageSize > 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}

	if t.Name != "" {
		db = db.Where("name = ?", t.Name)
	}
	db = db.Where("state = ?", t.State)

	// 查找符合筛选条件的数据
	err = db.Where("is_del = ?", 0).Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (t Tag) Get(db *gorm.DB) (Tag, error) {
	var tag Tag
	err := db.Where("id = ? AND is_del = ? AND state = ?", t.ID, 0, t.State).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return tag, err
	}

	return tag, nil
}

// 创建数据
func (t Tag) Create(db *gorm.DB) error {
	return db.Create(&t).Error
}

// 更新传入的字段
func (t Tag) Update(db *gorm.DB, values interface{}) error {
	return db.Model(&t).Where("id = ? AND is_del = ?", t.ID, 0).Updates(values).Error
}

// 删除数据
func (t Tag) Delete(db *gorm.DB) error {
	return db.Where("id = ? AND is_del = ?", t.Model.ID, 0).Delete(&t).Error
}
