package service

import (
	"slices"
	"socialat/be/storage"
	"socialat/be/utils"

	"gorm.io/gorm"
)

func (s *Service) GetCommentViewList(loginUser string, postId uint64) ([]*storage.CommentView, error) {
	var comments []storage.Comment
	if err := s.db.Where("post_id = ?", postId).Order("created_at desc").Find(&comments).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Error("GetPdsUserByUsername:get user info fail with error: ", err)
			return nil, err
		}
	}
	parentCommentViews := make([]*storage.CommentView, 0)
	authorMap := make(map[string]*storage.Author)
	for _, comment := range comments {
		if comment.ParentId <= 0 {
			var parentView *storage.CommentView
			parentView, authorMap = s.GetCommentView(loginUser, comment, authorMap)
			childs := make([]*storage.CommentView, 0)
			// get child
			for _, child := range comments {
				if child.ParentId == comment.Id {
					var childView *storage.CommentView
					childView, authorMap = s.GetCommentView(loginUser, child, authorMap)
					childs = append(childs, childView)
				}
			}
			parentView.Childs = childs
			slices.Reverse(parentView.Childs)
			parentCommentViews = append(parentCommentViews, parentView)
		}
	}
	return parentCommentViews, nil
}

func (s *Service) UpdateCommentLikeCount(commentId uint32, liked bool) error {
	var comment storage.Comment
	err := s.db.First(&comment, commentId).Error
	if err != nil {
		return err
	}
	if liked {
		comment.LikeCount++
	} else {
		comment.LikeCount--
	}
	err = s.db.Save(&comment).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetCommentView(loginUser string, comment storage.Comment, authorMap map[string]*storage.Author) (*storage.CommentView, map[string]*storage.Author) {
	author, exist := authorMap[comment.Username]
	if !exist {
		var err error
		author, err = s.GetAuthorByUsername(comment.Username)
		if err == nil {
			authorMap[comment.Username] = author
		}
	}
	// check comment like for loginUser
	liked := false
	if !utils.IsEmpty(loginUser) {
		liked, _ = s.CheckLoginLiked(loginUser, comment.Id, int(utils.CommentType))
	}
	commentView := &storage.CommentView{
		Id:        comment.Id,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		LikeCount: comment.LikeCount,
		Author:    author,
		ParentId:  comment.ParentId,
		Liked:     liked,
	}
	return commentView, authorMap
}
