package model

import "blog-service/pkg/app"

type Article struct {
	*Model
	Title         string `json:"name"`            // 文章标题
	Desc          string `json:"desc"`            // 文章简述
	Content       string `json:"content"`         // 文章内容
	CoverImageUrl string `json:"cover_image_url"` // 封面图地址
	State         uint8  `json:"state"`           // 状态：0为禁用、1为启用
}

func (a Article) TableName() string {
	return "blog_article"
}

type ArticleSwagger struct {
	List  []*Article
	Pager *app.Pager
}