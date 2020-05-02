// user-service/database.go
package main

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// CreateConnection ...
func CreateConnection() (*gorm.DB, error) {

	// 获取环境变量，连通数据库
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	// port := os.Getenv("DB_PORT")
	DBName := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")

	return gorm.Open(
		"postgres",
		fmt.Sprintf(
			"postgres://%s:%s@%s/%s?sslmode=disable",
			user, password, host, DBName,
		),
	)
}
