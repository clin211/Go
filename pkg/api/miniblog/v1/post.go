package v1

type CreatePostRequest struct {
	Title   string `json:"title" valid:"required,stringlength(1|255)"`
	Content string `json:"content" valid:"required,stringlength(1|4294967295)"`
}

type ListPostRequest struct {
	Username string `json:"username" valid:"alphanum,required,stringlength(1|255)"`
}

type PostByIDRequest struct {
	ID string `json:"id" valid:"required"`
}
