package storage

import (
	"time"

	"gorm.io/gorm"
)

const UserFieldUName = "user_name"
const UserFieldId = "id"
const RecipientId = "recipient_id"

type AuthType int

const (
	AuthLocalUsernamePassword AuthType = iota
	AuthMicroservicePasskey
)

type UserInfoStorage interface {
	CreateUserInfo(user *UserInfo) error
	UpdateUserInfo(user *UserInfo) error
}

type UserInfo struct {
	Id        uint64    `json:"id" gorm:"primarykey"`
	Username  string    `json:"username" gorm:"index:user_info_username_idx,unique"`
	Email     string    `json:"email"`
	FullName  string    `json:"fullName"`
	Bio       string    `json:"bio"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"createdAt"`
}
type AuthClaims struct {
	Id        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FullName  string `json:"fullName"`
	Bio       string `json:"bio"`
	Avatar    string `json:"avatar"`
	LoginType int    `json:"loginType"`
	Expire    int64  `json:"expire"`
	Role      int    `json:"role"`
	CreatedAt int64  `json:"createdAt"`
	LastLogin int64  `json:"lastLogin"`
}

func (p *psql) CreateUserInfo(user *UserInfo) error {
	return p.db.Create(user).Error
}

func (p *psql) UpdateUserInfo(user *UserInfo) error {
	return p.db.Save(user).Error
}

type UserInfoFilter struct {
	Sort
	Username string
}

func (f *UserInfoFilter) BindQuery(db *gorm.DB) *gorm.DB {
	db = f.Sort.BindQuery(db)
	return db
}

func (f *UserInfoFilter) BindCount(db *gorm.DB) *gorm.DB {
	return db
}

func (f *UserInfoFilter) BindFirst(db *gorm.DB) *gorm.DB {
	return db
}

func (f *UserInfoFilter) Sortable() map[string]bool {
	return map[string]bool{
		"username": true,
	}
}
