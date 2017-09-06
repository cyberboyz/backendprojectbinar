package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error

type Food struct {
	ID        uint   `json:"id"`
	FoodName string `json:"foodname"`
	Price  uint `json:"price"`
}

// START OMIT
func main() {
	db, err = gorm.Open("postgres", "host=localhost user=postgres dbname=gorm sslmode=disable password=postgres")

	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	db.AutoMigrate(&Food{})

	r := gin.Default()
	r.GET("/v1/food", GetAllFood) // HL
	r.GET("/v1/food/:id", GetFood)
	r.POST("/v1/food", CreateFood)
	r.PUT("/v1/food/:id", UpdateFood)
	r.PATCH("/v1/food/:id", UpdateFood)
	r.DELETE("/v1/food/:id", DeleteFood)
	r.Run(":8080")
}

// END OMIT

func DeleteFood(c *gin.Context) {
	id := c.Params.ByName("id")
	var food Food
	d := db.Where("id = ?", id).Delete(&food)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}

func UpdateFood(c *gin.Context) {
	var food Food
	id := c.Params.ByName("id")
	if err := db.Where("id = ?", id).First(&food).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	c.BindJSON(&food)
	db.Save(&food)
	c.JSON(200, food)
}

func CreateFood(c *gin.Context) {
	var food Food
	c.BindJSON(&food)
	db.Create(&food)
	c.JSON(200, food)
}

func GetFood(c *gin.Context) {
	id := c.Params.ByName("id")
	var food Food
	if err := db.Where("id = ?", id).First(&food).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, food)
	}
}

func GetAllFood(c *gin.Context) {
	var somefood []Food
	if err := db.Find(&somefood).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, somefood)
	}

}