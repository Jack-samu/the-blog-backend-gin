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

type LoginResp struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refreshToken"`
	UserInfo     *UserInfo `json:"userInfo"`
}

type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar,omitempty"`
	Posts    int64  `json:"posts"`
}

type RefreshResp struct {
	Token    string    `json:"token"`
	UserInfo *UserInfo `json:"userInfo"`
}
