package comment

import (
	"net/http"
	"testing"

	dao "github.com/Jack-samu/the-blog-backend-gin.git/internal/DAO"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/dtos"
	"github.com/Jack-samu/the-blog-backend-gin.git/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestCommentServiceFlow(t *testing.T) {
	db := setupTestDB(t)
	repo := dao.NewRepository(db)
	serv := service.NewService(repo)
	defer teardownTestDB(db)

	// 查询不存在的post评论
	commentsResp, err := serv.GetComments(int64(1))
	assert.Equal(t, http.StatusNotFound, err.Code)
	assert.Nil(t, commentsResp)

	// 查询post评论但传入错误参数
	commentsResp, err = serv.GetComments(int64(-1))
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Nil(t, commentsResp)

	// 查询reply评论但传入错误参数
	repliesResp, err := serv.GetReplies(int64(-1))
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Nil(t, repliesResp)

	//准备文章主体和用户主体
	userID, postID := preparation(serv, t)

	// 创建评论
	req := &dtos.CommentReq{
		ArticleID: int64(postID),
		Content:   "comment测试",
	}
	t.Logf("comment请求参数：%v\n", req)
	commentResp, err := serv.CreateComment(req, userID)
	assert.Nil(t, err)
	assert.Equal(t, uint(postID), commentResp.Comment.PostID)
	assert.Equal(t, userID, commentResp.Comment.Commenter.ID)
	// 查comments
	commentsResp, err = serv.GetComments(int64(postID))
	assert.Nil(t, err)
	assert.Equal(t, int(1), len(commentsResp.Comments))
	assert.Equal(t, "comment测试", commentsResp.Comments[0].Content)
	// 查replies
	repliesResp, err = serv.GetReplies(int64(commentResp.Comment.ID))
	assert.Nil(t, err)
	assert.Empty(t, repliesResp.Replies)

	// 修改comment
	req = &dtos.CommentReq{
		Content:   "comment测试，改",
		CommentID: int64(commentResp.Comment.ID),
	}
	t.Logf("comment请求参数：%v\n", req)
	commentResp, err = serv.ModifyComment(req, userID)
	assert.Nil(t, err)
	assert.Equal(t, "comment测试，改", commentResp.Comment.Content)
	// 查comments
	commentsResp, err = serv.GetComments(int64(postID))
	assert.Nil(t, err)
	assert.Equal(t, int(1), len(commentsResp.Comments))
	assert.Equal(t, "comment测试，改", commentsResp.Comments[0].Content)

	// 创建reply
	req = &dtos.CommentReq{
		ArticleID: 0,
		Content:   "reply测试",
		CommentID: int64(commentResp.Comment.ID),
	}
	t.Logf("reply请求参数：%v\n", req)
	replyResp, err := serv.CreateReply(req, userID)
	assert.Nil(t, err)
	assert.Equal(t, "reply测试", replyResp.Reply.Content)
	// 查询replies
	repliesResp, err = serv.GetReplies(int64(commentResp.Comment.ID))
	assert.Nil(t, err)
	assert.Equal(t, int(1), len(repliesResp.Replies))
	assert.Equal(t, "reply测试", repliesResp.Replies[0].Content)

	// 修改reply
	req = &dtos.CommentReq{
		ReplyID:   replyResp.Reply.ID,
		Content:   "reply测试，改",
		CommentID: int64(commentResp.Comment.ID),
	}
	t.Logf("reply请求参数：%v\n", req)
	replyResp, err = serv.ModifyReply(req, userID)
	assert.Nil(t, err)
	assert.Equal(t, "reply测试，改", replyResp.Reply.Content)
	// 查询replies
	repliesResp, err = serv.GetReplies(int64(commentResp.Comment.ID))
	assert.Nil(t, err)
	assert.Equal(t, int(1), len(repliesResp.Replies))
	assert.Equal(t, "reply测试，改", repliesResp.Replies[0].Content)

	// 删除comment
	err = serv.DeleteComment(int64(commentResp.Comment.ID), userID)
	assert.Nil(t, err)
	// 查replies
	repliesResp, err = serv.GetReplies(int64(commentResp.Comment.ID))
	assert.Equal(t, http.StatusNotFound, err.Code)
	assert.Nil(t, repliesResp)
	// 查comments
	commentsResp, err = serv.GetComments(int64(postID))
	assert.Nil(t, err)
	assert.Equal(t, int(0), len(commentsResp.Comments))
}
