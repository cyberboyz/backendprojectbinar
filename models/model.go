package models

import (
	"time"
)

type SuccessStatus struct {
	Message    string `json:"message"`
	Success    bool   `json:"success"`
	StatusCode int    `json:"status_code"`
}

type Users struct {
	ResponseUsers
	Password string `form:"password" json:"password"`
}

type ResponseUsersSignUp struct {
	ID    uint   `gorm:"primary_key" json:"id_user"`
	Email string `gorm:"unique_index" form:"email" json:"email"`
	Name  string `form:"name" json:"name"`
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

type NewResponseUserDetail struct {
	ResponseUsers
	Message        string   `json:"message"`
	ListCategories []string `json:"categories"`
	Success        bool     `json:"success"`
	StatusCode     int      `json:"status_code"`
}

type UpdateUsers struct {
	Email        string `json:"email" form:"email"`
	Name         string `json:"name" form:"name"`
	Address      string `json:"address" form:"address"`
	Bio          string `json:"bio" form:"bio"`
	IDAvatar     uint   `json:"id_avatar" form:"id_avatar"`
	IDCoverPhoto uint   `json:"id_cover_photo" form:"id_cover_photo"`
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
	Categories   string `gorm:"primary_key" json:"categories" form:"categories"`
	IDBackground uint   `json:"id_background" form:"id_background"`
}

type UsersCategories struct {
	ID         uint   `gorm:"primary_key" json:"id"`
	IDUser     uint   `json:"id_user, omitempty"`
	Categories string `json:"categories"`
}

type InputUsersCategories struct {
	ID         uint     `gorm:"primary_key" json:"id"`
	IDUser     uint     `form:"id_user" json:"id_user, omitempty"`
	Categories []string `form:"categories" json:"categories"`
}

type Posts struct {
	ID           uint      `gorm:"primary_key" form:"id_post" json:"id_post"`
	IDUser       uint      `form:"id_user" json:"id_user, omitempty"`
	IDBackground uint      `form:"id_background" json:"id_background, omitempty"`
	PostTitle    string    `form:"post_title" json:"post_title, omitempty"`
	Categories   string    `form:"categories" json:"categories, omitempty"`
	Content      string    `form:"content" json:"content, omitempty"`
	CreatedAt    time.Time `json:"published_at, omitempty"`
	UpdatedAt    time.Time `json:"updated_at, omitempty"`
}

type PostsUsersJoin struct {
	Posts
	UpdateUsers
}

type Bookmarks struct {
	ID        uint      `gorm:"primary_key" json:"id_bookmark"`
	IDUser    uint      `json:"id_user" form:"id_user"`
	IDPost    uint      `json:"id_post" form:"id_post"`
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
	*InputUsersCategories
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
