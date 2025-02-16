package storage

import "time"

type FollowerStorage interface {
	CreateFollower(follower *Follower) error
	UpdateFollower(follower *Follower) error
}

type Follower struct {
	Id         uint64    `json:"id" gorm:"primarykey"`
	Follower   string    `json:"follower" gorm:"not null;index"`
	Target     string    `json:"target" gorm:"not null;index"`
	Status     bool      `json:"status"`
	FollowedAt time.Time `json:"followedAt"`

	_ struct{} `gorm:"uniqueIndex:idx_follower_follow_target"`
}

type FollowingDisplay struct {
	Target   string `json:"target"`
	Avatar   string `json:"avatar"`
	FullName string `json:"fullName"`
}

func (p *psql) CreateFollower(follower *Follower) error {
	return p.db.Create(follower).Error
}

func (p *psql) UpdateFollower(follower *Follower) error {
	return p.db.Save(follower).Error
}
