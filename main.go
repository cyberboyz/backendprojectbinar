package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error

type Post struct {
	IDPost       uint   `json:"id_post"`
	IDUser       uint   `json:"id_user"`
	IDBackground uint   `json:"id_background"`
	PostTitle    string `json:"post_title"`
	Categories   string `json:"categories"`
	Content      string `json:"content"`
}

// START OMIT
func main() {
	db, err = gorm.Open("postgres", "host=localhost user=postgres dbname=gorm sslmode=disable password=postgres")

	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	db.AutoMigrate(&Post{})

	r := gin.Default()
	r.POST("/v1/posts", CreateNewPost)
	r.GET("/v1/posts", GetAllPosts) // HL
	r.GET("/v1/posts/:id_post", GetPost)
	r.PUT("/v1/posts/:id_post", UpdatePost)
	r.PATCH("/v1/posts/:id_post", UpdatePost)
	r.DELETE("/v1/posts/:id_post", DeletePost)
	r.Run(":8080")
}

// END OMIT

func CreateNewPost(c *gin.Context) {
	var post Post
	c.BindJSON(&post)
	db.Create(&post)
	c.JSON(200, post)
}

func GetAllPosts(c *gin.Context) {
	var someposts []Post
	if err := db.Find(&someposts).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, someposts)
	}

}

func GetPost(c *gin.Context) {
	id_post := c.Params.ByName("id_post")
	fmt.Println(id_post)
	var post Post
	if err := db.Where("id_post = ?", id_post).First(&post).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, post)
	}
}

func UpdatePost(c *gin.Context) {
	var post Post
	id_post := c.Params.ByName("id_post")
	if err := db.Where("id_post = ?", id_post).First(&post).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	c.BindJSON(&post)
	db.Save(&post)
	c.JSON(200, post)
}

func DeletePost(c *gin.Context) {
	id_post := c.Params.ByName("id_post")
	var post Post
	d := db.Where("id_post = ?", id_post).Delete(&post)
	fmt.Println(d)
	c.JSON(200, gin.H{"id_post #" + id_post: "deleted"})
}
