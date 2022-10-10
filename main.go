package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"html/template"
	"net/http"
	"strconv"
)

type User struct {
	Id   uint   `json:"id"`
	Name string `json:"FullName"`
	Age  uint   `json:"Age"`
}

var user_list = []User{}
var showUser = User{}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin12345"
	dbname   = "golangdb"
)

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)

	}
	user_list = []User{}
	user_var, err := db.Query(`SELECT * FROM users`)
	if err != nil {
		panic(err)
	}
	for user_var.Next() {

		var user User
		err = user_var.Scan(&user.Id, &user.Name, &user.Age)
		if err != nil {
			panic(err)
		}
		//fmt.Println(fmt.Sprintf("ID:%d and his name is %s", user.Id, user.Name))
		user_list = append(user_list, user)
	}
	defer user_var.Close()

	defer db.Close()

	t.ExecuteTemplate(w, "index", user_list)

}

func contact_list(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/contacts.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	t.ExecuteTemplate(w, "contact", nil)

}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header1.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	//fmt.Println(name)
	age := r.FormValue("age")
	//fmt.Println(age)
	i, err := strconv.Atoi(age)
	if err != nil {
		// ... handle error
		panic(err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)

	}

	defer db.Close()

	//insert, err := db.Query(fmt.Sprintf(`INSERT INTO users (name)  VALUES (%s)`, name))
	insert, err := db.Query(fmt.Sprintf("INSERT INTO users (name,age) VALUES ('%s',%d)", name, i))

	if err != nil {
		panic(err)

	}
	defer insert.Close()

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func get_user(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	t, err := template.ParseFiles("templates/show_user.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)

	}
	showu, err := db.Query(fmt.Sprintf("SELECT * FROM users WHERE id = '%s'", vars["id"]))

	if err != nil {
		panic(err)
	}

	user_show := User{}
	for showu.Next() {
		var user User
		err = showu.Scan(&user.Id, &user.Name, &user.Age)
		//fmt.Println(user.Id, user.Name, user.Age)
		if err != nil {
			panic(err)
		}
		user_show = user
	}
	defer db.Close()

	t.ExecuteTemplate(w, "user_show", user_show)

}

func delete_user(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)

	}

	del, err := db.Query(fmt.Sprintf("DELETE FROM users WHERE id = '%s'", vars["id"]))
	//fmt.Println(vars["id"])
	if err != nil {
		panic(err)
	}

	defer del.Close()

	defer db.Close()

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func delete_all(w http.ResponseWriter, r *http.Request) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)

	}

	del, err := db.Query(fmt.Sprintf("DELETE FROM users"))
	if err != nil {
		panic(err)
	}

	defer del.Close()

	defer db.Close()

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func HandleReq() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/contacts/", contact_list).Methods("GET")
	rtr.HandleFunc("/create/", create).Methods("GET")
	rtr.HandleFunc("/save_article/", save_article).Methods("POST")
	rtr.HandleFunc("/info/{id:[0-9]+}", get_user).Methods("GET")
	rtr.HandleFunc("/info/{id:[0-9]+}/delete/", delete_user).Methods("GET")
	rtr.HandleFunc("/delete/", delete_all).Methods("GET")

	http.Handle("/", rtr)

	http.ListenAndServe(":8080", nil)

}

func main() {
	HandleReq()

}
