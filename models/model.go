package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"os"
	"time"
)

type Users struct {
	ResponseUsers
	Password string `json:"password"`
}

type SuccessStatus struct {
	Success    bool `json:"success"`
	StatusCode int  `json:"status_code"`
}

type ResponseUsersSignUp struct {
	ID    uint   `gorm:"primary_key" json:"id_user"`
	Email string `gorm:"unique_index" json:"email"`
	Name  string `json:"name"`
}

type ResponseUsers struct {
	ResponseUsersSignUp
	Token        string    `json:"token"`
	Address      string    `json:"address"`
	Bio          string    `json:"bio"`
	IDAvatar     int       `json:"id_avatar"`
	IDCoverPhoto int       `json:"id_cover_photo"`
	CreatedAt    time.Time `json:"published_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Response struct {
	Message    string      `json:"message"`
	Success    bool        `json:"success"`
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data, omitempty"`
}

type Status struct {
	Message    string `json:"message"`
	Success    bool   `json:"success"`
	StatusCode int    `json:"status_code"`
}

type Login struct {
	Email    string `gorm:"unique_index" json:"email"`
	Password string `json:"password"`
}

type Categories struct {
	Categories   string `gorm:"primary_key" json:"categories"`
	IDBackground uint   `json:"id_background"`
}

type Posts struct {
	ID           uint      `gorm:"primary_key" json:"id_post"`
	IDUser       uint      `json:"id_user"`
	Address      string    `json:"address"`
	IDAvatar     int       `json:"id_avatar"`
	IDBackground uint      `json:"id_background"`
	PostTitle    string    `json:"post_title"`
	Categories   string    `json:"categories"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"published_at"`
}

type PostsUsersJoin struct {
	ID           uint      `gorm:"primary_key" json:"id_post"`
	IDUser       uint      `json:"id_user"`
	IDBackground uint      `json:"id_background"`
	PostTitle    string    `json:"post_title"`
	Categories   string    `json:"categories"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"published_at"`
}

type Bookmarks struct {
	ID        uint      `gorm:"primary_key" json:"id_bookmark"`
	IDUser    uint      `json:"id_user"`
	IDPost    uint      `json:"id_post"`
	CreatedAt time.Time `json:"published_at"`
}

type Authentication struct {
	Email    string `gorm:"primary_key" json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ResponseUser struct {
	Users   interface{} `json:"users, omitempty"`
	Message string      `json:"message"`
}

type ResponseAuth struct {
	Authentication interface{} `json:"data", omitempty"`
	Message        string      `json:"message"`
}

type ResponseBookmark struct {
	Bookmarks interface{} `json:"bookmarks, omitempty"`
	Message   string      `json:"message"`
}

type ResponsePost struct {
	Posts   interface{} `json:"posts, omitempty"`
	Message string      `json:"message"`
	SuccessStatus
}

type ResponseCategory struct {
	Categories interface{} `json:"categories, omitempty"`
	Message    string      `json:"message"`
}

var db *gorm.DB
var err error

func init() {

	db_url := os.Getenv("DATABASE_URL")
	if db_url == "" {
		db_url = "host=localhost user=postgres dbname=gorm sslmode=disable password=postgres"
	}

	db, err = gorm.Open("postgres", db_url)

	db.SingularTable(true)

	if err != nil {
		fmt.Println(err)
	}

	db.AutoMigrate(&Users{}, &Bookmarks{}, &Posts{}, &Categories{})

}
