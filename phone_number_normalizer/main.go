package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	_ "database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type phone struct {
	id     int
	number string
}

const (
	dbname = "normalizer_db"
	host   = "localhost"
	port   = "5432"
	user   = "postgres"
)

func main() {
	phoneNumber := flag.String("phone", "11111", "Introduce the phone number you want to save")
	normalize := flag.Bool("normalize", true, "Indicate you want to normalize all value to the same format")
	flag.Parse()
	fmt.Println(*phoneNumber)
	if *phoneNumber == "" {
		writer := os.Stdout
		writer.WriteString("Plesae introduce exact one phone number\n")
		os.Exit(1)
	}

	db := connectToDB()
	db.MustExec("INSERT INTO phone_number (id,number) VALUES (nextval('phone_id_seq'),$1)", *phoneNumber)

	if *normalize {
		normalizeNumbers(db)
	}
}

func connectToDB() *sqlx.DB {
	connectionInfo := fmt.Sprintf("host=%s port=%s dbname=%s user=%s sslmode=disable", host, port, dbname, user)
	return sqlx.MustConnect("postgres", connectionInfo)
}

func normalizeNumbers(db *sqlx.DB) {
	regex := regexp.MustCompile(`[^1234567890]`)
	rows, _ := db.Query("SELECT id,number FROM phone_number")
	var number string
	var id string
	for {
		if rows.Next() {
			err := rows.Scan(&id, &number)
			if err != nil {
				panic(err)
			}

			normalizedValue := strings.Join(strings.Fields(regex.ReplaceAllLiteralString(number, " ")), "")
			existQuery := "SELECT id FROM phone_number where number = $1"

			existsRows, _ := db.Query(existQuery, normalizedValue)
			for {
				if existsRows.Next() {
					var dupIndex string
					existsRows.Scan(&dupIndex)
					if dupIndex != id {
						deleteNumber := "DELETE FROM phone_number where id = $1"
						db.MustExec(deleteNumber, dupIndex)
					}
				} else {
					break
				}
			}
			updateQuery := "UPDATE phone_number SET number = $1 where id = $2"
			db.MustExec(updateQuery, normalizedValue, id)
		} else {
			break
		}
	}
}
