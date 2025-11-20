package dtos

import (
	"log"

	dao "github.com/Jack-samu/the-blog-backend-gin.git/internal/DAO"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/models"
)

func ToPostListItem(post *models.Post) PostListItem {
	p := PostListItem{
		ArticleBasic: ArticleBasic{
			Id:        post.ID,
			Title:     post.Title,
			CreatedAt: post.CreatedAt.String(),
			UpdatedAt: post.UpdatedAt.String(),
		},
		Excerpt: post.Excerpt,
		Views:   uint(post.ViewsCnt),
		Likes:   uint(post.LikeCnt),
		Cover:   post.Cover,
		Author:  post.Author.Username,
		// comment数量
		Comments: len(post.Comments),
	}

	for _, tag := range post.Tags {
		p.Tags = append(p.Tags, tag.Name)
	}

	return p
}

func ToArticleBasic(post *models.Post) ArticleBasic {
	article := ArticleBasic{
		Id:        post.ID,
		Title:     post.Title,
		CreatedAt: post.CreatedAt.String(),
		UpdatedAt: post.UpdatedAt.String(),
	}

	return article
}

func ToDraftsItem(draft *models.Draft) DraftItem {
	d := DraftItem{
		ArticleBasic: ArticleBasic{
			Id:        draft.ID,
			Title:     draft.Title,
			CreatedAt: draft.CreatedAt.String(),
			UpdatedAt: draft.UpdatedAt.String(),
		},
		Excerpt: draft.Excerpt,
		Cover:   draft.Cover,
		Author:  draft.Author.Username,
	}

	if draft.Category != nil {
		log.Printf("%v\n", draft.Category)
		d.Category = draft.Category.Name
	}

	for _, tag := range draft.Tags {
		d.Tags = append(d.Tags, tag.Name)
	}

	return d
}

func ToPostDetail(post *models.Post) PostDetailItem {
	p := PostDetailItem{
		PostListItem: PostListItem{
			ArticleBasic: ArticleBasic{
				Id:        post.ID,
				Title:     post.Title,
				CreatedAt: post.CreatedAt.String(),
				UpdatedAt: post.UpdatedAt.String(),
			},
			Excerpt:  post.Excerpt,
			Cover:    post.Cover,
			Views:    uint(post.ViewsCnt),
			Likes:    uint(post.LikeCnt),
			Comments: len(post.Comments),
		},
		Author: AuthorProfile{
			ID:       post.Author.ID,
			Username: post.Author.Username,
			Avatar:   post.Author.Avatar,
		},
		Content: post.Content,
	}

	if post.Category != nil {
		log.Printf("%v\n", post.Category)
		p.Category = post.Category.Name
	}

	if post.Tags != nil {
		for _, tag := range post.Tags {
			p.Tags = append(p.Tags, tag.Name)
		}
	}

	return p
}

func ToDraftDetail(draft *models.Draft) DraftDetail {
	d := DraftDetail{
		ArticleBasic: ArticleBasic{
			Id:        draft.ID,
			Title:     draft.Title,
			CreatedAt: draft.CreatedAt.String(),
			UpdatedAt: draft.UpdatedAt.String(),
		},
		Excerpt: draft.Excerpt,
		Cover:   draft.Cover,
		Content: draft.Content,
		Author: AuthorProfile{
			ID:       draft.Author.ID,
			Username: draft.Author.Username,
			Avatar:   draft.Author.Avatar,
		},
	}

	if draft.Category != nil {
		log.Printf("%v\n", draft.Category)
		d.Category = draft.Category.Name
	}

	return d
}

func ToPostList(posts []models.Post) []PostListItem {
	list := make([]PostListItem, len(posts))
	for i := range list {
		list[i] = ToPostListItem(&posts[i])
	}

	return list
}

func ToPostsBasic(posts []models.Post) []ArticleBasic {
	list := make([]ArticleBasic, len(posts))
	for i := range list {
		list[i] = ToArticleBasic(&posts[i])
	}

	return list
}

func ToDraftList(drafts []models.Draft) []DraftItem {
	list := make([]DraftItem, len(drafts))
	for i := range drafts {
		list[i] = ToDraftsItem(&drafts[i])
	}

	return list
}

func ToCommentsResp(comments []*models.Comment, r *dao.DAO) *CommentsResp {
	resp := &CommentsResp{
		Total: len(comments),
	}
	resp.Comments = make([]CommentItem, len(comments))
	for i := range comments {
		resp.Comments[i] = ToCommentItem(comments[i], nil, r)
	}

	return resp
}

func ToCommentItem(comment *models.Comment, user *models.User, r *dao.DAO) CommentItem {
	if user == nil {
		user = &comment.User
	}

	c := CommentItem{
		ID:        comment.ID,
		Content:   comment.Content,
		LikeCnt:   comment.LikeCnt,
		UpdatedAt: comment.UpdatedAt.GoString(),
		PostID:    comment.PostID,
		Replies:   len(comment.Replies),
		Commenter: AuthorProfile{
			ID:       user.ID,
			Username: user.Username,
			Avatar:   user.Avatar,
		},
		Liked: false,
	}

	if r != nil {
		// 直接查询
		isLiked, err := r.IsLiked(comment.UserID, "comment", comment.ID)
		if err != nil {
			c.Liked = false
			return c
		}
		c.Liked = isLiked
	}

	return c
}

func ToRepliesResp(replies []*models.Reply, r *dao.DAO) *RepliesResp {

	resp := &RepliesResp{
		Total: len(replies),
	}

	resp.Replies = make([]ReplyItem, len(replies))
	for i := range replies {
		resp.Replies[i] = ToReplyItem(replies[i], nil, r)
	}

	return resp
}

func ToReplyItem(reply *models.Reply, user *models.User, r *dao.DAO) ReplyItem {
	if user == nil {
		user = &reply.User
	}

	replyItem := ReplyItem{
		ID:        reply.ID,
		Content:   reply.Content,
		LikeCnt:   reply.LikeCnt,
		UpdatedAt: reply.UpdatedAt.GoString(),
		CommentID: reply.CommentID,
		Commenter: AuthorProfile{
			ID:       user.ID,
			Username: user.Username,
			Avatar:   user.Avatar,
		},
		Liked: false,
	}

	if r != nil {
		// 直接查询
		isLiked, err := r.IsLiked(reply.UserID, "reply", reply.ID)
		if err != nil {
			replyItem.Liked = false
			return replyItem
		}
		replyItem.Liked = isLiked
	}

	return replyItem
}
