package storage

import (
	"socialat/be/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

type PostStorage interface {
	CreatePost(post *Post) error
	UpdatePost(post *Post) error
}

type Post struct {
	Id           uint64    `json:"id" gorm:"primarykey"`
	UserId       uint64    `json:"userId" gorm:"index:post_user_id_idx"`
	Username     string    `json:"username"`
	Content      string    `json:"content"`
	ImageUrl     string    `json:"imageUrl"`
	VideoUrl     string    `json:"videoUrl"`
	OtherUrl     string    `json:"otherUrl"`
	LikeCount    uint64    `json:"likeCount"`
	CommentCount uint64    `json:"commentCount"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func HandlerPostFileUrls(siteRoot string, posts []PostView) []PostView {
	for index, post := range posts {
		if !utils.IsEmpty(post.ImageUrl) {
			urlArr := strings.Split(post.ImageUrl, utils.FileSeparate)
			post.ImageUrls = urlArr
		}
		posts[index] = post
	}
	return posts
}

func HandlerOnePostFileUrls(siteRoot string, post PostView) PostView {
	if !utils.IsEmpty(post.ImageUrl) {
		urlArr := strings.Split(post.ImageUrl, utils.FileSeparate)
		post.ImageUrls = urlArr
	}
	return post
}

type Author struct {
	Username string `json:"username"`
	FullName string `json:"fullName"`
	Avatar   string `json:"avatar"`
}

type PostView struct {
	*Post
	ImageUrls []string `json:"imageUrls"`
	Author    *Author  `json:"author"`
	Liked     bool     `json:"liked"`
}

type PostFilter struct {
	Sort
	UserId   uint64
	Username string
}

func (f *PostFilter) BindQuery(db *gorm.DB) *gorm.DB {
	db = f.Sort.BindQuery(db)
	return db
}

func (f *PostFilter) BindCount(db *gorm.DB) *gorm.DB {
	return db
}

func (f *PostFilter) BindFirst(db *gorm.DB) *gorm.DB {
	return db
}

func (f *PostFilter) Sortable() map[string]bool {
	return map[string]bool{
		"createdAt": true,
	}
}

func (p *psql) CreatePost(post *Post) error {
	return p.db.Create(post).Error
}

func (p *psql) UpdatePost(post *Post) error {
	return p.db.Save(post).Error
}
