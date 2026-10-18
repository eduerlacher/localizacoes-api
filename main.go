package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"

	"localizacoes-api/internal/controller"
	"localizacoes-api/internal/model"
)

func main() {
	godotenv.Load()

	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.SetTrustedProxies(nil)

	unidadeModel := model.New(db)
	unidadeController := controller.NewUnidadeController(unidadeModel)

	r.GET("/upas", unidadeController.ListUPAs)
	r.GET("/ubss", unidadeController.ListUBSs)

	log.Println("API running on http://localhost:8080")
	log.Fatal(r.Run(":8080"))
}
