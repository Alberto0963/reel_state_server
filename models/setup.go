package models

import (
	// "github.com/gin-gonic/gin"
	"gorm.io/gorm"
	// "net/http"
	// "log"
	"gorm.io/driver/mysql"
	"os"
	// "github.com/gin-gonic/gin"
	// "database/sql"
	// "net/http"
	// "log"
	// "os"
	"fmt"
	"log"
	// "os"

	// "github.com/jinzhu/gorm"
	//  _ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

// var DB *gorm.DB

// func ConnectDataBase(){

// 	err := godotenv.Load(".env")

// 	if err != nil {
// 	  log.Fatalf("Error loading .env file")
// 	}

// 	Dbdriver := os.Getenv("DB_DRIVER")
// 	DbHost := os.Getenv("DB_HOST")
// 	DbUser := os.Getenv("DB_USER")
// 	DbPassword := os.Getenv("DB_PASSWORD")
// 	DbName := os.Getenv("DB_NAME")
// 	DbPort := os.Getenv("DB_PORT")

// 	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)

// 	DB, err = gorm.Open(Dbdriver, DBURL)

// 	if err != nil {
// 		fmt.Println("Cannot connect to database ", Dbdriver)
// 		log.Fatal("connection error:", err)
// 	} else {
// 		fmt.Println("We are connected to the database ", Dbdriver)
// 	}

// 	DB.AutoMigrate(&User{})

// }

var Pool *gorm.DB // Declare a global variable to hold the connection pool

func InitDB() (*gorm.DB, error) {

	err := godotenv.Load("/home/reelstate/go/reel_state_server" + "/.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	Dbdriver := os.Getenv("DB_DRIVER")
	DbHost := os.Getenv("DB_HOST")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbPort := os.Getenv("DB_PORT")

	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
	// dsn := os.Getenv(DBURL)

	db, err := gorm.Open(mysql.Open(DBURL), &gorm.Config{})
	if err != nil {
		return nil, err
	} else {
		fmt.Println("We are connected to the database ", Dbdriver)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set the maximum number of connections in the pool
	sqlDB.SetMaxOpenConns(10)
	Pool = db
	return db, nil
}
