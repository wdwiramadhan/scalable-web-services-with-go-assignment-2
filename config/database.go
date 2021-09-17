package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var DB *sql.DB

type DBConfig struct{
	Host     string
	Port     int
	User     string
	DBName   string
	Password string
}

func BuildDBConfig() *DBConfig{
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	DBHost := os.Getenv("DB_HOST")
	DBUser := os.Getenv("DB_USER")
	DBPass := os.Getenv("DB_PASS")
	DBName := os.Getenv("DB_NAME")

	dbConfig := DBConfig{
		Host:     DBHost,
		Port:     3306,
		User:     DBUser,
		Password: DBPass,
		DBName:   DBName,
	}
	return &dbConfig
}

func DBUrl(dbConfig *DBConfig) string{
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)
}