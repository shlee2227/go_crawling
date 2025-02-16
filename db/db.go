package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // 드라이버 등록만 하고 실제로 사용하진 않아서 _붙임
)

var Conn *sql.DB

func Init() {
	//ENV
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error 
	Conn, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Println("DB 연결 실패", err)
		return 
	}
	if err := Conn.Ping(); err != nil {
		log.Println()
	}
	log.Println("DB 연결 성공")

}

func Close() {
	err := Conn.Close()
	if err != nil {
		log.Println("DB 연결 해제 실패", err)
		return 
	}
	log.Println("DB 연결 해제")
}



