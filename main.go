package main

import (
	"log"
	"os"

	"github.com/gertd/go-pluralize"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var plur = pluralize.NewClient()

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_NAME")
	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic(err)
	}
	tables := getTables(db, database)
	for _, v := range tables {
		os.WriteFile("./output/go/models/"+v.TableName+".go", []byte(tblToGo(v)), 0644)
	}
	classes := tablesToClasses(tables)

	for i, v := range classes {
		os.WriteFile("./output/dart/bin/models/"+tables[i].TableName+".dart", []byte(classToDart(v)), 0644)
	}
}
