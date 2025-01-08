package main

import (
	"database/sql"
	"github.com/Cococtel/Cococtel_BaBackend/internal/http"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func main() {
	/*
		godotenv.Load(".env")
		db, err := createDB()
		if err != nil {
			panic(e	rr)
		}
	*/
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "cococtel",
		AllowNativePasswords: true,
	}
	db, err := createDB(cfg)
	if err != nil {
		panic(err)
	}

	eng := gin.Default()
	router := http.InitRouter(eng, db)
	router.MapRoutes()
	if err := eng.Run(); err != nil {
		panic(err)
	}
}
func createDB(cfg mysql.Config) (*sql.DB, error) {
	//dbUser := "root"
	//dbPassword := ""
	//dbHost := "localhost:3306"
	//dbName := "cococtel"
	//var connectionString = fmt.Sprintf("%s:@tcp(%s)/%s?charset=utf8", dbUser, dbHost, dbName)

	return sql.Open("mysql", cfg.FormatDSN())
}
