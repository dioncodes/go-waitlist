package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var instance *gorm.DB

func Conn() *gorm.DB {
	if instance == nil {
		log.Println("Connecting to DB...")
		instance = connect()
	}

	return instance
}

func connect() *gorm.DB {
	dsn := os.Getenv("DOCKER_MYSQL_USER") + ":" + os.Getenv("DOCKER_MYSQL_PASSWORD") + "@tcp(db:3306)/" + os.Getenv("DOCKER_MYSQL_DATABASE") + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("connected to database")
	return db
}
