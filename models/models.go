package models

import (
	"strconv"

	"gorm.io/gorm"
)

//UserAuth 用户授权表
type UserAuth struct {
	gorm.Model
	Uid          uint64 `gorm:"uniqueIndex:idx_uid"`                 // uid
	IdentityType int    `gorm:"uniqueIndex:idx_uid,idx_identifier"`  // 1用户名 2邮箱 3手机号 4qq 5微信
	Identifier   string `gorm:"uniqueIndex:idx_identifier;size:128"` // 手机号 邮箱 用户名或第三方应用的唯一标识
	Certificate  string `gorm:"size:128"`                            // 密码凭证(站内的保存密码，站外的不保存或保存token)
	//CertExpireAt time.Time
	RefreshToken string `gorm:"size:128"`
	Openid       string `gorm:"index;size:128"` // wx openid
}

// table posts
type Post struct {
	gorm.Model
	Title       string // title
	Body        string // body
	View        int    // view count
	IsPublished bool   // published or not
	Tags        []*Tag `gorm:"-"` // tags of post
}

// table tags
type Tag struct {
	gorm.Model
	Name  string // tag name
	Total int    `gorm:"-"` // count of post
}

// table post_tags
type PostTag struct {
	gorm.Model
	PostId uint // post id
	TagId  uint // tag id
}

// Insert -
func (post *Post) Insert(db *gorm.DB) error {
	return db.Create(post).Error
}

// Update -
func (post *Post) Update(db *gorm.DB) error {
	return db.Model(post).Updates(map[string]interface{}{
		"title":        post.Title,
		"body":         post.Body,
		"is_published": post.IsPublished,
	}).Error
}

//UpdateView -
func (post *Post) UpdateView(db *gorm.DB) error {
	return db.Model(post).Updates(map[string]interface{}{
		"view": post.View,
	}).Error
}

//Delete -
func (post *Post) Delete(db *gorm.DB) error {
	return db.Delete(post).Error
}

//GetPostById -
func GetPostById(db *gorm.DB, id string) (*Post, error) {
	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}

	var post Post
	err = db.First(&post, "id = ?", pid).Error
	return &post, err
}
