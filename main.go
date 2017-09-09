package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Users struct {
	ID           uint      `gorm:"primary_key" json:"id_user"`
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	Bio          string    `json:"bio"`
	IDAvatar     int       `json:"id_avatar"`
	IDCoverPhoto int       `json:"id_cover_photo"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Posts struct {
	gorm.Model
	IDUser       uint   `json:"id_user"`
	IDBackground uint   `json:"id_background"`
	PostTitle    string `json:"post_title"`
	Categories   string `json:"categories"`
	Content      string `json:"content"`
}

type ResponseUser struct {
	Users   interface{} `json:"users, omitempty"`
	Message string      `json:"message"`
}

type ResponsePost struct {
	Posts   interface{} `json:"posts, omitempty"`
	Message string      `json:"message"`
}

type Login struct {
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var db *gorm.DB
var err error

func main() {
	db, err = gorm.Open("postgres", "host=localhost user=postgres dbname=dreamcatcher sslmode=disable password=postgres")
	db.SingularTable(true)
	if err != nil {
		log.Panic(err)
	}
	db.AutoMigrate(&Users{})

	router := gin.Default()

	v1 := router.Group("/v1")
	{
		usergroup := v1.Group("/profile")
		{
			usergroup.GET("/", UserGet)
			usergroup.GET("/:id", UserDetail)
			usergroup.POST("/", UserCreate)
			usergroup.PUT("/:id", UserUpdate)
			usergroup.PATCH("/:id", UserUpdate)
			usergroup.DELETE("/:id", UserDelete)
		}
		v1.POST("/login", LoginUser)
	}
	router.Run(":8080")
}

func LoginUser(c *gin.Context) {
	login := &Login{}
	if c.BindJSON(login) == nil {
		if login.User == "ngetes" && login.Password == "123" {
			c.JSON(http.StatusOK, gin.H{"status": "Login successful"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized"})
		}
	}
}

func Authorize(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	if authorization != "12345" {
		response := &ResponseUser{
			Message: "Unauthorized access",
		}
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}
}

func UserGet(c *gin.Context) {
	Authorize(c)

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
	Authorize(c)

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
	Authorize(c)

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
	Authorize(c)

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

	Authorize(c)

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
