package storage

import "time"

type CommentStorage interface {
	CreateComment(comment *Comment) error
	UpdateComment(comment *Comment) error
}

type Comment struct {
	Id        uint64    `json:"id" gorm:"primarykey"`
	PostId    uint64    `json:"postId" gorm:"index:follower_follower_id_idx"`
	UserId    uint64    `json:"userId" gorm:"index:follower_target_id_idx"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

func (p *psql) CreateComment(comment *Comment) error {
	return p.db.Create(comment).Error
}

func (p *psql) UpdateComment(comment *Comment) error {
	return p.db.Save(comment).Error
}
