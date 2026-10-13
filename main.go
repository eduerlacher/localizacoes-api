
package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	godotenv.Load()

	db, err := sql.Open("pgx", "host="+os.Getenv("DB_HOST")+
		" port="+os.Getenv("DB_PORT")+
		" user="+os.Getenv("DB_USER")+
		" password="+os.Getenv("DB_PASSWORD")+
		" dbname="+os.Getenv("DB_NAME")+
		" sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.GET("/upas", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, lat, long, title, phone, street, number, neighborhood, city, state, zip_code FROM upas")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
			return
		}
		defer rows.Close()

		var upas []map[string]interface{}

		for rows.Next() {
			var id, title, street, city, state string
			var lat, long float64
			var phone, number, neighborhood, zipCode sql.NullString

			err := rows.Scan(&id, &lat, &long, &title, &phone, &street, &number, &neighborhood, &city, &state, &zipCode)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read UPA"})
				return
			}

			upa := map[string]interface{}{
				"id": id, "lat": lat, "long": long, "title": title,
				"phone": phone.String,
				"address": map[string]interface{}{
					"street": street, "number": number.String,
					"neighborhood": neighborhood.String, "city": city,
					"state": state, "zipCode": zipCode.String,
				},
				"openingHours": "24/7",
			}

			upas = append(upas, upa)
		}
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read UPA rows"})
			return
		}

		c.JSON(http.StatusOK, upas)
	})

	log.Println("API running on http://localhost:8080")
	log.Fatal(r.Run(":8080"))
}
