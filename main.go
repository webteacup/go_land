package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

var db *sql.DB

// album represents data about a record album.
type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		fmt.Printf("getAlbums: %v", err)
	}
	var albums []Album
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			fmt.Printf("Error: %v", err)
		} else {
			albums = append(albums, alb)
		}
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("Error: %v", err)
	}
	c.IndentedJSON(http.StatusOK, albums)
}

// addAlbum adds a new album to the database.
func addAlbum(c *gin.Context) {
	var newAlbum Album
	if err := c.BindJSON(&newAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": "Successfully created!"})
}

// updateAlbum updates a exist album
func updateAlbum(c *gin.Context) {
	var updatedAlbum Album
	if err := c.BindJSON(&updatedAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := db.Exec("UPDATE album SET title = ?, artist = ?, price = ? WHERE id = ?", updatedAlbum.Title, updatedAlbum.Artist, updatedAlbum.Price, updatedAlbum.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "Successfully updated!"})
}

func main() {

	var err error
	db, err = sql.Open("mysql", "root@/recordings")
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected")

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.PUT("/albums", addAlbum)
	router.POST("/albums", updateAlbum)
	router.Run("localhost:8080")
}
