package main

import (
	"config"
	"database/sql"
	"fmt"
	"httpd"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

func main() {
	fmt.Println("Starting App...")
	// createing Database Comection
	db, err := sql.Open("mysql", config.MysqlURI())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Printf("Database Connected...")
	// defining route
	r := mux.NewRouter()

	r.HandleFunc("/posts", httpd.GetAll(db)).Methods("GET")
	r.HandleFunc("/post/{id}", httpd.GetOne(db)).Methods("GET")
	r.HandleFunc("/post", httpd.CreateOne(db)).Methods("POST")
	r.HandleFunc("/post/{id}", httpd.UpdateOne(db)).Methods("PUT")
	r.HandleFunc("/post/{id}", httpd.DeleteOne(db)).Methods("DELETE")

	http.ListenAndServe(":8000", r)
	fmt.Printf("App Stoped")
}
