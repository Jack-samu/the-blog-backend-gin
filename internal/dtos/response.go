package dtos

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

type ProfileResp struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Articles int64  `json:"articleCount"`
	Drafts   int64  `json:"draftCount"`
	Avatar   string `json:"avatar"`
}

type PhotoItem struct {
	ID  uint   `json:"id"`
	Img string `json:"name"`
}

type PhotosResp struct {
	Photos []PhotoItem `json:"photos"`
}

type ArticleBasic struct {
	Id        uint   `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type PostListItem struct {
	ArticleBasic
	Excerpt  string   `json:"excerpt"`
	Cover    string   `json:"cover"`
	Author   string   `json:"author"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	// 计数位
	Views    uint `json:"views"`
	Likes    uint `json:"likes"`
	Comments int  `json:"comments"`
}

type PostDetailItem struct {
	PostListItem
	Author  AuthorProfile `json:"author"`
	Content string        `json:"content"`
}

// 对于响应的构造可能需要一定的json定义
type CommentItem struct {
	//
}

type PostListResp struct {
	Posts       []PostListItem `json:"articles"`
	Cnt         uint           `json:"total"`
	CurrentPage uint           `json:"current_page"`
}

type PostDetailResp struct {
	Post PostDetailItem `json:"article"`
}

type DraftDetailResp struct {
	Draft DraftDetail `json:"article"`
}

type DraftsResp struct {
	Drafts      []DraftItem `json:"drafts"`
	Cnt         uint        `json:"total"`
	CurrentPage uint        `json:"current_page"`
}

type PostListPersonalResp struct {
	Posts       interface{} `json:"publishedArticles"`
	Cnt         uint        `json:"total"`
	CurrentPage uint        `json:"current_page"`
}

type SeriesResp struct {
	Categories interface{} `json:"categories"`
}

type SeriesItem struct {
	ID    uint           `json:"id"`
	Name  string         `json:"name"`
	Posts []ArticleBasic `json:"articles"`
}

type DraftDetail struct {
	ArticleBasic
	Excerpt  string        `json:"excerpt"`
	Cover    string        `json:"cover"`
	Category string        `json:"category"`
	Tags     []string      `json:"tags"`
	Author   AuthorProfile `json:"author"`
	Content  string        `json:"content"`
}

type DraftItem struct {
	ArticleBasic
	Excerpt  string   `json:"excerpt"`
	Cover    string   `json:"cover"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	Author   string   `json:"author"`
}

type AuthorProfile struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

func NewPostList(list []PostListItem, total, page int64) *PostListResp {
	return &PostListResp{
		Posts:       list,
		Cnt:         uint(total),
		CurrentPage: uint(page),
	}
}

func NewPostListPersonal(list []PostListItem, total, page int64) *PostListPersonalResp {
	return &PostListPersonalResp{
		Posts:       list,
		Cnt:         uint(total),
		CurrentPage: uint(page),
	}
}

func NewDraftsPersonal(list []DraftItem, total, page int64) *DraftsResp {
	return &DraftsResp{
		Drafts:      list,
		Cnt:         uint(total),
		CurrentPage: uint(page),
	}
}
