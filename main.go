package main

import (
	m "baru-dreamcatcher/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
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
	Email string `json:"email"`
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
	Email    string `json:"email"`
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

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db_url := os.Getenv("DATABASE_URL")
	if db_url == "" {
		db_url = "host=localhost user=postgres dbname=gorm sslmode=disable password=postgres"
	}

	db, err = gorm.Open("postgres", db_url)

	db.SingularTable(true)

	if err != nil {
		fmt.Println(err)
	}

	db.DropTable("users", "bookmarks", "posts", "categories")
	db.AutoMigrate(&m.Users{}, &m.Bookmarks{}, &m.Posts{}, &m.Categories{})

	router := gin.New()

	router.LoadHTMLGlob("htmlfile/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	v1 := router.Group("/v1")
	{

		v1.POST("/login/", LoginUser)
		v1.POST("/register/", RegisterUser)

		logged_in := v1.Group("")
		logged_in.Use(AuthorizeMiddleware)
		{
			logged_in.GET("/logout/", LogoutUser)
			post := logged_in.Group("/posts")
			{
				post.GET("/", PostGet)
				post.GET("/:id", PostDetail)
				post.POST("/", PostCreate)
				post.PUT("/:id", PostUpdate)
				post.PATCH("/:id", PostUpdate)
				post.DELETE("/:id", PostDelete)
			}
			usergroup := logged_in.Group("/profile")
			{
				usergroup.GET("/", UserGet)
				usergroup.GET("/:id", UserDetail)
				usergroup.POST("/", UserCreate)
				usergroup.PUT("/:id", UserUpdate)
				usergroup.PATCH("/:id", UserUpdate)
				usergroup.DELETE("/:id", UserDelete)
				// usergroup.GET("/posts", UserPostsGet)
			}
			bookmark := logged_in.Group("/bookmarks")
			{
				bookmark.POST("/", BookmarkCreate)
				bookmark.GET("/", BookmarkGet)
				bookmark.DELETE("/:id", BookmarkDelete)
			}
			category := logged_in.Group("/categories")
			{
				category.POST("/", CategoryCreate)
				category.GET("/", CategoryGet)
				category.GET("/:id", CategoryPostsList)
			}
		}
	}
	router.Run(":" + port)
}

func AuthorizeMiddleware(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	auth := &Users{}

	if authorization == "" {
		response := &m.ResponseUser{
			Message: "Cannot access the resource : You need to authenticate",
		}
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}

	err = db.Where("token = ? ", authorization).Find(&auth).Error

	if err != nil {
		response := &m.ResponseUser{
			Message: "Unauthorized access",
		}
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}
	c.Next()
}

func ValidateFormatEmail(email string) string {
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegexp.MatchString(email) {
		return "Error : Email is not correct"
	}
	return ""
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

	if register.Email == "" {
		response := &ResponseAuth{
			Message: "Error : Email cannot be empty",
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	validation := ValidateFormatEmail(register.Email)

	if validation != "" {
		response := &ResponseAuth{
			Message: validation,
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	register.Email = strings.ToLower(register.Email)

	if register.Password == "" {
		response := &ResponseAuth{
			Message: "Error : Password cannot be empty",
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	register.Password = string(hash)

	if register.Name == "" {
		response := &ResponseAuth{
			Message: "Error : Name cannot be empty",
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

	registerOutput := register.ResponseUsersSignUp

	response := &Response{
		Success:    true,
		StatusCode: http.StatusCreated,
		Message:    "New User has been created",
		Data:       registerOutput,
	}

	c.JSON(http.StatusCreated, response)
}

func randToken(length int) string {
	rand.Seed(time.Now().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func LoginUser(c *gin.Context) {
	login := &Users{}

	err = c.BindJSON(&login)
	inputPassword := login.Password
	login.Email = strings.ToLower(login.Email)

	db.Where("email = ? ", login.Email).Find(&login)

	err = bcrypt.CompareHashAndPassword([]byte(login.Password), []byte(inputPassword))

	if err != nil {
		response := &Response{
			Message:    "Error : Unauthorized User",
			Success:    false,
			StatusCode: http.StatusUnauthorized,
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	if login.Token == "" {
		login.Token = randToken(20)
		db.Model(login).Update("token", login.Token)
	}

	loginOutput := login.ResponseUsers

	response := &Response{
		Message:    "New User has been created",
		Success:    true,
		StatusCode: http.StatusOK,
		Data:       loginOutput,
	}

	c.JSON(http.StatusOK, response)
}

func LogoutUser(c *gin.Context) {
	logout := &m.Users{}

	authorization := c.Request.Header.Get("Authorization")

	db.Model(logout).Update("token", logout.Token).Where("token", authorization)

	response := &m.Response{
		Message:    "Logout successful",
		Success:    true,
		StatusCode: http.StatusOK,
	}

	c.JSON(http.StatusOK, response)
}

func UserGet(c *gin.Context) {

	users := []*Users{}
	err = db.Find(&users).Error

	if err != nil {
		response := &m.ResponseUser{
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	response := &m.ResponseUser{
		Message: "Get users",
		Users:   users,
	}

	c.JSON(http.StatusOK, response)
}

func UserDetail(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := &m.ResponseUser{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	user := &Users{}
	err = db.Where("id = ?", id).First(&user).Error

	if err != nil {
		response := &m.ResponseUser{
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	response := &m.ResponseUser{
		Message: "Get user",
		Users:   user,
	}

	c.JSON(http.StatusOK, response)
}

func UserCreate(c *gin.Context) {

	user := &Users{}
	err := c.BindJSON(&user)
	if err != nil {
		response := &m.ResponseUser{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	err = db.Create(user).Error

	if err != nil {
		response := &m.ResponseUser{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	response := &m.ResponseUser{
		Message: "User has been created",
		Users:   user,
	}

	c.JSON(http.StatusCreated, response)
}

func UserUpdate(c *gin.Context) {

	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	id := uint(id64)
	if err != nil {
		response := &m.ResponseUser{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	user := &Users{}
	err = c.BindJSON(&user)
	if err != nil {
		response := &m.ResponseUser{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	user.ID = id
	err = db.Save(user).Error

	if err != nil {
		response := &m.ResponseUser{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	response := &m.ResponseUser{
		Message: "User has been updated",
		Users:   user,
	}

	c.JSON(http.StatusOK, response)
}

func UserDelete(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := &m.ResponseUser{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	err = db.Where("id = ?", id).Delete(&Users{}).Error

	if err != nil {
		response := &m.ResponseUser{
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	response := &m.ResponseUser{
		Message: "User has been deleted",
	}

	c.JSON(http.StatusOK, response)
}

func BookmarkCreate(c *gin.Context) {

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

	posts := []*Posts{}
	err = db.Order("created_at desc").Find(&posts).Error
	// err = db.Table("posts").Order("created_at desc").Joins("JOIN users on users.id = posts.id_user").Scan(&posts).Error

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

	authorization := c.Request.Header.Get("Authorization")
	auth := &Users{}
	err = db.Where("token = ? ", authorization).Find(&auth).Error

	post.IDUser = auth.ID

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
// 		response := &m.ResponseUser{
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
