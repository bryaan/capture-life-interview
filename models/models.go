package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

var DB *gorm.DB

// UrlName - name of the blog post to show in browser url

type BlogPost struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	UrlName   string    `json:"urlName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// parent - Parent Blog Post
// author - Author of comment
// content - Content of comment
// DateCreated - timestamp of when comment was added

type ParentType int

const (
	PT_BlogPost ParentType = 1
	PT_Comment             = 2
)

type Comment struct {
	ID          uint       `json:"id" gorm:"primary_key"`
	ParentId    uint       `json:"parentId"`
	ParentType  ParentType `json:"parentType"`
	Author      string     `json:"author"`
	Content     string     `json:"content"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	SubComments []Comment  `json:"subComments"`
}

func ConnectDatabase() {
	database, err := gorm.Open("sqlite3", "test.db")

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&BlogPost{})
	database.AutoMigrate(&Comment{})

	DB = database
}
