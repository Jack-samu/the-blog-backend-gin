package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Jack-samu/the-blog-backend-gin.git/internal/dtos"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Register(c *gin.Context) {
	var req dtos.RegisterReq

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("%s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "请求参数无效" + err.Error(),
		})
		return
	}

	err := h.s.Register(req.Username, req.Email, req.Password, req.Bio, req.Avatar)
	if err != nil {
		switch err.Code {
		case 400:
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Msg,
			})
		case 500:
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err.Msg + err.Err.Error(),
			})
		}
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"msg": "注册成功",
		})
	}
}

func (h *Handler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "image part not found.",
		})
		return
	}

	imageType := c.PostForm("type")
	log.Printf("debug，%v\n", imageType)
	if imageType != "" {
		if !h.s.TypeAllow(imageType) {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": "仅支持PNG/JPG/JPEG/GIF格式",
			})
			return
		}
	}

	filename, errs := h.s.UploadImg(file, imageType, "images")
	if errs != nil {
		switch errs.Code {
		case 400:
			c.JSON(http.StatusBadRequest, gin.H{
				"err": errs.Msg,
			})
		case 500:
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": errs.Msg,
			})
		}
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"filename": filename,
		})
	}
}

func (h *Handler) Login(c *gin.Context) {

	var req dtos.LoginReq
	err := c.ShouldBindJSON(&req)

	if err != nil {
		log.Printf("%s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "无效参数",
		})
		return
	}

	loginResp, errs := h.s.Login(req.Username, req.Password)
	if errs != nil {
		switch errs.Code {
		case 400:
			c.JSON(http.StatusBadRequest, gin.H{
				"err": errs.Msg,
			})
		case 500:
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": errs.Msg + errs.Err.Error(),
			})
		}
	} else {
		c.JSON(http.StatusOK, loginResp)
	}
}

func (h *Handler) Logout(c *gin.Context) {
	user_id := c.GetString("user_id")
	if user_id == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "用户状态信息查询出错，请重试"})
	}

	last_activity, err := h.s.Logout(user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "退出出错" + err.Err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"last_activity": last_activity})
	}
}

func (h *Handler) RefreshTheToken(c *gin.Context) {
	user_id := c.GetString("user_id")
	if user_id == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "用户状态信息查询出错，请重试"})
	}

	refreshResp, errs := h.s.RefreshTheToken(user_id)
	if errs != nil {
		switch errs.Code {
		case 400:
			c.JSON(http.StatusBadRequest, gin.H{
				"err": errs.Msg,
			})
		case 500:
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": errs.Msg,
			})
		}
	} else {
		c.JSON(http.StatusOK, refreshResp)
	}
}

func (h *Handler) GetCaptcha(c *gin.Context) {

	var req dtos.GetCodeReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "无效参数",
		})
		return
	}

	errs := h.s.SendCaptcha(req.Username)
	if errs != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": ""})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "验证码已发送"})
}

func (h *Handler) VerifyCaptcha(c *gin.Context) {
	var req dtos.VerifyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "无效参数",
		})
		return
	}

	err := h.s.Verify(req.Username, req.Code)
	if err != nil {
		switch err.Code {
		case 400:
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Msg,
			})
		case 500:
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err.Msg,
			})
		}
		log.Println("啥？我漏了？")
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "校验通过，密码重设链接已发送到邮箱"})
}

func (h *Handler) ResetPassword(c *gin.Context) {
	token := c.Param("token")

	payload, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "token无效或者已过期"})
		return
	}

	switch c.Request.Method {
	case http.MethodGet:
		c.HTML(http.StatusOK, `            
				<form method="post">
					<input type="password" name="password" placeholder="输入新密码" required />
					<input type="password" name="pwdRepeat" placeholder="确认新密码" required />
					<input type="submit" value="重置密码">                   
				</form>`, nil)
	case http.MethodPost:
		password := c.PostForm("password")
		pwdConfirm := c.PostForm("pwdRepeat")

		if password != pwdConfirm {
			c.JSON(http.StatusBadRequest, gin.H{"err": "密码不一致"})
			return
		}

		err := h.s.Reset(payload.ID, password)
		if err != nil {
			switch err.Code {
			case http.StatusBadRequest:
				c.JSON(http.StatusBadRequest, gin.H{"err": err.Msg})
			case http.StatusInternalServerError:
				c.JSON(http.StatusInternalServerError, gin.H{"err": err.Err.Error()})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{"msg": "密码已重置"})
		}
	}
}

func (h *Handler) Profile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		log.Printf("获取用户信息的请求中id为'%s'\n", id)
		c.JSON(http.StatusBadRequest, gin.H{"err": "缺少id参数，检查路由"})
		return
	}

	profileResp, err := h.s.Profile(id)
	if err != nil {
		switch err.Code {
		case http.StatusBadRequest:
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Msg})
		case http.StatusInternalServerError:
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Err.Error()})
		}
	} else {
		c.JSON(http.StatusOK, profileResp)
	}
}

func (h *Handler) GetPhotos(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		log.Printf("获取用户信息的请求中id为'%s'\n", id)
		c.JSON(http.StatusBadRequest, gin.H{"err": "缺少id参数，检查路由"})
		return
	}

	photosResp, err := h.s.GetPhotos(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Err.Error()})
	} else {
		c.JSON(http.StatusOK, photosResp)
	}
}

func (h *Handler) DeleteImg(c *gin.Context) {
	id := c.Param("id")
	user_id := c.GetString("user_id")
	if user_id == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "用户状态信息查询出错，请重试"})
	}

	if id == "" {
		log.Printf("获取用户信息的请求中id为'%s'\n", id)
		c.JSON(http.StatusBadRequest, gin.H{"err": "缺少id参数，检查路由"})
		return
	}

	id_num, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "服务器错误"})
		return
	}

	errs := h.s.DeleteImg(uint(id_num), user_id)
	if errs != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": errs.Err.Error()})
	} else {
		c.JSON(http.StatusCreated, gin.H{"msg": "文件已删除"})
	}
}

func (h *Handler) SetAvatar(c *gin.Context) {
	user_id := c.GetString("user_id")
	if user_id == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "用户状态信息查询出错，请重试"})
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "image part not found.",
		})
		return
	}

	imageType := c.PostForm("type")
	log.Printf("debug，%v\n", imageType)
	if imageType != "" {
		if !h.s.TypeAllow(imageType) {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": "仅支持PNG/JPG/JPEG/GIF格式",
			})
			return
		}
	}

	filename, errs := h.s.SetAvatar(file, imageType, "images", user_id)
	if errs != nil {
		switch errs.Code {
		case 400:
			c.JSON(http.StatusBadRequest, gin.H{
				"err": errs.Msg,
			})
		case 500:
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": errs.Msg,
			})
		}
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"filename": filename,
		})
	}
}

// UploadImage的路由保护版
func (h *Handler) UploadImg(c *gin.Context) {
	user_id := c.GetString("user_id")
	if user_id == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "用户状态信息查询出错，请重试"})
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "image part not found.",
		})
		return
	}

	imageType := c.PostForm("type")
	log.Printf("debug，%v\n", imageType)
	if imageType != "" {
		if !h.s.TypeAllow(imageType) {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": "仅支持PNG/JPG/JPEG/GIF格式",
			})
			return
		}
	}

	filename, errs := h.s.SaveImgWithUser(file, imageType, "images", user_id)
	if errs != nil {
		switch errs.Code {
		case 400:
			c.JSON(http.StatusBadRequest, gin.H{
				"err": errs.Msg,
			})
		case 500:
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": errs.Msg,
			})
		}
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"filename": filename,
		})
	}
}
