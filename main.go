package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
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
			log.Println("ERRO SQL UPA:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Database query failed",
				"details": err.Error(),
			})
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
				log.Println("ERRO SCAN UPA:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read UPA"})
				return
			}

			upa := map[string]interface{}{
				"id": id, "lat": lat, "long": long, "title": title,
				"phone": phone.String,
				"address": map[string]interface{}{
					"street": street, "number": number.String,
					"neighborhood": neighborhood.String, "city": city,
					"state": strings.TrimSpace(state), "zipCode": zipCode.String,
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

	r.GET("/ubss", func(c *gin.Context) {
		rows, err := db.Query(`
        SELECT id, lat, long, title, phone,
               street, number, neighborhood, city, state, zip_code,
               monday_open, monday_close,
               tuesday_open, tuesday_close,
               wednesday_open, wednesday_close,
               thursday_open, thursday_close,
               friday_open, friday_close,
               saturday_open, saturday_close,
               sunday_open, sunday_close
        FROM ubss
    `)
		if err != nil {
			log.Println("ERRO SQL UBS:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed", "details": err.Error()})
			return
		}
		defer rows.Close()

		ubss := []map[string]interface{}{}

		formatTime := func(t sql.NullString) interface{} {
			if !t.Valid || len(t.String) < 5 {
				return nil
			}
			return t.String[:5]
		}

		for rows.Next() {
			var id, title, street, city, state string
			var lat, long float64
			var phone, number, neighborhood, zipCode sql.NullString
			var mondayOpen, mondayClose sql.NullString
			var tuesdayOpen, tuesdayClose sql.NullString
			var wednesdayOpen, wednesdayClose sql.NullString
			var thursdayOpen, thursdayClose sql.NullString
			var fridayOpen, fridayClose sql.NullString
			var saturdayOpen, saturdayClose sql.NullString
			var sundayOpen, sundayClose sql.NullString

			err := rows.Scan(
				&id, &lat, &long, &title, &phone,
				&street, &number, &neighborhood, &city, &state, &zipCode,
				&mondayOpen, &mondayClose,
				&tuesdayOpen, &tuesdayClose,
				&wednesdayOpen, &wednesdayClose,
				&thursdayOpen, &thursdayClose,
				&fridayOpen, &fridayClose,
				&saturdayOpen, &saturdayClose,
				&sundayOpen, &sundayClose,
			)
			if err != nil {
				log.Println("ERRO SCAN UBS:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read UBS"})
				return
			}

			openingHours := map[string]interface{}{
				"monday":    map[string]interface{}{"open": formatTime(mondayOpen), "close": formatTime(mondayClose)},
				"tuesday":   map[string]interface{}{"open": formatTime(tuesdayOpen), "close": formatTime(tuesdayClose)},
				"wednesday": map[string]interface{}{"open": formatTime(wednesdayOpen), "close": formatTime(wednesdayClose)},
				"thursday":  map[string]interface{}{"open": formatTime(thursdayOpen), "close": formatTime(thursdayClose)},
				"friday":    map[string]interface{}{"open": formatTime(fridayOpen), "close": formatTime(fridayClose)},
				"saturday":  nil,
				"sunday":    nil,
			}

			if saturdayOpen.Valid || saturdayClose.Valid {
				openingHours["saturday"] = map[string]interface{}{
					"open": formatTime(saturdayOpen), "close": formatTime(saturdayClose),
				}
			}

			if sundayOpen.Valid || sundayClose.Valid {
				openingHours["sunday"] = map[string]interface{}{
					"open": formatTime(sundayOpen), "close": formatTime(sundayClose),
				}
			}

			ubs := map[string]interface{}{
				"id":    id,
				"lat":   lat,
				"long":  long,
				"title": title,
				"phone": phone.String,
				"address": map[string]interface{}{
					"street":       street,
					"number":       number.String,
					"neighborhood": neighborhood.String,
					"city":         city,
					"state":        strings.TrimSpace(state),
					"zipCode":      zipCode.String,
				},
				"openingHours": openingHours,
			}

			ubss = append(ubss, ubs)
		}
		if err := rows.Err(); err != nil {
			log.Println("UBS rows error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read UBS records"})
			return
		}

		c.JSON(http.StatusOK, ubss)
	})
	log.Println("API running on http://localhost:8080")
	log.Fatal(r.Run(":8080"))
}