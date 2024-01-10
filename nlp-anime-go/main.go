package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
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
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("OK")
	}
	ok, err := db.Exec("CREATE TABLE IF NOT EXISTS USER(ID INTEGER,LOGIN TEXT,PASSWORD TEXT, RANK INTEGER);")
	if err != nil {
		panic(err)
		fmt.Println(ok.LastInsertId())
	} else {
		fmt.Println("OK")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "html/homepage.html")
	})

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "html/login.html")
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
			_, err := db.Exec("insert into USER(ID, LOGIN, PASSWORD, RANK) values (1, $1, $2, 0)", name, password)
			if err != nil {
				panic(err)
			} else {
				fmt.Println("OK")
			}
			http.ServeFile(w, r, "static/about.html")
		}
	})
	http.HandleFunc("/bebrik", func(w http.ResponseWriter, r *http.Request) {

		name := r.FormValue("userlogin")
		password := r.FormValue("userpassword")

		var count int
		err2 := db.QueryRow("SELECT COUNT(*) FROM USER WHERE LOGIN = $1", name).Scan(&count)
		if err2 != nil {
			panic(err2)
		}

		if count > 0 {
			var cellContent string
			err3 := db.QueryRow("SELECT PASSWORD FROM USER WHERE LOGIN = $1", name).Scan(&cellContent)
			if err3 != nil {
				panic(err3)
			}
			if cellContent != password {
				http.ServeFile(w, r, "static/not_find.html")
			}
			cookie := http.Cookie{
				Name:  "name",
				Value: name,
				Path:  "/",
			}
			http.SetCookie(w, &cookie)
			http.ServeFile(w, r, "html/index.html")
		} else {
			http.ServeFile(w, r, "static/not_find.html")
		}

		fmt.Println(name, password)
	})
	http.HandleFunc("/zalupa_slonika", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "html/index.html")
	})
	http.HandleFunc("/comedy", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "html/comedy.html")
	})
	http.HandleFunc("/abobus", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("name")

		if err != nil || cookie.Value == "" {
			http.ServeFile(w, r, "static/login.html")
			return
		}
		var UserName = cookie.Value
		fmt.Println(UserName)
		var cnt = 0
		err5 := db.QueryRow("SELECT RANK FROM USER WHERE LOGIN = $1", UserName).Scan(&cnt)
		fmt.Println(cnt)
		cnt++
		if err5 != nil {
			panic(err5)
		} else {
			fmt.Println("OK")
		}
		fmt.Println(UserName)
		fmt.Println(cnt)
		_, err6 := db.Exec("UPDATE USER SET RANK = $1 WHERE LOGIN = $2", cnt, UserName)
		if err6 != nil {
			panic(err6)
		} else {
			fmt.Println("OK")
		}

		query := r.FormValue("request")
		fmt.Println(query)
		values := map[string]string{"query": query}
		json_data, err := json.Marshal(values)

		if err != nil {
			log.Fatal(err)
		}

		resp, err4 := http.Post("http://localhost:3000/", "application/json",
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
		} else if res["key"] == "Comedy" {
			http.ServeFile(w, r, "static/comedy.html")
		} else if res["key"] == "Kids" {
			http.ServeFile(w, r, "static/kids.html")
		} else if res["key"] == "Drama" {
			http.ServeFile(w, r, "static/drama.html")
		} else if res["key"] == "Adventure" {
			http.ServeFile(w, r, "static/adventure.html")
		} else if res["key"] == "Fantasy" {
			http.ServeFile(w, r, "static/fantasy.html")
		} else if res["key"] == "Sci-Fi" {
			http.ServeFile(w, r, "static/scifi.html")
		} else if res["key"] == "Music" {
			http.ServeFile(w, r, "static/music.html")
		} else if res["key"] == "Slice" {
			http.ServeFile(w, r, "static/slice.html")
		} else if res["key"] == "Action" {
			http.ServeFile(w, r, "static/action.html")
		}
	})
	h1 := func(w http.ResponseWriter, r *http.Request) {
		templ := template.Must(template.ParseFiles("html/lederboard.html"))
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
	fmt.Println("Server is listening...")
	http.ListenAndServe(":"+"3001", nil)
}
