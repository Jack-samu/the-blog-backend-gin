package handler

import (
	"net/http"
	"strconv"

	"github.com/Jack-samu/the-blog-backend-gin.git/internal/dtos"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetComments(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "服务器错误"})
		return
	}

	resp, errs := h.s.GetComments(id)
	if errs != nil {
		switch errs.Code {
		case http.StatusBadRequest:
			c.JSON(http.StatusBadRequest, gin.H{"err": errs.Msg})
		case http.StatusNotFound:
			c.JSON(http.StatusInternalServerError, gin.H{"err": "没有找到对应的文章，无法刷新评论"})
		case http.StatusInternalServerError:
			c.JSON(http.StatusInternalServerError, gin.H{"err": "服务器错误：" + errs.Err.Error()})
		}
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func (h *Handler) GetReplies(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "服务器错误"})
		return
	}

	resp, errs := h.s.GetReplies(id)
	if errs != nil {
		switch errs.Code {
		case http.StatusBadRequest:
			c.JSON(http.StatusBadRequest, gin.H{"err": errs.Msg})
		case http.StatusNotFound:
			c.JSON(http.StatusInternalServerError, gin.H{"err": "没有找到对应的评论，无法加载"})
		case http.StatusInternalServerError:
			c.JSON(http.StatusInternalServerError, gin.H{"err": "服务器错误：" + errs.Err.Error()})
		}
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func (h *Handler) CreateComment(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "用户id读取失败"})
	}

	var req dtos.CommentReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "无效参数"})
		return
	}

	if req.ArticleID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "要进行评论的文章ID为空"})
		return
	}

	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "评论主体内容不能为空"})
		return
	}

	commentResp, errs := h.s.CreateComment(&req, userID)
	if errs != nil {
		switch errs.Code {
		case http.StatusNotFound:
			c.JSON(http.StatusBadRequest, gin.H{"err": "主体文章没找到"})
		case http.StatusInternalServerError:
			c.JSON(http.StatusInternalServerError, gin.H{"err": errs.Err.Error()})
		}
	} else {
		c.JSON(http.StatusCreated, commentResp)
	}
}

func (h *Handler) CreateReply(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "用户id读取失败"})
	}

	var req dtos.CommentReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "无效参数"})
		return
	}

	if req.CommentID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "要进行评论的评论ID为空"})
		return
	}

	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "评论主体内容不能为空"})
		return
	}

	replyResp, errs := h.s.CreateReply(&req, userID)
	if errs != nil {
		switch errs.Code {
		case http.StatusNotFound:
			c.JSON(http.StatusBadRequest, gin.H{"err": "主体文章没找到"})
		case http.StatusInternalServerError:
			c.JSON(http.StatusInternalServerError, gin.H{"err": errs.Err.Error()})
		}
	} else {
		c.JSON(http.StatusCreated, replyResp)
	}
}

func (h *Handler) ModifyComment(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "用户id读取失败"})
	}

	var req dtos.CommentReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "无效参数"})
		return
	}

	if req.ArticleID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "文章ID为空"})
		return
	}

	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "评论主体内容不能为空"})
		return
	}

	commentResp, errs := h.s.ModifyComment(&req, userID)
	if errs != nil {
		switch errs.Code {
		case http.StatusNotFound:
			c.JSON(http.StatusBadRequest, gin.H{"err": "主体文章没找到"})
		case http.StatusInternalServerError:
			c.JSON(http.StatusInternalServerError, gin.H{"err": errs.Err.Error()})
		}
	} else {
		c.JSON(http.StatusCreated, commentResp)
	}
}

func (h *Handler) ModifyReply(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "用户id读取失败"})
	}

	var req dtos.CommentReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "无效参数"})
		return
	}

	if req.ReplyID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "评论ID为空"})
		return
	}

	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "评论主体内容不能为空"})
		return
	}

	commentResp, errs := h.s.ModifyReply(&req, userID)
	if errs != nil {
		switch errs.Code {
		case http.StatusNotFound:
			c.JSON(http.StatusBadRequest, gin.H{"err": "主体文章没找到"})
		case http.StatusInternalServerError:
			c.JSON(http.StatusInternalServerError, gin.H{"err": errs.Err.Error()})
		}
	} else {
		c.JSON(http.StatusCreated, commentResp)
	}
}

func (h *Handler) DeleteComment(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "无效参数"})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "用户id解析出错"})
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "提取评论id出错"})
		return
	}

	errs := h.s.DeleteComment(id, userID)
	if errs != nil {
		switch errs.Code {
		case http.StatusNotFound:
			c.JSON(http.StatusInternalServerError, gin.H{"err": "没找到该评论，请重试"})
		case http.StatusInternalServerError:
			c.JSON(http.StatusInternalServerError, gin.H{"err": errs.Err.Error()})
		}
	} else {
		c.JSON(http.StatusCreated, gin.H{"msg": "评论已删除"})
	}
}

func (h *Handler) DeleteReply(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "无效参数"})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "用户id解析出错"})
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "提取评论id出错"})
		return
	}

	errs := h.s.DeleteReply(id, userID)
	if errs != nil {
		switch errs.Code {
		case http.StatusNotFound:
			c.JSON(http.StatusInternalServerError, gin.H{"err": "没找到该评论，请重试"})
		case http.StatusInternalServerError:
			c.JSON(http.StatusInternalServerError, gin.H{"err": errs.Err.Error()})
		}
	} else {
		c.JSON(http.StatusCreated, gin.H{"msg": "评论已删除"})
	}
}
