package v1

type CreatePostRequest struct {
	Title   string `json:"title" valid:"required,stringlength(1|255)"`
	Content string `json:"content" valid:"required,stringlength(1|4294967295)"`
}
