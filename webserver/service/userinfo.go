package service

import (
	"fmt"
	"socialat/be/storage"
	"socialat/be/utils"

	"gorm.io/gorm"
)

func (s *Service) GetUserInfoByUsername(username string) (*storage.UserInfo, error) {
	var user storage.UserInfo
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.NewError(fmt.Errorf("user info not found"), utils.ErrorNotFound)
		}
		log.Error("GetPdsUserByUsername:get user info fail with error: ", err)
		return nil, err
	}
	return &user, nil
}
