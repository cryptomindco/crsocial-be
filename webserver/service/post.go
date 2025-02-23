package service

import (
	"fmt"
	"socialat/be/storage"
	"socialat/be/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

func (s *Service) CreateNewPost(userId uint64, username string, postContent string, fileType string, fileName []string) (*storage.Post, error) {
	post := storage.Post{
		UserId:    userId,
		Username:  username,
		Content:   postContent,
		CreatedAt: time.Now(),
	}
	post.UpdatedAt = post.CreatedAt
	if !utils.IsEmpty(fileType) {
		fileMainType := utils.GetFileType(fileType)
		filesStr := strings.Join(fileName, utils.FileSeparate)
		if utils.IsVideoFileType(fileMainType) {
			post.VideoUrl = filesStr
		} else if utils.IsImageFileType(fileMainType) {
			post.ImageUrl = filesStr
		} else {
			post.OtherUrl = filesStr
		}
	}
	// create new post
	err := s.db.Create(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *Service) GetAuthorByUsername(username string) (*storage.Author, error) {
	userInfo, err := s.GetUserInfoByUsername(username)
	if err != nil {
		return nil, err
	}
	author := storage.Author{
		Username: userInfo.Username,
		FullName: userInfo.FullName,
		Avatar:   userInfo.Avatar,
	}
	return &author, nil
}

func (s *Service) GetPostListByUsername(username string) ([]storage.Post, error) {
	var posts []storage.Post
	if err := s.db.Where("username = ?", username).Order("created_at desc").Find(&posts).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return make([]storage.Post, 0), nil
		}
		log.Error("GetPdsUserByUsername:get user info fail with error: ", err)
		return nil, err
	}
	return posts, nil
}

func (s *Service) GetAllPost(limit int) ([]storage.Post, error) {
	var posts = make([]storage.Post, 0)
	if err := s.db.Order("created_at desc").Limit(limit).Find(&posts).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return make([]storage.Post, 0), nil
		}
		log.Error("GetPdsUserByUsername:get user info fail with error: ", err)
		return make([]storage.Post, 0), err
	}
	return posts, nil
}

func (s *Service) GetPostOfFollowing(username string) ([]storage.Post, error) {
	var posts []storage.Post
	query := fmt.Sprintf(`SELECT * FROM posts WHERE username IN (SELECT target FROM followers WHERE follower = '%s' AND status) ORDER BY created_at desc`, username)
	err := s.db.Raw(query).Scan(&posts).Error
	if err != nil {
		return make([]storage.Post, 0), err
	}
	return posts, nil
}

func (s *Service) HandleForPostList(siteRoot, loginUser string, postViews []storage.PostView) []storage.PostView {
	postViews = storage.HandlerPostFileUrls(siteRoot, postViews)
	if utils.IsEmpty(loginUser) {
		return postViews
	}
	// check like
	for index, postView := range postViews {
		liked, _ := s.CheckLoginLiked(loginUser, postView.Id)
		postView.Liked = liked
		postViews[index] = postView
	}
	return postViews
}

func (s *Service) HandleForPost(siteRoot, loginUser string, postView storage.PostView) storage.PostView {
	postView = storage.HandlerOnePostFileUrls(siteRoot, postView)
	if utils.IsEmpty(loginUser) {
		return postView
	}
	// check like
	liked, _ := s.CheckLoginLiked(loginUser, postView.Id)
	postView.Liked = liked
	return postView
}

func (s *Service) CheckLoginLiked(username string, postId uint64) (bool, error) {
	var liked bool
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM likes WHERE username = '%s' AND post_id = %d)`, username, postId)
	err := s.db.Raw(query).Scan(&liked).Error
	if err != nil {
		return false, err
	}
	return liked, nil
}
