package httpd

import (
	"database/sql"
	"net/http"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

type Post struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Author string `json:"author"`
}

func GetAll(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var posts []Post
		result, err := db.Query("SELECT * FROM goshit")
		if err != nil {
			log.Fatal(err)
		}
		for result.Next() {
			var post Post
			err := result.Scan(&post.ID, &post.Title, &post.Body, &post.Author)
			if err != nil {
				log.Fatal(err)
			}
			posts = append(posts, post)
		}
		json.NewEncoder(w).Encode(posts)
	}
}
func GetOne(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)

		result, err := db.Query("SELECT * FROM goshit WHERE id = ?", params["id"])
		if err != nil {
			log.Fatal(err)
		}
		var post Post
		for result.Next() {
			err := result.Scan(&post.ID,
				&post.Title,
				&post.Body,
				&post.Author)
			if err != nil {
				log.Fatal(err)
			}
		}
		if post.ID != 0 {
			json.NewEncoder(w).Encode(post)
		} else {
			fmt.Fprintf(w, "{}")
		}
	}
}
func CreateOne(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stmt, err := db.Prepare("INSERT INTO goshit(title, body, author) VALUES(?,?,?)")
		if err != nil {
			log.Fatal(err)
		}
		bod, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		var post Post
		json.Unmarshal(bod, &post)
		title := post.Title
		body := post.Body
		author := post.Author
		_, err = stmt.Exec(title, body, author)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "New post was created")
	}
}
func UpdateOne(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get params value
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)

		// prepare statement
		stmt, err := db.Prepare("UPDATE goshit SET title = ?, body = ?, author = ? WHERE id = ?")
		if err != nil {
			log.Fatal(err)
		}
		// extracting values from body
		bdy, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		var post Post

		json.Unmarshal(bdy, &post)
		title := post.Title
		body := post.Body
		author := post.Author

		// execute statement
		_, err = stmt.Exec(title, body, author, params["id"])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Post Updated")
	}
}
func DeleteOne(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		stmt, err := db.Prepare("DELETE FROM goshit WHERE id = ?")
		if err != nil {
			log.Fatal(err)
		}
		_, err = stmt.Exec(params["id"])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Post deleted")
	}
}
