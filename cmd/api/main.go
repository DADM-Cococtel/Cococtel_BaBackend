package main

import (
	"database/sql"
	"fmt"
	"github.com/Cococtel/Cococtel_BaBackend/internal/http"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	/*
		godotenv.Load(".env")
		db, err := createDB()
		if err != nil {
			panic(err)
		}
	*/
	eng := gin.Default()
	router := http.InitRouter(eng, nil)
	router.MapRoutes()
	if err := eng.Run(); err != nil {
		panic(err)
	}
}
func createDB() (*sql.DB, error) {
	dbUser := os.Getenv("dbUser")
	dbPassword := os.Getenv("dbPassword")
	dbHost := os.Getenv("dbHost")
	dbName := os.Getenv("dbName")
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", dbUser, dbPassword, dbHost, dbName)
	return sql.Open("mysql", connectionString)
}
