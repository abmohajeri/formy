package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var db *gorm.DB

func DatabaseSetup() {
	dsn := "host=" + os.Getenv("DATABASE_HOST") + " user=" + os.Getenv("DATABASE_USERNAME") + " password=" + os.Getenv("DATABASE_PASSWORD") + " dbname=" + os.Getenv("DATABASE_DB_NAME") + " port=" + os.Getenv("DATABASE_PORT") + " sslmode=disable TimeZone=Asia/Tehran"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}

func GetDB() *gorm.DB {
	return db
}
