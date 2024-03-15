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
	router.Run("localhost:8080")
}
