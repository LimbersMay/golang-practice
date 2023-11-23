package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

func findAlbums(c *gin.Context) {

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	var albums []Album

	db.Find(&albums)

	c.IndentedJSON(200, albums)
}

func createNewAlbum(c *gin.Context) {

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	var newAlbum Album

	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	db.Create(&newAlbum)

	c.IndentedJSON(201, newAlbum)
}

func main() {

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Album{})

	router := gin.Default()
	router.GET("/albums", findAlbums)

	router.POST("/albums", createNewAlbum)

	router.Run("localhost:3000")
}
