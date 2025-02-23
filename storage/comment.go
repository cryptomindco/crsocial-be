package storage

import "time"

type CommentStorage interface {
	CreateComment(comment *Comment) error
	UpdateComment(comment *Comment) error
}

type Comment struct {
	Id        uint64    `json:"id" gorm:"primarykey"`
	PostId    uint64    `json:"postId" gorm:"index:comment_post_id_idx"`
	Username  string    `json:"username" gorm:"index:comment_username_idx"`
	Content   string    `json:"content"`
	ParentId  uint64    `json:"parentId"`
	LikeCount uint64    `json:"likeCount"`
	CreatedAt time.Time `json:"createdAt"`
}

type CommentView struct {
	Id        uint64         `json:"id" gorm:"primarykey"`
	Author    *Author        `json:"author"`
	Content   string         `json:"content"`
	LikeCount uint64         `json:"likeCount"`
	ParentId  uint64         `json:"parentId"`
	Childs    []*CommentView `json:"childs"`
	Liked     bool           `json:"liked"`
	CreatedAt time.Time      `json:"createdAt"`
}

func (p *psql) CreateComment(comment *Comment) error {
	return p.db.Create(comment).Error
}

func (p *psql) UpdateComment(comment *Comment) error {
	return p.db.Save(comment).Error
}
