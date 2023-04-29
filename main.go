package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type stammform struct {
	id int

	value  string
	wortId int
}

type wort struct {
	id int

	name         string
	wortschatzId int
}

type wortschatz struct {
	id int

	name string
}

type jsonVOC struct {
	Ws          string   `json:"ws"`
	Stammformen []string `json:"stammformen"`
	Name        string   `json:"name"`
}

type jsonDB struct {
	Vocs  []jsonVOC   `json:"vocs"`
	Texts []jsonTexte `json:"texte"`
}

type jsonTexte struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

func main() {
	db, err := sql.Open("mysql", "root:root@(127.0.0.1:3308)/TestDb?parseTime=true")

	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	{ // Creating tables
		query := `DROP TABLE IF EXISTS stammform;`
		execSQLQuery(db, query)

		query = `
			CREATE TABLE IF NOT EXISTS stammform (
				id INT AUTO_INCREMENT,
				value TEXT NOT NULL,
				wordID INT NOT NULL,
				PRIMARY KEY (id)
			);`
		execSQLQuery(db, query)

		query = `DROP TABLE IF EXISTS wort;`
		execSQLQuery(db, query)

		query = `
			CREATE TABLE IF NOT EXISTS wort (
				id INT AUTO_INCREMENT,
				name TEXT NOT NULL,
				wortschatzId TEXT NOT NULL,
				PRIMARY KEY (id)
			);`
		execSQLQuery(db, query)

		query = `DROP TABLE IF EXISTS wortschatz;`
		execSQLQuery(db, query)
		// DONT THINK ITS NEEDED
		// query = `
		// 	CREATE TABLE IF NOT EXISTS wortschatz (
		// 		id INT AUTO_INCREMENT,
		// 		name TEXT NOT NULL,
		// 		PRIMARY KEY (id)
		// 	);`
		// execSQLQuery(db, query)
	}

	var jsondb jsonDB

	{ // insert vocs into db
		jsonFile, err := os.Open("Vocabel.json")

		if err != nil {
			log.Fatal(err)
		}

		defer jsonFile.Close()

		byteValue, err := ioutil.ReadAll(jsonFile)

		json.Unmarshal(byteValue, &jsondb)

		for i := 0; i < len(jsondb.Vocs); i++ {
			curVoc := jsondb.Vocs[i]

			result, err := db.Exec(`INSERT INTO wort (name, wortschatzId) VALUES (?,?)`, curVoc.Name, curVoc.Ws)

			if err != nil {
				log.Fatal(err)
			}

			id, err := result.LastInsertId()

			for j := 0; j < len(curVoc.Stammformen); j++ {
				curStammform := curVoc.Stammformen[j]

				result, err := db.Exec(`INSERT INTO stammform (value, wordID) VALUES (?,?)`, curStammform+"%", id)

				if err != nil {
					log.Fatal(err)
				}

				fmt.Println(result.LastInsertId())
			}
		}
	}

	for _, text := range jsondb.Texts {
		fmt.Println("text: ", text.Name)

		{ // wertet den aids aus
			testText := text.Text

			words := strings.FieldsFunc(testText, func(r rune) bool {
				return r == ' ' || r == ',' || r == '.' || r == '!'
			})

			totalKnow := 0

			query := "SELECT COUNT(*) FROM stammform WHERE ? LIKE value"

			for _, word := range words {
				var count int

				err = db.QueryRow(query, word).Scan(&count)

				if err != nil {
					log.Fatal(err)
				}

				if count > 0 {
					totalKnow = totalKnow + 1

					fmt.Println(word, totalKnow)
				}
			}

			var percentage float64

			percentage = float64(totalKnow) / float64(len(words)) * 100

			fmt.Printf("total word count: %d\n", len(words))
			fmt.Printf("total known word percentage: %f\n", percentage)
		}
	}

}

func execSQLQuery(db *sql.DB, query string) {
	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
	}
}
