package models

import (
	"context"
	"strconv"

	"ginfra/datasource"
	
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	//_ "github.com/go-sql-driver/mysql"
)

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

// Post
func (post *Post) Insert(ctx context.Context) error {
	db, _ := datasource.GormWithContext(ctx)
	return db.Create(post).Error
}

func (post *Post) Update(ctx context.Context) error {
	db, _ := datasource.GormWithContext(ctx)
	return db.Model(post).Updates(map[string]interface{}{
		"title":        post.Title,
		"body":         post.Body,
		"is_published": post.IsPublished,
	}).Error
}

func (post *Post) UpdateView(ctx context.Context) error {
	db, _ := datasource.GormWithContext(ctx)
	return db.Model(post).Updates(map[string]interface{}{
		"view": post.View,
	}).Error
}

func (post *Post) Delete(ctx context.Context) error {
	db, _ := datasource.GormWithContext(ctx)
	return db.Delete(post).Error
}

func GetPostById(ctx context.Context, id string) (*Post, error) {
	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}

	db, _ := datasource.GormWithContext(ctx)
	var post Post
	err = db.First(&post, "id = ?", pid).Error
	return &post, err
}

// func (post *Post) Excerpt() template.HTML {
// 	//you can sanitize, cut it down, add images, etc
// 	policy := bluemonday.StrictPolicy() //remove all html tags
// 	sanitized := policy.Sanitize(string(blackfriday.MarkdownCommon([]byte(post.Body))))
// 	runes := []rune(sanitized)
// 	if len(runes) > 300 {
// 		sanitized = string(runes[:300])
// 	}
// 	excerpt := template.HTML(sanitized + "...")
// 	return excerpt
// }

// func ListPublishedPost(tag string, pageIndex, pageSize int) ([]*Post, error) {
// 	return _listPost(tag, true, pageIndex, pageSize)
// }
//
// func ListAllPost(tag string) ([]*Post, error) {
// 	return _listPost(tag, false, 0, 0)
// }
//
// func _listPost(tag string, published bool, pageIndex, pageSize int) ([]*Post, error) {
// 	var posts []*Post
// 	var err error
// 	if len(tag) > 0 {
// 		tagId, err := strconv.ParseUint(tag, 10, 64)
// 		if err != nil {
// 			return nil, err
// 		}
// 		var rows *sql.Rows
// 		if published {
// 			if pageIndex > 0 {
// 				rows, err = DB.Raw("select p.* from posts p inner join post_tags pt on p.id = pt.post_id where pt.tag_id = ? and p.is_published = ? order by created_at desc limit ? offset ?", tagId, true, pageSize, (pageIndex-1)*pageSize).Rows()
// 			} else {
// 				rows, err = DB.Raw("select p.* from posts p inner join post_tags pt on p.id = pt.post_id where pt.tag_id = ? and p.is_published = ? order by created_at desc", tagId, true).Rows()
// 			}
// 		} else {
// 			rows, err = DB.Raw("select p.* from posts p inner join post_tags pt on p.id = pt.post_id where pt.tag_id = ? order by created_at desc", tagId).Rows()
// 		}
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer rows.Close()
// 		for rows.Next() {
// 			var post Post
// 			DB.ScanRows(rows, &post)
// 			posts = append(posts, &post)
// 		}
// 	} else {
// 		if published {
// 			if pageIndex > 0 {
// 				err = DB.Where("is_published = ?", true).Order("created_at desc").Limit(pageSize).Offset((pageIndex - 1) * pageSize).Find(&posts).Error
// 			} else {
// 				err = DB.Where("is_published = ?", true).Order("created_at desc").Find(&posts).Error
// 			}
// 		} else {
// 			err = DB.Order("created_at desc").Find(&posts).Error
// 		}
// 	}
// 	return posts, err
// }
//
// func MustListMaxReadPost() (posts []*Post) {
// 	posts, _ = ListMaxReadPost()
// 	return
// }
//
// func ListMaxReadPost() (posts []*Post, err error) {
// 	err = DB.Where("is_published = ?", true).Order("view desc").Limit(5).Find(&posts).Error
// 	return
// }
//
// // Tag
// func (tag *Tag) Insert() error {
// 	return DB.FirstOrCreate(tag, "name = ?", tag.Name).Error
// }
//
// func ListTag() ([]*Tag, error) {
// 	var tags []*Tag
// 	rows, err := DB.Raw("select t.*,count(*) total from tags t inner join post_tags pt on t.id = pt.tag_id inner join posts p on pt.post_id = p.id where p.is_published = ? group by pt.tag_id", true).Rows()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var tag Tag
// 		DB.ScanRows(rows, &tag)
// 		tags = append(tags, &tag)
// 	}
// 	return tags, nil
// }
//
// func MustListTag() []*Tag {
// 	tags, _ := ListTag()
// 	return tags
// }
//
// func ListTagByPostId(id string) ([]*Tag, error) {
// 	var tags []*Tag
// 	pid, err := strconv.ParseUint(id, 10, 64)
// 	if err != nil {
// 		return nil, err
// 	}
// 	rows, err := DB.Raw("select t.* from tags t inner join post_tags pt on t.id = pt.tag_id where pt.post_id = ?", uint(pid)).Rows()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var tag Tag
// 		DB.ScanRows(rows, &tag)
// 		tags = append(tags, &tag)
// 	}
// 	return tags, nil
// }
//
// func CountTag() int {
// 	var count int
// 	DB.Model(&Tag{}).Count(&count)
// 	return count
// }
//
// func ListAllTag() ([]*Tag, error) {
// 	var tags []*Tag
// 	err := DB.Model(&Tag{}).Find(&tags).Error
// 	return tags, err
// }
//
// // post_tags
// func (pt *PostTag) Insert() error {
// 	return DB.FirstOrCreate(pt, "post_id = ? and tag_id = ?", pt.PostId, pt.TagId).Error
// }
//
// func DeletePostTagByPostId(postId uint) error {
// 	return DB.Delete(&PostTag{}, "post_id = ?", postId).Error
// }
