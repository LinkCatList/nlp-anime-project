package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Login string
	Rank  int
}

func main() {
	fmt.Println("OK")

	db, err := sql.Open("sqlite3", "bebra.db")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Ok open db")
	}
	ok, err := db.Exec("CREATE TABLE IF NOT EXISTS USER(ID INTEGER,LOGIN TEXT,PASSWORD TEXT, RANK INTEGER);")
	if err != nil {
		panic(err)
		fmt.Println(ok.LastInsertId())
	} else {
		fmt.Println("Ok create db")
	}
	_, err1 := db.Exec(`insert into USER(ID, LOGIN, PASSWORD, RANK) values (1, "okokokok", "1234", 282)`)
	if err1 != nil {
		panic(err1)
	} else {
		fmt.Println("Ok insert user")
	}
	h1 := func(w http.ResponseWriter, r *http.Request) {
		templ := template.Must(template.ParseFiles("index.html"))
		rows, err := db.Query("SELECT * FROM USER")
		if err != nil {
			panic(err)
		}

		bebra := []User{}
		for rows.Next() {
			var id, rank int
			var login, password string
			err := rows.Scan(&id, &login, &password, &rank)
			if err != nil {
				panic(err)
			}
			bebra = append(bebra, User{login, rank})
		}
		fmt.Println(bebra)

		sort.SliceStable(bebra, func(i, j int) bool {
			return bebra[i].Rank > bebra[j].Rank
		})
		anime := map[string][]User{"Users": {}}
		for _, value := range bebra {
			anime["Users"] = append(anime["Users"], value)
		}
		fmt.Println(anime)
		templ.Execute(w, anime)
	}

	http.HandleFunc("/raiting", h1)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
