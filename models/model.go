package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"os"
	"time"
)

type SuccessStatus struct {
	Message    string `json:"message"`
	Success    bool   `json:"success"`
	StatusCode int    `json:"status_code"`
}

type Users struct {
	ResponseUsers
	Password string `json:"password"`
}

type ResponseUsersSignUp struct {
	ID    uint   `gorm:"primary_key" json:"id_user"`
	Email string `gorm:"unique_index" json:"email"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type ResponseUsers struct {
	ResponseUsersSignUp
	Address      string    `json:"address"`
	Bio          string    `json:"bio"`
	IDAvatar     uint      `json:"id_avatar"`
	IDCoverPhoto uint      `json:"id_cover_photo"`
	TotalPosts   uint      `json:"total_posts"`
	CreatedAt    time.Time `json:"published_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UpdateUsers struct {
	Email        string `json:"email"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	Bio          string `json:"bio"`
	IDAvatar     uint   `json:"id_avatar"`
	IDCoverPhoto uint   `json:"id_cover_photo"`
}

type UpdatePosts struct {
	ID           uint      `gorm:"primary_key" json:"id_post"`
	IDBackground uint      `json:"id_background, omitempty"`
	PostTitle    string    `json:"post_title, omitempty"`
	Categories   string    `json:"categories, omitempty"`
	Content      string    `json:"content, omitempty"`
	CreatedAt    time.Time `json:"published_at, omitempty"`
	UpdatedAt    time.Time `json:"updated_at, omitempty"`
}

type ResponseLogin struct {
	Message    string `json:"message"`
	Success    bool   `json:"success"`
	StatusCode int    `json:"status_code"`
	Token      string `json:"token"`
	Email      string `json:"email"`
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

type UsersCategories struct {
	ID         uint   `gorm:"primary_key" json:"id"`
	IDUser     uint   `json:"id_user, omitempty"`
	Categories string `json:"categories"`
}

type Posts struct {
	ID           uint      `gorm:"primary_key" json:"id_post"`
	IDUser       uint      `json:"id_user, omitempty"`
	IDBackground uint      `json:"id_background, omitempty"`
	PostTitle    string    `json:"post_title, omitempty"`
	Categories   string    `json:"categories, omitempty"`
	Content      string    `json:"content, omitempty"`
	CreatedAt    time.Time `json:"published_at, omitempty"`
	UpdatedAt    time.Time `json:"updated_at, omitempty"`
}

type PostsUsersJoin struct {
	Posts
	UpdateUsers
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
	Message    string      `json:"message"`
	Success    bool        `json:"success"`
	StatusCode int         `json:"status_code"`
	Users      interface{} `json:"users, omitempty"`
}

type ResponseUserDetail struct {
	Message    string `json:"message"`
	Success    bool   `json:"success"`
	StatusCode int    `json:"status_code"`
	ResponseUsers
}

type ResponseUserUpdate struct {
	Message    string `json:"message"`
	Success    bool   `json:"success"`
	StatusCode int    `json:"status_code"`
	*UpdateUsers
}

type ResponseUserGet struct {
	Message    string `json:"message"`
	Success    bool   `json:"success"`
	StatusCode int    `json:"status_code"`
	ResponseUsers
}

type ResponseAuth struct {
	Authentication interface{} `json:"data", omitempty"`
	Message        string      `json:"message"`
}

type ResponseBookmark struct {
	Bookmarks interface{} `json:"bookmarks, omitempty"`
	Message   string      `json:"message"`
}

type ResponseAddCategories struct {
	*UsersCategories
	Message string `json:"message"`
}

type ResponsePost struct {
	Posts      interface{} `json:"posts, omitempty"`
	Message    string      `json:"message"`
	Success    bool        `json:"success"`
	StatusCode int         `json:"status_code"`
}

type ResponsePostUpdate struct {
	*Posts
	Message    string `json:"message"`
	Success    bool   `json:"success"`
	StatusCode int    `json:"status_code"`
}

type ResponseDetailPost struct {
	*PostsUsersJoin
	EditMode   bool   `json:"edit_mode"`
	Message    string `json:"message"`
	Success    bool   `json:"success"`
	StatusCode int    `json:"status_code"`
}

type ResponseCategory struct {
	Categories interface{} `json:"categories, omitempty"`
	Message    string      `json:"message"`
	Success    bool        `json:"success"`
	StatusCode int         `json:"status_code"`
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

	// Kode yang dikomen untuk delete tabel
	// db.DropTable("users", "bookmarks", "posts", "categories")
	// db.AutoMigrate(&m.Users{}, &m.Bookmarks{}, &m.Posts{}, &m.Categories{})
	db.AutoMigrate(&Users{}, &Bookmarks{}, &Posts{}, &Categories{}, &UsersCategories{})

}
