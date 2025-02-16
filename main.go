package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"go_crawling/crawler"
	"go_crawling/db"
)

func ping (c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func db_ping (c *gin.Context) {
	err := db.Conn.Ping()
	if err != nil {
		c.JSON(500, gin.H{"message": "DB 연결 실패"})
	}
	c.JSON(200, gin.H{"message": "DB 연결됨"})
}

func searchAndStore (c *gin.Context) {
	searchWord := c.Query("search_word")
	fmt.Println(searchWord)
	err := crawler.SearchAndStoreData(searchWord)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message":"검색 및 저장에 실패"})
	
	}
	c.JSON(http.StatusCreated, gin.H{"message":"검색 및 저장 성공"})
}
	


func main() {
	// ENV 로드
	if err := godotenv.Load(); err != nil {
		log.Fatal("ENV 로드 실패")
	}

	// DB init
	db.Init()
	defer db.Close()

	// SERVER
	r := gin.Default()

	// basic
	r.GET("/ping", ping)

	// DB
	r.GET("/db_ping", db_ping)

	// crawling 
	r.POST("/search", searchAndStore)
	// r.GET("/data")

	log.Println("Server 시작")
	r.Run(":8000")
	log.Println("Server 시작")
}
