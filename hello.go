package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {

	albums := getAlbumsFromDB()

	c.IndentedJSON(http.StatusOK, albums)
}

func createAlbum(c *gin.Context) {
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	insertAlbum(newAlbum.ID, newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func createTable(db *sql.DB) {
	createStudentTableSQL := `CREATE TABLE Album (
    		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    		"title" TEXT,
    		"artist" TEXT,
    		"price" FLOAT
	  );` // SQL Statement for Create Table

	log.Println("Create album table...")
	statement, err := db.Prepare(createStudentTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("album table created")
}

func insertAlbum(id string, title string, artist string, price float64) {

	db, _ := sql.Open("sqlite3", "./albums.db")

	log.Println("Inserting student record ...")
	insertStudentSQL := `INSERT INTO Album(id, title, artist, price) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insertStudentSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(id, title, artist, price)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func getAlbumsFromDB() []album {

	db, _ := sql.Open("sqlite3", "./albums.db")
	rows, err := db.Query("SELECT id, title, artist, price FROM Album")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var albums []album

	for rows.Next() {
		var id string
		var title string
		var artist string
		var price float64
		rows.Scan(&id, &title, &artist, &price)
		albums = append(albums, album{ID: id, Title: title, Artist: artist, Price: price})
	}

	return albums
}

func main() {

	sqliteDatabase, _ := sql.Open("sqlite3", "./albums.db") // Open the created SQLite File
	defer sqliteDatabase.Close()                            // Defer Closing the database

	router := gin.Default()
	router.GET("/albums", getAlbums)

	router.POST("/albums", createAlbum)

	router.GET("/albums/:id", getAlbumByID)

	router.Run("localhost:8080")
}
