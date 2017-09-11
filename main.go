package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
	"time"
)

type Users struct {
	ID           uint      `gorm:"primary_key" json:"id_user"`
	Email        string    `gorm:"primary_key" json:"email" binding:"required"`
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	Bio          string    `json:"bio"`
	IDAvatar     int       `json:"id_avatar"`
	IDCoverPhoto int       `json:"id_cover_photo"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Password     string    `json:"password" binding:"required"`
}

type Categories struct {
	Categories   string `gorm:"primary_key" json:"categories"`
	IDBackground uint   `json:"id_background"`
}

type Posts struct {
	ID           uint      `gorm:"primary_key" json:"id_post"`
	IDUser       uint      `json:"id_user"`
	IDBackground uint      `json:"id_background"`
	PostTitle    string    `json:"post_title"`
	Categories   string    `json:"categories"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
}

type Bookmarks struct {
	ID        uint      `gorm:"primary_key" json:"id_bookmark"`
	IDUser    uint      `json:"id_user"`
	IDPost    uint      `json:"id_post"`
	CreatedAt time.Time `json:"created_at"`
}

type Authentication struct {
	Email    string `gorm:"primary_key" json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SuccessStatus struct {
	Success    bool `json:"success"`
	StatusCode int  `json:"status_code"`
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

func main() {
	// db, err = gorm.Open("postgres", "host=localhost user=postgres dbname=gorm sslmode=disable password=postgres")
	db, err = gorm.Open("postgres", "postgres://ssvuoibpdkugsp:c73e9d4dfbc63b197d4c2336c77adcc4d974afccf584f8a5e793d1e7ac90d242@ec2-50-17-236-15.compute-1.amazonaws.com:5432/dso2rra5vg6aa")

	db.SingularTable(true)
	if err != nil {
		fmt.Println(err)
	}
	db.AutoMigrate(&Users{}, &Bookmarks{}, &Posts{}, &Categories{})

	router := gin.Default()

	v1 := router.Group("/v1")
	{
		post := v1.Group("/posts")
		{
			post.GET("/", PostGet)
			post.GET("/:id", PostDetail)
			post.POST("/", PostCreate)
			post.PUT("/:id", PostUpdate)
			post.PATCH("/:id", PostUpdate)
			post.DELETE("/:id", PostDelete)
		}
		usergroup := v1.Group("/profile")
		{
			usergroup.GET("/", UserGet)
			usergroup.GET("/:id", UserDetail)
			usergroup.POST("/", UserCreate)
			usergroup.PUT("/:id", UserUpdate)
			usergroup.PATCH("/:id", UserUpdate)
			usergroup.DELETE("/:id", UserDelete)
			// usergroup.GET("/posts", UserPostsGet)
		}
		bookmark := v1.Group("/bookmarks")
		{
			bookmark.POST("/", BookmarkCreate)
			bookmark.GET("/", BookmarkGet)
			bookmark.DELETE("/:id", BookmarkDelete)
		}
		category := v1.Group("/categories")
		{
			category.POST("/", CategoryCreate)
			category.GET("/", CategoryGet)
			category.GET("/:id", CategoryPostsList)
		}
		v1.POST("/login/", LoginUser)
		v1.POST("/register/", RegisterUser)
		v1.POST("/logout/", LogoutUser)
	}
	router.Run(":8080")
}

func RegisterUser(c *gin.Context) {
	register := &Users{}

	err := c.BindJSON(&register)
	if err != nil {
		response := &ResponseAuth{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	err = db.Create(register).Error

	if err != nil {
		response := &ResponseAuth{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	response := &ResponseAuth{
		Message:        "New User has been created",
		Authentication: register,
	}

	c.JSON(http.StatusCreated, response)
}

func LoginUser(c *gin.Context) {
	login := &Users{}
	if c.BindJSON(login) == nil {
		err = db.Select("password").Where("email = ? AND password = ?", login.Email, login.Password).Find(&login).Error
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": "2n41jt01fk1-cj2190c129je211x910s19k112i012d",
				"email": login.Email, "success": true, "status_code": http.StatusOK})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized", "success": false, "status_code": http.StatusUnauthorized})
		}
	}
}

func LogoutUser(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	if authorization != "2n41jt01fk1-cj2190c129je211x910s19k112i012d" {
		response := &ResponseUser{
			Message: "Unauthorized access",
		}
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful",
		"success": true, "status_code": http.StatusOK})
}

func Authorize(c *gin.Context) {
	// fmt.Println(c)
	authorization := c.Request.Header.Get("Authorization")
	if authorization != "2n41jt01fk1-cj2190c129je211x910s19k112i012d" {
		response := &ResponseUser{
			Message: "Unauthorized access",
		}
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}
}

func UserGet(c *gin.Context) {

	authorization := c.Request.Header.Get("Authorization")
	if authorization != "2n41jt01fk1-cj2190c129je211x910s19k112i012d" {
		response := &ResponseUser{
			Message: "Unauthorized access",
		}
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}

	users := []*Users{}
	err = db.Find(&users).Error

	if err != nil {
		response := &ResponseUser{
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	response := &ResponseUser{
		Message: "Get users",
		Users:   users,
	}

	c.JSON(http.StatusOK, response)
}

func UserDetail(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	if authorization != "2n41jt01fk1-cj2190c129je211x910s19k112i012d" {
		response := &ResponseUser{
			Message: "Unauthorized access",
		}
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := &ResponseUser{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	user := &Users{}
	err = db.Where("id = ?", id).First(&user).Error

	if err != nil {
		response := &ResponseUser{
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	response := &ResponseUser{
		Message: "Get user",
		Users:   user,
	}

	c.JSON(http.StatusOK, response)
}

func UserCreate(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	if authorization != "2n41jt01fk1-cj2190c129je211x910s19k112i012d" {
		response := &ResponseUser{
			Message: "Unauthorized access",
		}
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}

	user := &Users{}
	err := c.BindJSON(&user)
	if err != nil {
		response := &ResponseUser{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	err = db.Create(user).Error

	if err != nil {
		response := &ResponseUser{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	response := &ResponseUser{
		Message: "User has been created",
		Users:   user,
	}

	c.JSON(http.StatusCreated, response)
}

func UserUpdate(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	if authorization != "2n41jt01fk1-cj2190c129je211x910s19k112i012d" {
		response := &ResponseUser{
			Message: "Unauthorized access",
		}
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}

	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	id := uint(id64)
	if err != nil {
		response := &ResponseUser{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	user := &Users{}
	err = c.BindJSON(&user)
	if err != nil {
		response := &ResponseUser{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	user.ID = id
	err = db.Save(user).Error

	if err != nil {
		response := &ResponseUser{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	response := &ResponseUser{
		Message: "User has been updated",
		Users:   user,
	}

	c.JSON(http.StatusOK, response)
}

func UserDelete(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	if authorization != "2n41jt01fk1-cj2190c129je211x910s19k112i012d" {
		response := &ResponseUser{
			Message: "Unauthorized access",
		}
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := &ResponseUser{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	err = db.Where("id = ?", id).Delete(&Users{}).Error

	if err != nil {
		response := &ResponseUser{
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	response := &ResponseUser{
		Message: "User has been deleted",
	}

	c.JSON(http.StatusOK, response)
}

func BookmarkCreate(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	if authorization != "2n41jt01fk1-cj2190c129je211x910s19k112i012d" {
		response := &ResponseUser{
			Message: "Unauthorized access",
		}
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}

	bookmarks := &Bookmarks{}
	err := c.BindJSON(&bookmarks)
	if err != nil {
		response := &ResponseBookmark{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	err = db.Create(bookmarks).Error

	if err != nil {
		response := &ResponseBookmark{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	response := &ResponseBookmark{
		Message:   "Bookmark has been created",
		Bookmarks: bookmarks,
	}

	c.JSON(http.StatusCreated, response)
}

func BookmarkGet(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	if authorization != "2n41jt01fk1-cj2190c129je211x910s19k112i012d" {
		response := &ResponseUser{
			Message: "Unauthorized access",
		}
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}

	bookmarks := []*Bookmarks{}
	err = db.Find(&bookmarks).Error

	if err != nil {
		response := &ResponseBookmark{
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	response := &ResponseBookmark{
		Message:   "Get bookmarks",
		Bookmarks: bookmarks,
	}

	c.JSON(http.StatusOK, response)
}

func BookmarkDelete(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	if authorization != "2n41jt01fk1-cj2190c129je211x910s19k112i012d" {
		response := &ResponseUser{
			Message: "Unauthorized access",
		}
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := &ResponseBookmark{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	err = db.Where("id = ?", id).Delete(&Bookmarks{}).Error

	if err != nil {
		response := &ResponseBookmark{
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	response := &ResponseBookmark{
		Message: "Bookmark has been deleted",
	}

	c.JSON(http.StatusOK, response)
}

func PostGet(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	if authorization != "2n41jt01fk1-cj2190c129je211x910s19k112i012d" {
		response := &ResponseUser{
			Message: "Unauthorized access",
		}
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}

	posts := []*Posts{}
	err = db.Find(&posts).Error

	if err != nil {
		response := &ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	response := &ResponsePost{
		posts, "Get posts", SuccessStatus{true, 200},
	}

	c.JSON(http.StatusOK, response)
}

func PostDetail(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	if authorization != "2n41jt01fk1-cj2190c129je211x910s19k112i012d" {
		response := &ResponseUser{
			Message: "Unauthorized access",
		}
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := &ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	post := &Posts{}
	err = db.Where("id = ?", id).First(&post).Error

	if err != nil {
		response := &ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	response := &ResponsePost{
		post,
		"Post has been obtained",
		SuccessStatus{true, 200},
	}

	c.JSON(http.StatusOK, response)
}

func PostCreate(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	if authorization != "2n41jt01fk1-cj2190c129je211x910s19k112i012d" {
		response := &ResponseUser{
			Message: "Unauthorized access",
		}
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}

	post := &Posts{}
	err := c.BindJSON(&post)
	if err != nil {
		response := &ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	err = db.Create(post).Error

	if err != nil {
		response := &ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	response := &ResponsePost{
		post,
		"Post has been created",
		SuccessStatus{true, 200},
	}

	c.JSON(http.StatusCreated, response)
}

func PostUpdate(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	if authorization != "2n41jt01fk1-cj2190c129je211x910s19k112i012d" {
		response := &ResponseUser{
			Message: "Unauthorized access",
		}
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}

	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	id := uint(id64)
	if err != nil {
		response := &ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	post := &Posts{}
	err = c.BindJSON(&post)
	if err != nil {
		response := &ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	post.ID = id
	err = db.Save(post).Error

	if err != nil {
		response := &ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	response := &ResponsePost{
		post,
		"Post has been updated",
		SuccessStatus{true, 200},
	}

	c.JSON(http.StatusOK, response)
}

func PostDelete(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	if authorization != "2n41jt01fk1-cj2190c129je211x910s19k112i012d" {
		response := &ResponseUser{
			Message: "Unauthorized access",
		}
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := &ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	err = db.Where("id = ?", id).Delete(&Posts{}).Error

	if err != nil {
		response := &ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	response := &ResponsePost{
		Message: "Post has been deleted",
	}

	c.JSON(http.StatusOK, response)
}

func CategoryCreate(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	if authorization != "2n41jt01fk1-cj2190c129je211x910s19k112i012d" {
		response := &ResponseUser{
			Message: "Unauthorized access",
		}
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}

	category := &Categories{}
	err := c.BindJSON(&category)
	if err != nil {
		response := &ResponseCategory{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	err = db.Create(category).Error

	if err != nil {
		response := &ResponseCategory{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	response := &ResponseCategory{
		Message:    "Post has been created",
		Categories: category,
	}

	c.JSON(http.StatusCreated, response)
}

func CategoryGet(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	if authorization != "2n41jt01fk1-cj2190c129je211x910s19k112i012d" {
		response := &ResponseUser{
			Message: "Unauthorized access",
		}
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}

	categories := []*Categories{}
	err = db.Find(&categories).Error

	if err != nil {
		response := &ResponseCategory{
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	response := &ResponseCategory{
		Message:    "Get categories",
		Categories: categories,
	}

	c.JSON(http.StatusOK, response)
}

func CategoryPostsList(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	if authorization != "2n41jt01fk1-cj2190c129je211x910s19k112i012d" {
		response := &ResponseUser{
			Message: "Unauthorized access",
		}
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}

	id := c.Param("id")

	posts := []*Posts{}
	err = db.Where("categories = ?", id).Find(&posts).Error

	if err != nil {
		response := &ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	response := &ResponsePost{
		Message: "Get posts by category",
		Posts:   posts,
	}

	c.JSON(http.StatusOK, response)
}

// func UserPostsGet(c *gin.Context) {
// 	authorization := c.Request.Header.Get("Authorization")
// 	if authorization != "2n41jt01fk1-cj2190c129je211x910s19k112i012d" {
// 		response := &ResponseUser{
// 			Message: "Unauthorized access",
// 		}
// 		c.JSON(http.StatusUnauthorized, response)
// 		c.Abort()
// 		return
// 	}

// 	posts := []*Posts{}

// 	err = db.Where("id = ?", posts.).Find(&posts).Error

// 	if err != nil {
// 		response := &ResponsePost{
// 			Message: err.Error(),
// 		}
// 		c.JSON(http.StatusServiceUnavailable, response)
// 		c.Abort()
// 		return
// 	}

// 	response := &ResponsePost{
// 		posts, "Get posts", SuccessStatus{true, 200},
// 	}

// 	c.JSON(http.StatusOK, response)
// }
