package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model           // this includes createAt, updatedAt, ID, DeletedAt
	FirstName  string    `gorm:"size:60" json:"first_name" validate:"required"`
	LastName   string    `gorm:"size:60" json:"last_name" validate:"required"`
	Email      *string   `gorm:"unique;size:60" json:"email" validate:"email,required"`
	Password   string    `gorm:"size:250;<-:create" json:"-" validate:"required"`
	Phone      string    `gorm:"unique;size:16" json:"phone" validate:"required"`
	Username   string    `gorm:"unique;size:60" json:"username" validate:"required"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	Comments   []Comment `json:"comments"`
	Articles   []Article `json:"articles"`
}

type Category struct {
	gorm.Model
	CategoryName string `json:"category_name"`
}

type Article struct {
	gorm.Model
	Body       string     `json:"body"`
	Likes      []User     `gorm:"many2many:article_users_likes;" json:"likes"`
	Readers    []User     `gorm:"many2many:article_users_readers;" json:"readers"`
	Categories []Category `gorm:"many2many:category_articles;" json:"categories"`
	UserId     int        `json:"user_id"`
	User       User       `json:"user"`
}

type Comment struct {
	gorm.Model
	Body   string `json:"body"`
	Likes  []User `gorm:"many2many:article_users_likes;" json:"likes"`
	UserId int    `json:"user_id"`
	User   User   `json:"user"`
}
