package db

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB


func RegisterRoutes(r *gin.Engine){
	r.GET("/db_ping", db_ping)
}

func Init() {
	//ENV
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul", dbHost, dbUser, dbPassword, dbName, dbPort)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("DB 연결 실패", err)
		return
	}
	log.Println("DB 연결 성공")
}

func Close() {
    sqlDB, err := DB.DB()
    if err != nil {
        log.Println("DB 객체 획득 실패", err)
        return
    }
    err = sqlDB.Close()
    if err != nil {
        log.Println("DB 연결 해제 실패", err)
        return
    }
    log.Println("DB 연결 해제")
}

func db_ping(c *gin.Context) {
    sqlDB, err := DB.DB()
    if err != nil {
        c.JSON(500, gin.H{"message": "DB 객체 획득 실패"})
        return
    }
    err = sqlDB.Ping()
    if err != nil {
        c.JSON(500, gin.H{"message": "DB 연결 실패"})
        return
    }
    c.JSON(200, gin.H{"message": "DB 연결됨"})
}
