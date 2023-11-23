package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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
	http.HandleFunc("/register.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/register.html")
	})
	http.HandleFunc("/postform", func(w http.ResponseWriter, r *http.Request) {

		name := r.FormValue("userlogin")
		password := r.FormValue("userpassword")

		var count int
		err1 := db.QueryRow("SELECT COUNT(*) FROM USER WHERE LOGIN = $1", name).Scan(&count)
		if err1 != nil {
			panic(err1)
		}
		if count > 0 {
			fmt.Println("Exists")
			http.ServeFile(w, r, "static/used.html")
		} else {
			fmt.Println("new user ", name)
			_, err := db.Exec("insert into USER(ID, LOGIN, PASSWORD) values (1, $1, $2)", name, password)
			if err != nil {
				panic(err)
			} else {
				fmt.Println("OK")
			}
			http.ServeFile(w, r, "static/about.html")
		}
	})
	http.HandleFunc("/bebrik", func(w http.ResponseWriter, r *http.Request) {

		// TODO: проверить, существует ли аккаунт DONE

		name := r.FormValue("userlogin")
		password := r.FormValue("userpassword")

		var count int
		err2 := db.QueryRow("SELECT COUNT(*) FROM USER WHERE LOGIN = $1", name).Scan(&count)
		if err2 != nil {
			panic(err2)
		}

		// TODO: проверить, совпадает ли введенный пароль и проль в бд DONE

		if count > 0 {
			var cellContent string
			err3 := db.QueryRow("SELECT PASSWORD FROM USER WHERE LOGIN = $1", name).Scan(&cellContent)
			if err3 != nil {
				panic(err3)
			}
			if cellContent != password {
				http.ServeFile(w, r, "static/not_find.html")
			}
			http.ServeFile(w, r, "static/index.html")
		} else {
			http.ServeFile(w, r, "static/not_find.html")
		}
		fmt.Println(name, password)
	})
	http.HandleFunc("/login.html", func(w http.ResponseWriter, r *http.Request) {

		http.ServeFile(w, r, "static/login.html")
	})
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {

		http.ServeFile(w, r, "static/index.html")
	})
	http.HandleFunc("/abobus", func(w http.ResponseWriter, r *http.Request) {
		query := r.FormValue("request")
		fmt.Println(query)
		values := map[string]string{"query": query}
		json_data, err := json.Marshal(values)

		if err != nil {
			log.Fatal(err)
		}

		resp, err4 := http.Post("http://127.0.0.1:8080/", "application/json",
			bytes.NewBuffer(json_data))
		if err4 != nil {
			panic(err4)
		}

		defer resp.Body.Close()
		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		fmt.Println(res["key"])
		if res["key"] == "Hentai" {
			http.ServeFile(w, r, "static/hentai.html")
		} else {
			http.ServeFile(w, r, "static/comedy.html")
		}
	})
	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)
}
