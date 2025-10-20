package service

import (
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Jack-samu/the-blog-backend-gin.git/internal/dtos"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/errs"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/models"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 用户查询和创建的调度处理
func (s *Service) Register(username, email, password, bio, avatar string) *errs.ErrorResp {
	exists, err := s.r.ExistByEmail(email)
	if err != nil {
		log.Printf("用户查询出错：%v\n", err)
		return errs.NewError(http.StatusInternalServerError, "用户查询出错", err)
	}

	if exists {
		return errs.NewError(http.StatusBadRequest, "用户邮箱已注册", nil)
	}

	exists, err = s.r.ExistByUsername(username)
	if err != nil {
		log.Printf("用户查询出错：%v\n", err)
		return errs.NewError(http.StatusInternalServerError, "用户查询出错", err)
	}

	if exists {
		return errs.NewError(http.StatusBadRequest, "昵称重复", nil)
	}

	user := &models.User{
		ID:            uuid.New().String(),
		Username:      username,
		Email:         email,
		Bio:           bio,
		FailedLogin:   0,
		CaptchaReqCnt: 0,
		CreatedAt:     time.Now(),
	}

	err = user.SetPassword(password)
	if err != nil {
		log.Printf("哈希加密出错：%v\n", err.Error())
		return errs.NewError(http.StatusInternalServerError, "用户注册中哈希加密出错", err)
	}

	if err = s.r.CreateUser(user); err != nil {
		log.Printf("user创建出错：%v\n", err.Error())
		return errs.NewError(http.StatusInternalServerError, "用户注册失败", err)
	}

	if avatar != "" {
		a := &models.Img{
			Name:      avatar,
			IsAvatar:  true,
			CreatedAt: time.Now(),
		}

		// 添加头像的注册
		if err = s.r.AddAvatar(user, a); err != nil {
			log.Printf("头像添加出错：%v\n", err.Error())
			return errs.NewError(http.StatusInternalServerError, "用户注册失败", err)
		}
	}

	return nil

}

func (s *Service) Login(username, password string) (*dtos.LoginResp, *errs.ErrorResp) {

	// 用户查询
	user, avatar, err := s.r.GetUserByNameWithAvatar(username)
	if err != nil {
		log.Printf("用户查询出错：%v\n", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewError(http.StatusBadRequest, "用户不存在", nil)
		}
		return nil, errs.NewError(http.StatusInternalServerError, "用户查询出错", err)
	}

	// 登录锁定检查
	if user.FailedLogin >= 5 {
		return nil, errs.NewError(http.StatusBadRequest, "该用户登录失败次数过多，已锁定，请明天再试", nil)
	}

	if !user.CheckPassword(password) {
		if err = s.r.IncreaseFailedLogin(user); err != nil {
			log.Printf("用户登陆失败计数位重设失败：%v\n", err)
		}
		return nil, errs.NewError(http.StatusBadRequest, "密码错误", nil)
	}

	posts, err := s.r.GetUserPosts(user.ID)
	if err != nil {
		log.Printf("获取用户文章数量失败：%v\n", err)
		posts = 0
	}

	// token生成
	token, err := utils.GenerateToken(user.ID, time.Hour)
	if err != nil {
		return nil, errs.NewError(http.StatusInternalServerError, "token生成失败，请重试", nil)
	}

	// refreshtoken生成
	refreshToken, err := utils.GenerateToken(user.ID, 24*time.Hour)
	if err != nil {
		return nil, errs.NewError(http.StatusInternalServerError, "refreshToken生成失败，请重试", nil)
	}

	return &dtos.LoginResp{
		Token:        token,
		RefreshToken: refreshToken,
		UserInfo: &dtos.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Avatar:   avatar,
			Posts:    posts,
		},
	}, nil
}

func (s *Service) RefreshTheToken(id string) (*dtos.RefreshResp, *errs.ErrorResp) {
	// 查询用户是否存在
	user, avatar, err := s.r.GetUserByNameWithAvatar(id)
	if err != nil {
		log.Printf("用户查询出错：%v\n", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewError(http.StatusBadRequest, "用户不存在", nil)
		}
		return nil, errs.NewError(http.StatusInternalServerError, "用户查询出错", err)
	}

	posts, err := s.r.GetUserPosts(user.ID)
	if err != nil {
		log.Printf("获取用户文章数量失败：%v\n", err)
		posts = 0
	}

	token, err := utils.GenerateToken(id, time.Hour)
	if err != nil {
		return nil, errs.NewError(http.StatusInternalServerError, "token生成失败，请重试", nil)
	}

	return &dtos.RefreshResp{
		Token: token,
		UserInfo: &dtos.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Avatar:   avatar,
			Posts:    posts,
		},
	}, nil
}

func (s *Service) SendCaptcha(username string) *errs.ErrorResp {
	user, err := s.r.GetUserByName(username)
	if err != nil {
		log.Printf("用户查询出错：%v\n", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NewError(http.StatusBadRequest, "用户不存在", nil)
		}
		return errs.NewError(http.StatusInternalServerError, "", err)
	}

	// 检查是否发送过
	val, ok := s.u.Load(user.ID)
	if ok {
		captcha, ok := val.(utils.Captcha)
		if !ok {
			return errs.NewError(http.StatusInternalServerError, "加载存储的验证码报错", err)
		}

		// 检查是否还没过请求冷却
		if captcha.Since.Before(time.Now()) {
			return errs.NewError(http.StatusBadRequest, "让我歇会", nil)
		}
	}

	captcha, err := utils.GenerateCaptcha()
	if err != nil {
		return errs.NewError(http.StatusInternalServerError, "token生成出错", err)
	}

	log.Println("生成的6位验证码：", captcha.Code)
	s.u.Store(user.ID, captcha)

	// 验证码邮件发送
	err = utils.SendCaptcha(user.Email, captcha.Code)
	if err != nil {
		return errs.NewError(http.StatusInternalServerError, "验证码发送失败", err)
	}

	// 最后对于验证码的请求进行计数位修改
	err = s.r.IncreaseCaptchaCnt(user)
	if err != nil {
		return errs.NewError(http.StatusInternalServerError, "重置请求码计数位出错", err)
	}

	return nil
}

func (s *Service) Verify(username, code string) *errs.ErrorResp {
	user, err := s.r.GetUserByName(username)
	if err != nil {
		log.Printf("用户查询出错：%v\n", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NewError(http.StatusBadRequest, "用户不存在", nil)
		}
		return errs.NewError(http.StatusInternalServerError, "", err)
	}

	val, ok := s.u.Load(user.ID)
	if !ok {
		return errs.NewError(http.StatusInternalServerError, "加载存储的验证码报错", err)
	}

	captcha, ok := val.(utils.Captcha)
	if !ok {
		return errs.NewError(http.StatusInternalServerError, "加载存储的验证码报错", err)
	}

	if captcha.Since.Add(10 * time.Minute).Before(time.Now()) {
		return errs.NewError(http.StatusBadRequest, "验证码已过期", nil)
	}

	if code != captcha.Code {
		return errs.NewError(http.StatusBadRequest, "验证码错误", nil)
	}

	token, err := utils.GenerateToken(user.ID, 20*time.Minute)
	if err != nil {
		return errs.NewError(http.StatusInternalServerError, "token生成失败，请重试", nil)
	}

	err = utils.SendResetLink(user.Email, token)
	if err != nil {
		return errs.NewError(http.StatusInternalServerError, "密码重置邮件发送失败", nil)
	}

	return nil
}

func (s *Service) Reset(userID, pwd string) *errs.ErrorResp {
	user, err := s.r.GetUserById(userID)
	if err != nil {
		log.Printf("用户查询出错：%v\n", err)
		return errs.NewError(http.StatusInternalServerError, "用户查询出错，请重试", nil)
	}

	if len(pwd) < 8 {
		return errs.NewError(http.StatusBadRequest, "密码长度不该小于8", nil)
	}

	if user.CheckPassword(pwd) {
		return errs.NewError(http.StatusBadRequest, "好啊，原密码", nil)
	}

	err = user.SetPassword(pwd)

	if err != nil {
		return errs.NewError(http.StatusInternalServerError, "", err)
	}

	return nil
}

func (s *Service) Logout(id string) (string, *errs.ErrorResp) {
	user, err := s.r.GetUserById(id)
	if err != nil {
		log.Printf("用户查询出错：%v\n", err)
		return "", errs.NewError(http.StatusInternalServerError, "用户查询出错，请重试", nil)
	}

	last_activity, err := s.r.SetLastActivity(user)
	if err != nil {
		return "", errs.NewError(http.StatusInternalServerError, "", nil)
	}
	return last_activity, nil
}

func (s *Service) UploadImg(file *multipart.FileHeader, imageType, uploadPath string) (string, *errs.ErrorResp) {
	mimetype := file.Header.Get("Content-type")
	ext := filepath.Ext(file.Filename)
	log.Println("debug，", mimetype, ext)
	if !s.TypeAllow(mimetype) {
		return "", errs.NewError(400, "仅支持PNG/JPG/JPEG/GIF格式", nil)
	}

	// 大小控制，后续再来添加类型方面的大小控制吧
	if file.Size > 10*1024*1024 {
		return "", errs.NewError(400, "超过最大可承受范围啦，吊毛", nil)
	}

	if ext == "" {
		switch mimetype {
		case "image/png":
			ext = ".png"
		case "image/jpg":
			ext = ".jpg"
		case "image/jpeg":
			ext = ".jpeg"
		case "image/gif":
			ext = ".gif"
		default:
			return "", errs.NewError(400, "不支持的图片格式", nil)
		}
	}

	filename := uuid.NewString() + strings.ToLower(ext)
	dst := filepath.Join(uploadPath, filename)

	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", errs.NewError(500, "文件保存失败，图片文件夹不存在", err)
	}

	if err := s.SaveImg(file, dst); err != nil {
		return "", errs.NewError(500, "文件保存失败", err)
	}

	return filename, nil
}

func (s *Service) TypeAllow(mimetype string) bool {
	lower_mime := strings.ToLower(mimetype)
	types := map[string]bool{
		"image/png":  true,
		"image/jpg":  true,
		"image/jpeg": true,
		"image/gif":  true,
	}
	return types[lower_mime]
}

// 只进行图片的本地文件保存动作
func (s *Service) SaveImg(f *multipart.FileHeader, filename string) error {
	src, err := f.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
