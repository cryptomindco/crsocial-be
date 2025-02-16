package service

import (
	"socialat/be/storage"

	"gorm.io/gorm"
)

func (s *Service) GetFollowByName(follower, target string) (*storage.Follower, error) {
	var follow storage.Follower
	if err := s.db.Where("follower = ? AND target = ?", follower, target).First(&follow).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &follow, nil
}

func (s *Service) GetFollowingOfUser(username string) ([]storage.Follower, error) {
	var follows []storage.Follower
	if err := s.db.Where("follower = ? AND status = ?", username, true).Find(&follows).Error; err != nil {
		return nil, err
	}
	return follows, nil
}
