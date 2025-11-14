package dtos

type RegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Bio      string `json:"bio" binding:"required"`
	Avatar   string `json:"avatar"`
}

type LoginReq struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=8"`
}

type GetCodeReq struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
}

type VerifyReq struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Code     string `json:"verificationCode" binding:"required"`
}

type ArticleReq struct {
	// 判别项
	Id uint `json:"id,omitempty"`
	// publish的基础表单数据
	Title   string `json:"title"`
	Content string `json:"content"`
	Excerpt string `json:"excerpt"`
	Cover   string `json:"cover"`

	// 附加型表单
	Category string   `json:"category,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}
