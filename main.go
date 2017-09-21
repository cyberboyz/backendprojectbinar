package main

import (
	// cont "baru-dreamcatcher/controllers"
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

	// Kode yang dikomen untuk delete tabel
	// db.DropTable("users", "bookmarks", "posts", "categories")
	db.AutoMigrate(&m.Users{}, &m.Bookmarks{}, &m.Posts{}, &m.Categories{})

	router := gin.New()

	// Menampilkan page di halaman utama terkait resource yang bisa diakses beserta method-nya
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
				post.DELETE("/:id", PostDelete)
			}
			usergroup := logged_in.Group("/profile")
			{
				usergroup.GET("/", UserGet)
				usergroup.GET("/:id", UserDetail)
				usergroup.PUT("/:id", UserUpdate)
				usergroup.DELETE("/:id", UserDelete)
				usergroup.GET("/:id/posts", UserPostsGet)
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
		Message:    "Registration successful : New User has been created",
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

	response := &m.ResponseLogin{
		Token:      login.Token,
		Email:      login.Email,
		Message:    "New User has been created",
		Success:    true,
		StatusCode: http.StatusOK,
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
		Message:    "Get users : All Users have been shown",
		Success:    true,
		StatusCode: http.StatusOK,
		Users:      users,
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

	output := user.ResponseUsers

	response := &m.ResponseUser{
		Message:    "Get user : Certain user detail has been shown",
		Success:    true,
		StatusCode: http.StatusOK,
		Users:      output,
	}

	c.JSON(http.StatusOK, response)
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

	user := &m.UpdateUsers{}
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
	err = db.Table("users").Update(user).Error

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

	posts := []*m.PostsUsersJoin{}
	err = db.Raw("SELECT * FROM posts join users on posts.id_user = users.id ORDER BY posts.created_at desc").Scan(&posts).Error
	// fmt.Println(posts.ResponseUsers.IDAvatar)

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

	authorization := c.Request.Header.Get("Authorization")
	auth := &m.Users{}

	err = db.Where("token = ? ", authorization).Find(&auth).Error

	if err != nil {
		response := &m.ResponsePost{
			Message: "Cannot execute the query to search user based on token",
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	post := &m.PostsUsersJoin{}
	err = db.Raw("SELECT * FROM posts join users on posts.id_user = users.id WHERE posts.id = ?", id).Scan(&post).Error

	if err != nil {
		response := &m.ResponsePost{
			Message: "Cannot execute the join query from posts and users table",
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	editMode := false
	if auth.ID == post.IDUser {
		editMode = true
	}

	response := &m.ResponseDetailPost{
		post,
		"Post has been obtained",
		editMode,
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

	authorization := c.Request.Header.Get("Authorization")
	auth := &m.Users{}
	checkuserID := &m.Posts{}

	err = db.Where("token = ? ", authorization).Find(&auth).Error
	err = db.Where("id = ?", id).Find(&checkuserID).Error

	if checkuserID.IDUser != auth.ID {
		response := &m.ResponsePost{
			Message: "Error : Cannot change other people posts",
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	post.IDUser = auth.ID
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
		Message:    "Category has been created",
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
		Message:    "Categories have been shown",
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

func UserPostsGet(c *gin.Context) {

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

	posts := []*m.Posts{}
	err = db.Order("created_at desc").Where("id_user = ?", id).Find(&posts).Error
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
