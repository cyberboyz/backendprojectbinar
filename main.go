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
	db.AutoMigrate(&m.Users{}, &m.Bookmarks{}, &m.Posts{}, &m.Categories{}, &m.UsersCategories{})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

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
		{
			logged_in.GET("/logout/", LogoutUser)
			logged_in.GET("/ownprofile", ShowOwnProfile)
			logged_in.POST("/3categoriesposts", SeveralCategoriesPostsList)
			post := logged_in.Group("/posts")
			{
				// comment := post.Group("/comments")
				// {
				// 	// comment.GET("/", CommentGet)
				// 	// comment.GET("/:id", CommentDetail)
				// 	// comment.POST("/", CommentCreate).Use(AuthorizeMiddleware)
				// 	// comment.PUT("/:id", CommentUpdate).Use(AuthorizeMiddleware)
				// 	// comment.DELETE("/:id", CommentDelete).Use(AuthorizeMiddleware)
				// }
				post.GET("/", PostGet)
				post.GET("/:id", PostDetail)
				post.POST("/", PostCreate).Use(AuthorizeMiddleware)
				post.PUT("/:id", PostUpdate).Use(AuthorizeMiddleware)
				post.DELETE("/:id", PostDelete).Use(AuthorizeMiddleware)
			}
			usergroup := logged_in.Group("/profile")
			{
				usergroup.GET("/", UserGet)
				usergroup.GET("/:id", UserDetail)
				usergroup.PUT("/:id", UserUpdate).Use(AuthorizeMiddleware)
				usergroup.DELETE("/:id", UserDelete).Use(AuthorizeMiddleware)
				usergroup.GET("/:id/posts", UserPostsGet).Use(AuthorizeMiddleware)
				usergroup.POST("/categories", AddCategoriesByUser).Use(AuthorizeMiddleware)
			}
			bookmark := logged_in.Group("/bookmarks")
			{
				bookmark.POST("/", BookmarkCreate).Use(AuthorizeMiddleware)
				bookmark.GET("/", BookmarkGet).Use(AuthorizeMiddleware)
				bookmark.DELETE("/:id", BookmarkDelete).Use(AuthorizeMiddleware)
			}
			category := logged_in.Group("/categories")
			{
				category.GET("/:id", CategoryPostsList)
				category.POST("/", CategoryCreate).Use(AuthorizeMiddleware)
				category.GET("/", CategoryGet).Use(AuthorizeMiddleware)
			}
			logged_in.GET("/ownposts", ShowOwnPosts).Use(AuthorizeMiddleware)
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

	err := c.Bind(register)
	if err != nil {
		response := &m.ResponseAuth{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	// fmt.Println(register.Email)
	// fmt.Println(register.Password)
	// fmt.Println(register.Name)

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

	register.Token = randToken(20)

	err = db.Create(register).Error

	if err != nil {
		response := &m.Response{
			Success:    false,
			StatusCode: http.StatusCreated,
			Message:    "Error : Email has already been registered",
		}
		c.JSON(http.StatusOK, response)
		c.Abort()
		return
	}

	output := register.ResponseUsers
	output.TotalPosts = 0

	var listCategories []string

	response := &m.NewResponseUserDetail{
		Message:        "Get user : Certain user detail has been shown",
		Success:        true,
		StatusCode:     http.StatusOK,
		ResponseUsers:  output,
		ListCategories: listCategories,
	}

	c.JSON(http.StatusOK, response)
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

	c.Bind(login)
	// c.BindJSON(login)

	// fmt.Println(login.Password)

	inputPassword := login.Password
	login.Email = strings.ToLower(login.Email)
	// fmt.Println(login.Email)

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
	// fmt.Println(login.Email)

	var id_user uint
	db.Raw("SELECT id from users where email = ?", login.Email).Row().Scan(&id_user)
	// fmt.Println(id_user)

	var total_post uint
	db.Raw("SELECT count(id) FROM posts where id_user = ?", id_user).Row().Scan(&total_post)
	// fmt.Println(total_post)

	var userCategories string
	categoriesRows, _ := db.Raw("SELECT categories from users_categories where id_user = ? GROUP BY categories", id_user).Rows()
	defer categoriesRows.Close()

	var listCategories []string

	for categoriesRows.Next() {
		categoriesRows.Scan(&userCategories)
		// fmt.Println(userCategories)
		listCategories = append(listCategories, userCategories)
	}

	output := login.ResponseUsers
	output.TotalPosts = total_post
	// output.ListCategories = listCategories

	response := &m.NewResponseUserDetail{
		Message:        "Get user : Certain user detail has been shown",
		Success:        true,
		StatusCode:     http.StatusOK,
		ResponseUsers:  output,
		ListCategories: listCategories,
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

	outputUser := []m.ResponseUsers{}

	for _, element := range users {
		outputUser = append(outputUser, element.ResponseUsers)
	}

	response := &m.ResponseUser{
		Message:    "Get users : All Users have been shown",
		Success:    true,
		StatusCode: http.StatusOK,
		Users:      outputUser,
	}

	c.JSON(http.StatusOK, response)

	// categories := []*m.Categories{}
	// db.Raw("SELECT categories from users_categories where id_user = ?", element.ResponseUsers.ID).Scan(&categories)

	// if err != nil {
	// 	response := &m.ResponseCategory{
	// 		Message: err.Error(),
	// 	}
	// 	c.JSON(http.StatusServiceUnavailable, response)
	// 	c.Abort()
	// 	return
	// }

	// response := &m.ResponseCategory{
	// 	Message:    "Categories have been shown",
	// 	Success:    true,
	// 	StatusCode: http.StatusOK,
	// 	Categories: categories,
	// }

	// c.JSON(http.StatusOK, response)
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

	var total_post uint
	db.Raw("SELECT count(id) FROM posts where id_user = ?", id).Row().Scan(&total_post)

	if err != nil {
		response := &m.ResponseUser{
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	var userCategories string
	categoriesRows, err := db.Raw("SELECT categories from users_categories where id_user = ? GROUP BY categories", id).Rows()
	defer categoriesRows.Close()

	var listCategories []string

	for categoriesRows.Next() {
		categoriesRows.Scan(&userCategories)
		// fmt.Println(userCategories)
		listCategories = append(listCategories, userCategories)
	}

	output := user.ResponseUsers
	output.TotalPosts = total_post
	// output.ListCategories = listCategories

	response := &m.NewResponseUserDetail{
		Message:        "Get user : Certain user detail has been shown",
		Success:        true,
		StatusCode:     http.StatusOK,
		ResponseUsers:  output,
		ListCategories: listCategories,
	}

	c.JSON(http.StatusOK, response)
}

func UserUpdate(c *gin.Context) {

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

	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	id := uint(id64)
	if err != nil {
		response := &m.SuccessStatus{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	user := &m.UpdateUsers{}
	err = c.Bind(user)
	if err != nil {
		response := &m.SuccessStatus{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	err = db.Table("users").Where("id = ?", id).Update(user).Error

	if err != nil {
		response := &m.SuccessStatus{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	response := &m.ResponseUserUpdate{
		Message:     "User has been updated",
		Success:     true,
		StatusCode:  http.StatusOK,
		UpdateUsers: user,
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

func ShowOwnProfile(c *gin.Context) {

	authorization := c.Request.Header.Get("Authorization")
	auth := &m.Users{}
	err = db.Where("token = ? ", authorization).Find(&auth).Error

	user := &m.Users{}
	err = db.Where("id = ?", auth.ID).First(&user).Error

	var total_post uint
	db.Raw("SELECT count(id) FROM posts where id_user = ?", auth.ID).Row().Scan(&total_post)

	if err != nil {
		response := &m.ResponseUser{
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	var userCategories string
	categoriesRows, _ := db.Raw("SELECT categories from users_categories where id_user = ? GROUP BY categories", auth.ID).Rows()
	defer categoriesRows.Close()

	var listCategories []string

	for categoriesRows.Next() {
		categoriesRows.Scan(&userCategories)
		// fmt.Println(userCategories)
		listCategories = append(listCategories, userCategories)
	}

	output := user.ResponseUsers
	output.TotalPosts = total_post
	// output.ListCategories = listCategories

	response := &m.NewResponseUserDetail{
		Message:        "Get user : Certain user detail has been shown",
		Success:        true,
		StatusCode:     http.StatusOK,
		ResponseUsers:  output,
		ListCategories: listCategories,
	}

	c.JSON(http.StatusOK, response)

}

func BookmarkCreate(c *gin.Context) {

	bookmarks := &m.Bookmarks{}
	err := c.Bind(bookmarks)
	if err != nil {
		response := &m.ResponseBookmark{
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
		response := &m.ResponseBookmark{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	var post_title string
	db.Table("posts").Where("id = ?", bookmarks.IDPost).Select("post_title").Row().Scan(&post_title)

	// fmt.Println(bookmarks.IDPost)
	// fmt.Println(post_title)
	if post_title == "" {
		response := &m.ResponseBookmark{
			Message: "Error : You cannot bookmark this post, it does not exist",
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	bookmarks.IDUser = auth.ID

	var id_user string
	db.Table("bookmarks").Where("id_user = ? AND id_post = ? ", bookmarks.IDUser, bookmarks.IDPost).Select("id_user").Row().Scan(&id_user)

	if id_user != "" {
		response := &m.ResponseBookmark{
			Message: "Error : This user has bookmarked this post",
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

	authorization := c.Request.Header.Get("Authorization")
	auth := &m.Users{}
	err = db.Where("token = ? ", authorization).Find(&auth).Error

	if err != nil {
		response := &m.ResponseBookmark{
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	err = db.Where("id_user = ?", auth.ID).Find(&bookmarks).Error

	if err != nil {
		response := &m.ResponseBookmark{
			Message: "This user has not bookmarked any posts",
		}
		c.JSON(http.StatusOK, response)
		c.Abort()
		return
	}

	posts := m.PostsUsersJoin{}
	outputBookmark := []m.PostsUsersJoin{}

	for _, element := range bookmarks {
		db.Raw("SELECT * FROM posts join users on posts.id_user = users.id WHERE posts.id = ? ORDER BY posts.created_at desc", element.IDPost).Scan(&posts)
		// fmt.Println(element.IDPost)
		// fmt.Println(posts.PostTitle)
		outputBookmark = append(outputBookmark, posts)
	}

	// db.Raw("SELECT * FROM posts join users on posts.id_user = users.id WHERE posts.id = ? ORDER BY posts.created_at desc", element.IDPost).Scan(&outputBookmark)

	response := &m.ResponsePost{
		outputBookmark,
		"Get posts",
		true,
		200,
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

	authorization := c.Request.Header.Get("Authorization")
	auth := &m.Users{}
	db.Where("token = ? ", authorization).Find(&auth)

	checkDeleteBookmark := db.Where("id_post = ? AND id_user = ?", id, auth.ID).Delete(m.Bookmarks{}).RowsAffected
	// fmt.Println(id)
	// fmt.Println(auth.ID)

	if checkDeleteBookmark == 0 {
		response := &m.ResponseBookmark{
			Message: "Error : Either the bookmark you want to delete does not exist or you do not have enough authorization to delete other people bookmark",
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
		posts, "Get posts", true, 200,
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

	db.Where("token = ? ", authorization).Find(&auth)

	post := &m.PostsUsersJoin{}
	err = db.Raw("SELECT * FROM posts join users on posts.id_user = users.id WHERE posts.id = ?", id).Scan(&post).Error

	if err != nil {
		response := &m.ResponsePost{
			Message: "The post does not exist",
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
		editMode,
		"Post has been obtained",
		true,
		200,
	}

	c.JSON(http.StatusOK, response)
}

func PostCreate(c *gin.Context) {

	post := &m.Posts{}

	err := c.Bind(post)
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
		true,
		200,
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
	err = c.Bind(post)
	if err != nil {
		response := &m.ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	checkuserID := &m.Posts{}

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

	err = db.Where("id = ?", id).Find(&checkuserID).Error

	if checkuserID.IDUser == 0 {
		response := &m.ResponsePost{
			Message: "Error : The post does not exist",
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

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

	postusers := &m.PostsUsersJoin{}
	err = db.Raw("SELECT * FROM posts join users on posts.id_user = users.id WHERE posts.id = ?", id).Scan(&postusers).Error

	if err != nil {
		response := &m.ResponsePost{
			Message: "Cannot execute the join query from posts and users table",
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	editMode := false
	if auth.ID == postusers.IDUser {
		editMode = true
	}
	if err != nil {
		response := &m.ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	response := &m.ResponseDetailPost{
		postusers,
		editMode,
		"Post has been updated",
		true,
		200,
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
		Message:    "Post has been deleted",
		Success:    true,
		StatusCode: http.StatusOK,
	}

	c.JSON(http.StatusOK, response)
}

func CategoryCreate(c *gin.Context) {

	category := &m.Categories{}
	err := c.Bind(category)
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
		Success:    true,
		StatusCode: http.StatusOK,
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
		Success:    true,
		StatusCode: http.StatusOK,
		Categories: categories,
	}

	c.JSON(http.StatusOK, response)
}

func CategoryPostsList(c *gin.Context) {

	id := c.Param("id")

	posts := []*m.PostsUsersJoin{}
	err = db.Raw("SELECT * FROM posts join users on posts.id_user = users.id WHERE posts.categories = ? ORDER BY posts.created_at desc", id).Scan(&posts).Error

	if err != nil {
		response := &m.ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	response := &m.ResponsePost{
		posts, "Posts based on categories have been shown", true, 200,
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
		posts,
		"Get posts",
		true,
		200,
	}

	c.JSON(http.StatusOK, response)
}

func AddCategoriesByUser(c *gin.Context) {

	inputCategoriesByUser := &m.InputUsersCategories{}
	err := c.Bind(inputCategoriesByUser)

	if err != nil {
		response := &m.ResponseAddCategories{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	authorization := c.Request.Header.Get("Authorization")
	auth := &m.Users{}
	err = db.Where("token = ? ", authorization).Find(&auth).Error

	inputCategoriesByUser.IDUser = auth.ID
	// var categoryString string

	for _, element := range inputCategoriesByUser.Categories {
		categoriesByUser := &m.UsersCategories{}
		categoriesByUser.IDUser = auth.ID

		categoriesByUser.Categories = element
		db.Save(&categoriesByUser)
	}

	response := &m.ResponseAddCategories{
		InputUsersCategories: inputCategoriesByUser,
		Message:              "New categories have been added for this user",
	}

	c.JSON(http.StatusCreated, response)
}

func ShowOwnPosts(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	auth := &m.Users{}
	err = db.Where("token = ? ", authorization).Find(&auth).Error

	id := auth.ID

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
		posts,
		"Get posts",
		true,
		200,
	}

	c.JSON(http.StatusOK, response)
}

func SeveralCategoriesPostsList(c *gin.Context) {

	inputCategories := &m.InputCategories{}
	err := c.Bind(inputCategories)

	if err != nil {
		response := &m.ResponseAddCategories{
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	posts := []m.PostsUsersJoin{}

	outputPosts := []m.PostsUsersJoin{}

	for _, element := range inputCategories.Categories {
		db.Raw("SELECT * FROM posts join users on posts.id_user = users.id WHERE posts.categories = ? ORDER BY posts.created_at desc", element).Scan(&posts)
		for _, postElement := range posts {
			outputPosts = append(outputPosts, postElement)
		}
	}

	if err != nil {
		response := &m.ResponsePost{
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		c.Abort()
		return
	}

	response := &m.ResponsePost{
		outputPosts, "Posts based on categories have been shown", true, 200,
	}

	c.JSON(http.StatusOK, response)
}
