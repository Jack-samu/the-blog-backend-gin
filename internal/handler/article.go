package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Jack-samu/the-blog-backend-gin.git/internal/dtos"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetArticles(c *gin.Context) {
	pageParam := c.Query("page")
	if pageParam == "" {
		pageParam = "1"
	}
	perPageParam := c.Query("per_page")
	if perPageParam == "" {
		perPageParam = "10"
	}

	page, _ := strconv.ParseInt(pageParam, 10, 64)
	perPage, _ := strconv.ParseInt(perPageParam, 10, 64)

	postsResp, err := h.s.GetPosts(page, perPage)
	if err != nil {
		if err.Err != nil {
			c.JSON(err.Code, gin.H{"err": err.Err.Error()})
		} else {
			c.JSON(err.Code, gin.H{"err": err.Msg})
		}
	} else {
		c.JSON(http.StatusOK, postsResp)
	}
}

func (h *Handler) GetArticle(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		// 没有明确指向就重定向回主页吧
		c.Redirect(http.StatusPermanentRedirect, "/")
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "获取指向文章出错"})
		return
	}

	postResp, errs := h.s.GetPost(uint(id))
	if errs != nil {
		if errs.Err != nil {
			c.JSON(errs.Code, gin.H{"err": errs.Err.Error()})
		} else {
			c.JSON(errs.Code, gin.H{"err": errs.Msg})
		}
	} else {
		c.JSON(http.StatusOK, postResp)
	}
}

func (h *Handler) GetDraftEditable(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		// 没有明确指向就重定向回主页吧
		c.Redirect(http.StatusPermanentRedirect, "/")
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "获取指向文章出错"})
		return
	}

	draftResp, errs := h.s.GetDraft(uint(id))
	if errs != nil {
		if errs.Err != nil {
			c.JSON(errs.Code, gin.H{"err": errs.Err.Error()})
		} else {
			c.JSON(errs.Code, gin.H{"err": errs.Msg})
		}
	} else {
		c.JSON(http.StatusOK, draftResp)
	}
}

func (h *Handler) GetPostsOfUser(c *gin.Context) {
	pageParam := c.Query("page")
	if pageParam == "" {
		pageParam = "1"
	}
	perPageParam := c.Query("per_page")
	if perPageParam == "" {
		perPageParam = "10"
	}

	user_id := c.GetString("user_id")
	if user_id == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "用户状态信息查询出错，请重试"})
		return
	}

	page, _ := strconv.ParseInt(pageParam, 10, 64)
	perPage, _ := strconv.ParseInt(perPageParam, 10, 64)

	postsResp, err := h.s.GetPostsOfUser(page, perPage, user_id)
	if err != nil {
		if err.Err != nil {
			c.JSON(err.Code, gin.H{"err": err.Err.Error()})
		} else {
			c.JSON(err.Code, gin.H{"err": err.Msg})
		}
	} else {
		c.JSON(http.StatusOK, postsResp)
	}
}

func (h *Handler) GetDraftOfUser(c *gin.Context) {
	pageParam := c.Query("page")
	if pageParam == "" {
		pageParam = "1"
	}
	perPageParam := c.Query("per_page")
	if perPageParam == "" {
		perPageParam = "10"
	}

	user_id := c.GetString("user_id")
	if user_id == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "用户状态信息查询出错，请重试"})
		return
	}

	page, _ := strconv.ParseInt(pageParam, 10, 64)
	perPage, _ := strconv.ParseInt(perPageParam, 10, 64)

	draftsResp, err := h.s.GetDraftsOfUser(page, perPage, user_id)
	if err != nil {
		if err.Err != nil {
			c.JSON(err.Code, gin.H{"err": err.Err.Error()})
		} else {
			c.JSON(err.Code, gin.H{"err": err.Msg})
		}
	} else {
		c.JSON(http.StatusOK, draftsResp)
	}
}

func (h *Handler) GetSeries(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "参数不明确"})
		return
	}

	postOfCategory, err := h.s.GetSeries(id)
	if err != nil {
		if err.Err != nil {
			c.JSON(err.Code, gin.H{"err": err.Err.Error()})
		} else {
			c.JSON(err.Code, gin.H{"err": err.Msg})
		}
	} else {
		c.JSON(http.StatusOK, postOfCategory)
	}
}

// 发布功能，重点
func (h *Handler) PublishArticle(c *gin.Context) {

	req := new(dtos.ArticleReq)
	if err := c.ShouldBindJSON(req); err != nil {
		log.Printf("%s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"err": "参数缺失或无效"})
		return
	}

	user_id := c.GetString("user_id")
	if user_id == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "用户状态信息查询出错，请重试"})
		return
	}

	post_id, err := h.s.PublishArticle(req, user_id)
	if err != nil {
		if err.Err != nil {
			c.JSON(err.Code, gin.H{"err": err.Err.Error()})
		}
	} else {
		c.JSON(http.StatusCreated, gin.H{"msg": "发表成功", "id": post_id})
	}
}

// 保存草稿
func (h *Handler) SaveDraft(c *gin.Context) {

	req := new(dtos.ArticleReq)
	if err := c.ShouldBindJSON(req); err != nil {
		log.Printf("%s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"err": "参数缺失或无效"})
		return
	}

	user_id := c.GetString("user_id")
	if user_id == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "用户状态信息查询出错，请重试"})
		return
	}

	draft_id, err := h.s.SaveDraft(req, user_id)
	if err != nil {
		if err.Err != nil {
			c.JSON(err.Code, gin.H{"err": err.Err.Error()})
		}
	} else {
		c.JSON(http.StatusCreated, gin.H{"msg": "草稿保存成功", "id": draft_id})
	}
}

func (h *Handler) DeletePost(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "id为空"})
		return
	}

	post_id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Printf("转换%s格式出错：%s\n", idParam, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"err": "服务器错误"})
	}

	if err := h.s.DeletePost(uint(post_id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "服务器错误"})
	} else {
		c.JSON(http.StatusCreated, gin.H{"msg": "文章已删除"})
	}
}

func (h *Handler) DeleteDraft(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "id为空"})
		return
	}

	draft_id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Printf("转换%s格式出错：%s\n", idParam, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"err": "服务器错误"})
	}

	if err := h.s.DeletePost(uint(draft_id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "服务器错误"})
	} else {
		c.JSON(http.StatusCreated, gin.H{"msg": "草稿已删除"})
	}
}
