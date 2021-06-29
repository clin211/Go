package model

import "blog-service/pkg/app"

type Tag struct {
	*Model
	Name  string `json:"name"`  // 标签名称
	State uint8  `json:"state"` // 状态: 0为禁用、1为启用
}

func (a Tag) TableName() string {
	return "blog_tag"
}

type TagSwagger struct {
	List []*Tag
	Pager *app.Pager
}