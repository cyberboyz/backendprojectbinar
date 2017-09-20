package controllers

import (
	m "baru-dreamcatcher/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var err error
var db *gorm.DB

func AuthorizeMiddleware(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	auth := &m.Users{}

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
	register := &m.Users{}

	err := c.BindJSON(&register)
	if err != nil {
		response := &m.ResponseAuth{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	if register.Email == "" {
		response := &m.ResponseAuth{
			Message: "Error : Email cannot be empty",
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	validation := ValidateFormatEmail(register.Email)

	if validation != "" {
		response := &m.ResponseAuth{
			Message: validation,
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	register.Email = strings.ToLower(register.Email)

	if register.Password == "" {
		response := &m.ResponseAuth{
			Message: "Error : Password cannot be empty",
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	register.Password = string(hash)

	if register.Name == "" {
		response := &m.ResponseAuth{
			Message: "Error : Name cannot be empty",
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	err = db.Create(register).Error

	if err != nil {
		response := &m.ResponseAuth{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	registerOutput := register.ResponseUsersSignUp

	response := &m.Response{
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
	login := &m.Users{}

	err = c.BindJSON(&login)
	inputPassword := login.Password
	login.Email = strings.ToLower(login.Email)

	db.Where("email = ? ", login.Email).Find(&login)

	err = bcrypt.CompareHashAndPassword([]byte(login.Password), []byte(inputPassword))

	if err != nil {
		response := &m.Response{
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

	response := &m.Response{
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

	users := []*m.Users{}
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

	user := &m.Users{}
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

	user := &m.Users{}
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

	user := &m.Users{}
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

	err = db.Where("id = ?", id).Delete(&m.Users{}).Error

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

	bookmarks := &m.Bookmarks{}
	err := c.BindJSON(&bookmarks)
	if err != nil {
		response := &m.ResponseBookmark{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	err = db.Create(bookmarks).Error

	if err != nil {
		response := &m.ResponseBookmark{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	response := &m.ResponseBookmark{
		Message:   "Bookmark has been created",
		Bookmarks: bookmarks,
	}

	c.JSON(http.StatusCreated, response)
}

func BookmarkGet(c *gin.Context) {

	bookmarks := []*m.Bookmarks{}
	err = db.Find(&bookmarks).Error

	if err != nil {
		response := &m.ResponseBookmark{
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	response := &m.ResponseBookmark{
		Message:   "Get bookmarks",
		Bookmarks: bookmarks,
	}

	c.JSON(http.StatusOK, response)
}

func BookmarkDelete(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := &m.ResponseBookmark{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	err = db.Where("id = ?", id).Delete(&m.Bookmarks{}).Error

	if err != nil {
		response := &m.ResponseBookmark{
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	response := &m.ResponseBookmark{
		Message: "Bookmark has been deleted",
	}

	c.JSON(http.StatusOK, response)
}

func PostGet(c *gin.Context) {

	posts := []*m.Posts{}
	err = db.Order("created_at desc").Find(&posts).Error
	// err = db.Table("posts").Order("created_at desc").Joins("JOIN users on users.id = posts.id_user").Scan(&posts).Error

	if err != nil {
		response := &m.ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	response := &m.ResponsePost{
		posts, "Get posts", m.SuccessStatus{true, 200},
	}

	c.JSON(http.StatusOK, response)
}

func PostDetail(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := &m.ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	post := &m.Posts{}
	err = db.Where("id = ?", id).First(&post).Error

	if err != nil {
		response := &m.ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	response := &m.ResponsePost{
		post,
		"Post has been obtained",
		m.SuccessStatus{true, 200},
	}

	c.JSON(http.StatusOK, response)
}

func PostCreate(c *gin.Context) {

	post := &m.Posts{}
	err := c.BindJSON(&post)
	if err != nil {
		response := &m.ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	authorization := c.Request.Header.Get("Authorization")
	auth := &m.Users{}
	err = db.Where("token = ? ", authorization).Find(&auth).Error

	post.IDUser = auth.ID

	err = db.Create(post).Error

	if err != nil {
		response := &m.ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	response := &m.ResponsePost{
		post,
		"Post has been created",
		m.SuccessStatus{true, 200},
	}

	c.JSON(http.StatusCreated, response)
}

func PostUpdate(c *gin.Context) {

	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	id := uint(id64)
	if err != nil {
		response := &m.ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	post := &m.Posts{}
	err = c.BindJSON(&post)
	if err != nil {
		response := &m.ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	post.ID = id
	err = db.Save(post).Error

	if err != nil {
		response := &m.ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	response := &m.ResponsePost{
		post,
		"Post has been updated",
		m.SuccessStatus{true, 200},
	}

	c.JSON(http.StatusOK, response)
}

func PostDelete(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := &m.ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	err = db.Where("id = ?", id).Delete(&m.Posts{}).Error

	if err != nil {
		response := &m.ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	response := &m.ResponsePost{
		Message: "Post has been deleted",
	}

	c.JSON(http.StatusOK, response)
}

func CategoryCreate(c *gin.Context) {

	category := &m.Categories{}
	err := c.BindJSON(&category)
	if err != nil {
		response := &m.ResponseCategory{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	err = db.Create(category).Error

	if err != nil {
		response := &m.ResponseCategory{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	response := &m.ResponseCategory{
		Message:    "Post has been created",
		Categories: category,
	}

	c.JSON(http.StatusCreated, response)
}

func CategoryGet(c *gin.Context) {

	categories := []*m.Categories{}
	err = db.Find(&categories).Error

	if err != nil {
		response := &m.ResponseCategory{
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	response := &m.ResponseCategory{
		Message:    "Get categories",
		Categories: categories,
	}

	c.JSON(http.StatusOK, response)
}

func CategoryPostsList(c *gin.Context) {

	id := c.Param("id")

	posts := []*m.Posts{}
	err = db.Where("categories = ?", id).Find(&posts).Error

	if err != nil {
		response := &m.ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	response := &m.ResponsePost{
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
// 		response := &m.ResponsePost{
// 			Message: err.Error(),
// 		}
// 		c.JSON(http.StatusServiceUnavailable, response)
// 		c.Abort()
// 		return
// 	}

// 	response := &m.ResponsePost{
// 		posts, "Get posts", m.SuccessStatus{true, 200},
// 	}

// 	c.JSON(http.StatusOK, response)
// }
