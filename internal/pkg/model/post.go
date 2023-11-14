package model

import "time"

// PostM 是数据库中 post 记录 struct 格式的映射.
type PostM struct {
	ID        int64     `gorm:"column:id;primary_key" json:"id"`   //id
	Username  string    `gorm:"column:username" json:"username"`   //用户名
	PostID    string    `gorm:"column:postID" json:"postID"`       //帖子ID
	Title     string    `gorm:"column:title" json:"title"`         //标题
	Content   string    `gorm:"column:content" json:"content"`     //内容
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"` //创建时间
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"` //更新时间
}

// TableName 用来指定映射的 MySQL 表名.
func (p *PostM) TableName() string {
	return "post"
}
