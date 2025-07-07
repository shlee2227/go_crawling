package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/shlee2227/go_crawling/internal/middleware/db"

	controller "github.com/shlee2227/go_crawling/internal/controller/search"
	repository "github.com/shlee2227/go_crawling/internal/repository/search"
	service "github.com/shlee2227/go_crawling/internal/service/search"
)

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func main() {
	// ENV 로드s
	if err := godotenv.Load(); err != nil {
		log.Fatal("ENV 로드 실패")
	}

	// DB init
	db.Init()
	defer db.Close()

	// 의존성 주입
	repository := repository.NewRepository(db.DB)
	service := service.NewService(repository)
	controller := controller.NewNaverController(service)


	// SERVER
	r := gin.Default()

	// basic
	r.GET("/ping", ping)

	// Routes
	db.RegisterRoutes(r)
	controller.RegisterRoutes(r)

	log.Println("Server 시작")
	r.Run(":8000")
	log.Println("Server 시작")
}
