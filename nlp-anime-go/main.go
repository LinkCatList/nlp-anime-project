package main

import (
	"bytes"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
	"unicode"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Login string
	Rank  int
	Id    int
}
type Data struct {
	Val  []string `json:"lables"`
	Data []int    `json:"data"`
}

func CorrectLogin(s string) string {
	var t string
	for _, val := range s {
		if val != ' ' {
			t += string(val)
		}
	}
	return t
}

func generateLineItems() Data {
	var items Data
	items.Data = append(items.Data, 1)
	items.Data = append(items.Data, 2)
	items.Data = append(items.Data, 1)
	items.Data = append(items.Data, 3)

	items.Val = append(items.Val, "1")
	items.Val = append(items.Val, "2")
	items.Val = append(items.Val, "3")
	items.Val = append(items.Val, "4")

	return items
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

	ok2, err2 := db.Exec("CREATE TABLE IF NOT EXISTS QUERIES(LOGIN TEXT, DATE TEXT);")
	if err2 != nil {
		fmt.Println("error while create db queries")
		fmt.Println(ok2)
	} else {
		fmt.Println("OK")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "html/homepage.html")
	})

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "html/login.html")
	})
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "html/register.html")
	})
	http.HandleFunc("/postform", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		password := r.FormValue("password")
		name = CorrectLogin(name)
		fmt.Println(name, password)

		var count int
		err1 := db.QueryRow("SELECT COUNT(*) FROM USER WHERE LOGIN = $1", name).Scan(&count)
		if err1 != nil {
			panic(err1)
		}
		if count > 0 {
			fmt.Println("Exists")
			http.ServeFile(w, r, "html/login.html")
		} else {
			hash := sha512.Sum512([]byte(name + password))
			hashedPassword := hex.EncodeToString(hash[:])
			fmt.Println("new user ", name)
			_, err := db.Exec("insert into USER(ID, LOGIN, PASSWORD, RANK) values (1, $1, $2, 0)", name, hashedPassword)
			if err != nil {
				panic(err)
			} else {
				fmt.Println("OK")
			}
			http.ServeFile(w, r, "html/homepage.html")
		}
	})
	http.HandleFunc("/bebrik", func(w http.ResponseWriter, r *http.Request) {

		name := r.FormValue("name")
		password := r.FormValue("password")

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

			hash := sha512.Sum512([]byte(name + password))
			hashedPassword := hex.EncodeToString(hash[:])

			if hashedPassword != cellContent {
				http.ServeFile(w, r, "html/login.html")
			}
			cookie := http.Cookie{
				Name:  "name",
				Value: hashedPassword,
				Path:  "/",
			}
			http.SetCookie(w, &cookie)
			http.ServeFile(w, r, "html/index.html")
		} else {
			http.ServeFile(w, r, "html/login.html")
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
			http.ServeFile(w, r, "html/login.html")
			return
		}
		var UserName = cookie.Value
		fmt.Println(UserName)
		var cnt = 0
		err5 := db.QueryRow("SELECT RANK FROM USER WHERE PASSWORD = $1", UserName).Scan(&cnt)
		_, err228 := db.Exec("insert into QUERIES(LOGIN, DATE) values ($1, $2)", UserName, time.Now())
		if err228 != nil {
			panic(err228)
		} else {
			fmt.Println("OK insert into Queries")
		}
		fmt.Println(cnt)
		cnt++
		if err5 != nil {
			http.ServeFile(w, r, "html/register.html")
			panic(err5)
		} else {
			fmt.Println("OK")
		}
		fmt.Println(UserName)
		fmt.Println(cnt)
		_, err6 := db.Exec("UPDATE USER SET RANK = $1 WHERE PASSWORD = $2", cnt, UserName)
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
		var res map[string]string
		json.NewDecoder(resp.Body).Decode(&res)
		fmt.Println(res["key"])
		tmp := ""
		for _, val := range res["key"] {
			if unicode.IsLetter(val) {
				tmp += string(val)
			}
		}
		path := "html/" + strings.ToLower(tmp) + ".html" // если перестанет работать то виноват терминейт!!!!!
		http.ServeFile(w, r, path)
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
			bebra = append(bebra, User{login, rank, 0})
		}
		fmt.Println(bebra)

		sort.SliceStable(bebra, func(i, j int) bool {
			return bebra[i].Rank > bebra[j].Rank
		})
		anime := map[string][]User{"Users": {}}
		var curId int
		for _, value := range bebra {
			curId++
			var Ok User
			Ok.Login = value.Login
			Ok.Rank = value.Rank
			Ok.Id = curId
			anime["Users"] = append(anime["Users"], Ok)
		}

		fmt.Println(anime)
		templ.Execute(w, anime)
	}
	http.HandleFunc("/raiting", h1)

	h2 := func(w http.ResponseWriter, r *http.Request) {
		cookie, err1 := r.Cookie("name")
		fmt.Println(err1)
		if err1 != nil {
			http.ServeFile(w, r, "html/login.html")
			return
		} else {
			templ := template.Must(template.ParseFiles("html/profile.html"))
			var Name string
			err := db.QueryRow("SELECT LOGIN FROM USER WHERE PASSWORD = $1", cookie.Value).Scan(&Name)
			if err != nil {
				fmt.Println("not ok")
			}
			var Rnk int
			err2 := db.QueryRow("SELECT RANK FROM USER WHERE PASSWORD = $1", cookie.Value).Scan(&Rnk)
			if err2 != nil {
				fmt.Println("not ok")
			}
			var People User
			People.Login = Name
			People.Rank = Rnk
			templ.Execute(w, People)
		}
	}
	http.HandleFunc("/profile", h2)

	h3 := func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("name")
		rows, err := db.Query("SELECT `DATE` FROM QUERIES WHERE LOGIN = $1", cookie.Value)
		if err != nil {
			fmt.Println("error while get dates from queries db")
		}
		defer rows.Close()
		var items []string
		for rows.Next() {
			var it string
			if err := rows.Scan(&it); err != nil {
				fmt.Println(err)
			}
			items = append(items, it[:10])
		}
		fmt.Println(items)
		var a Data
		mp := make(map[string]int)
		for _, val := range items {
			mp[val]++
		}
		for key, val := range mp {
			a.Val = append(a.Val, key)
			a.Data = append(a.Data, val)
		}
		fmt.Println(a)
		data, _ := json.Marshal(a)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
	http.HandleFunc("/get_data", h3)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":"+"3001", nil)
}

/*
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣀⣙⣆⠀⠈⢳⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣤⠶⠟⠛⠉⠁⠉⠛⠃⠀⠈⣿⠻⠷⠶⣦⣤⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣴⠟⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡄⠀⠀⠀⠀⠀⠈⠙⢿⣦⣄⣀⣤⣀⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⠾⠋⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣇⠀⠀⠀⠀⠀⠀⠀⠈⢿⣿⡉⢹⣿⣿⣿⣷⣶⣶⣤⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⡾⠋⠀⣼⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢻⡄⠀⠀⠀⠀⠀⠀⠀⠘⣿⣧⠀⢩⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⣴⠏⠀⠀⢸⡏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⣧⠀⠀⠀⠀⠀⠀⢀⠀⠘⣿⣷⡀⢉⣿⣿⣿⣿⡏⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⣠⣾⠃⠀⠀⠀⣿⠃⠀⠀⠀⠀⠀⠀⠀⠁⠀⠀⣿⠀⠀⠀⢹⣇⠀⠀⠀⠀⠀⠘⣇⠀⠘⢿⣷⡉⠉⣿⣿⣿⠁⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⣀⣠⣴⣾⡿⠁⠀⠀⠀⠀⣿⠀⠀⠀⠀⠀⠀⠀⢰⡇⠀⠀⠸⣇⠀⠀⠀⢻⡄⠀⠀⠀⠀⠀⢻⡀⠀⠈⠻⣿⣮⡉⢹⡏⠀⠀⠀⠀⠀⠀⠀⠀⠀
⢰⣿⣿⣿⣿⣿⠁⠀⢀⠀⠁⠀⢻⡆⠀⠀⠀⠀⠀⠀⢸⣧⠀⠀⠀⢻⡄⠀⠀⠀⢿⡀⠀⠀⠀⠀⠸⡇⠀⠀⠀⠘⡿⣿⣿⣷⡀⠀⠀⠀⠀⠀⠀⠀⠀
⠘⣿⣿⣿⣿⠇⠀⠀⣾⠀⠀⠀⢸⣧⠀⠀⠀⠀⠀⠀⠈⣿⣦⠀⠀⠈⢿⣄⠀⠀⠈⢷⡀⠀⠀⠀⠀⣷⠀⠀⠀⠀⢷⡀⠙⢿⣷⡀⠀⠀⠀⠀⠀⠀⠀
⠀⢻⣿⣿⡏⠀⢠⡀⢻⠀⠀⠀⢸⣿⣦⡀⠀⠀⠀⠀⠀⢿⡉⢷⡄⠀⠘⢿⣦⡀⠀⠈⢷⡀⠀⠀⠀⢻⠀⠀⠀⠀⠈⣧⠀⠈⢻⣷⡀⠀⠀⠀⠀⠀⠀
⠀⠘⣿⡿⠀⠀⣸⠀⣸⡇⠀⠀⢸⡇⠈⢷⣄⡀⠀⠀⠀⢺⣇⠀⠙⢦⣄⠈⢷⡹⢦⡀⠈⣷⠀⠀⠀⢸⡇⠀⠀⠀⠀⠸⣇⠀⠀⠹⣷⡀⠀⠀⠀⠀⠀
⠀⠀⣸⡇⠀⠀⡯⢠⣿⢿⡄⠀⢸⡇⠀⠀⠈⠛⠶⣦⣄⣀⣹⣿⡓⠳⠎⠛⠲⠿⢦⣽⣶⣼⣇⠀⠀⢸⡇⠀⠀⠀⠀⠀⢻⡄⠀⠀⢻⣧⠀⠀⠀⠀⠀
⠀⢠⣿⠀⠀⠀⡇⣼⠏⠀⠻⣆⢘⣧⣴⠖⠋⠀⠀⠀⠀⠉⠁⠉⠛⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⡁⠀⢸⡇⠀⠀⠀⠀⠀⠘⣷⠀⠀⠈⣿⣇⠀⠀⠀⠀
⠀⣼⡟⠀⠀⠀⣿⡟⠀⠀⠀⠙⠳⠥⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣠⣤⣤⣼⡇⠀⢸⡇⠀⠀⠀⠀⠀⠀⢹⡇⠀⠀⡟⢿⣆⠀⠀⠀
⢀⣿⡇⠀⠀⠀⣿⠇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣠⡤⢶⣾⣿⣿⣏⡹⠿⣇⠀⢸⡇⠀⠀⠀⠀⠀⠀⠘⣧⠀⠀⣧⠸⣿⡀⠀⠀
⢸⣿⢣⠀⠀⠀⣿⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣶⣯⠵⠶⠛⠉⠁⠀⠀⠀⠀⢿⠀⢸⡇⠀⠀⠀⠀⠀⠀⠀⢻⠀⠀⣿⠀⣿⣧⠀⠀
⣸⡏⢹⠀⠀⠀⢿⡇⠀⠀⠀⣠⣤⣶⣾⣿⣻⣿⡿⠖⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⠀⢀⠀⣀⠀⢸⡀⢸⠀⠀⠀⠀⠀⠀⠀⠀⢸⡆⠀⡿⢰⡏⣿⡀⠀
⣿⡇⢸⡄⠀⠀⢸⣿⢀⣴⣟⣡⡽⠟⠛⠋⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⡆⠸⣗⠻⠗⠻⠇⢸⡇⣸⠁⠀⠀⠀⠀⠀⠀⠀⢸⡇⠀⡇⢸⠇⢸⣧⠀
⣿⡅⠘⣇⠀⠀⠀⣿⡘⠛⠉⠁⠀⠀⠀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠁⠀⠀⠀⠀⠀⠀⢈⡇⣿⠀⠀⠀⠀⠀⠀⠀⠀⢸⡇⢠⣃⡿⠀⠀⣿⠀
⣿⢷⡀⢹⡄⠀⠀⢹⡇⠀⠀⣸⡆⠶⠄⠛⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣠⡶⠃⠀⠀⠀⠀⠀⢸⣧⡇⠀⠀⠀⠀⠀⠀⠀⠀⢸⡇⠘⣼⠃⠀⠀⢻⡄
⣿⠈⣧⠈⢷⠀⠀⠈⣿⠀⠀⠈⠀⠀⠀⠀⢀⠀⠀⢀⣀⣤⠴⠖⢚⣩⠽⠋⠀⠀⠀⠀⠀⠀⠀⠀⣿⠁⠀⠀⠀⠀⠀⠀⠀⠀⣸⠀⣰⠏⠀⠀⠀⢼⡇
⣿⠀⠘⣇⠘⣧⠀⠀⢸⡇⠀⠀⠀⠀⠀⠀⠉⠛⠛⠷⠖⠒⠒⠛⠉⠀⠀⠀⠀⠀⠀⠀⠀⣀⣴⢁⡏⠀⠀⠀⠀⠀⠀⠀⠀⠀⡿⢰⠟⠀⠀⠀⠀⣿⡄
⣿⡄⠀⠘⣦⠘⣇⠀⠈⣿⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣠⣴⠾⠋⢀⣽⡇⠀⠀⠀⠀⠀⠀⠀⠀⢨⡷⠋⠀⠀⠀⠀⠀⣿⠀
⢸⣧⠀⠀⠘⢧⡘⢧⡀⠘⠻⠶⢤⣤⣀⣀⣀⡀⠀⠀⠀⠀⠀⠀⣀⣀⣠⣴⣾⠟⠋⢀⣠⠶⢻⡏⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⠃⠀⠀⠀⠀⠀⣰⡏⠀
⠀⢿⡆⠀⠀⡈⢳⣄⠱⣄⠀⠀⠀⠀⠀⣽⠉⠉⢉⣉⠙⢿⣉⠉⠻⣿⡿⠋⢀⣠⠖⠋⠁⠀⣾⠁⠀⠀⠀⠀⠀⠀⠀⠀⢰⡟⠀⠀⠀⠀⠀⣠⡟⠀⠀
⠀⠈⢿⡄⠐⣧⠀⠙⢦⡈⠀⠀⠀⠀⠀⢻⣆⠀⠀⠙⢦⣀⠉⠳⢤⣘⣧⠶⠋⠁⠀⠀⠀⣰⡿⠀⠀⠀⠀⠀⠀⠀⢀⣠⠿⠃⠀⠀⠀⢀⣴⠟⠁⠀⠀
⠀⠀⠈⢿⣄⢿⣧⡀⠀⠛⢦⣄⠀⠀⠀⢸⣿⣷⣄⡀⠀⠉⠳⠶⣶⠞⠁⠀⠀⠀⢀⣠⣾⣿⠀⠀⠀⠀⠀⠀⣠⣴⡏⠁⠀⠀⢀⣠⡴⠟⠁⠀⠀⠀⠀
⠀⠀⠀⠀⠻⣾⣿⣛⣦⣄⠀⠈⠛⠲⠦⣄⣿⡇⠈⠙⠛⠶⠶⢶⣿⠀⠀⠀⢀⣴⣿⣿⣿⣯⣀⣀⣤⣤⣶⣿⣿⣿⣿⡛⠛⠋⠉⠉⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠘⣻⣷⣄⣙⡛⠶⠦⣤⣤⣄⣸⣷⡄⠀⠀⠀⢠⠏⣻⠀⣠⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠛⣿⣿⣿⣿⣷
*/
