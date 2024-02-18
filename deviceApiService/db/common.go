package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreateDb(host string, dbName string) {
	dsn := fmt.Sprintf("host=%s port=5432 user=%s password=%s", host, "postgres", "admin")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Exec("CREATE DATABASE " + dbName + ";")
}

func GetDb(host string, pgPort string, dbName string) *gorm.DB {

	CreateDb(host, dbName)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, pgPort, "postgres", "admin", dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
