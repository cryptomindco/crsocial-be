package storage

import "time"

type LikeStorage interface {
	CreateLike(like *Like) error
	UpdateLike(like *Like) error
}

type Like struct {
	Id       uint64    `json:"id" gorm:"primarykey"`
	Username string    `json:"Username" gorm:"not null;index"`
	PostId   uint64    `json:"postId" gorm:"not null;index"`
	LikedAt  time.Time `json:"likedAt"`

	_ struct{} `gorm:"uniqueIndex:idx_like_post_user,priority:1"`
}

func (p *psql) CreateLike(like *Like) error {
	return p.db.Create(like).Error
}

func (p *psql) UpdateLike(like *Like) error {
	return p.db.Save(like).Error
}
