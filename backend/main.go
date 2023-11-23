package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("OK")
	}
	ok, err := db.Exec("CREATE TABLE USER(ID INTEGER,LOGIN TEXT,PASSWORD TEXT);")
	if err != nil {
		panic(err)
		fmt.Println(ok.LastInsertId())
	} else {
		fmt.Println("OK")
	}

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/about.html")
	})
	http.HandleFunc("/postform", func(w http.ResponseWriter, r *http.Request) {

		// TODO: проверить, что пользователь уже существует

		name := r.FormValue("userlogin")
		password := r.FormValue("userpassword")
		_, err := db.Exec("insert into USER(ID, LOGIN, PASSWORD) values (1, $1, $2)", name, password)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("OK")
		}
		http.ServeFile(w, r, "static/about.html")
	})
	http.HandleFunc("/login.html", func(w http.ResponseWriter, r *http.Request) {

		http.ServeFile(w, r, "static/login.html")
	})
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {

		http.ServeFile(w, r, "static/index.html")
	})

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)
}
