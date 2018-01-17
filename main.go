package main

import (
	"database/sql"
	"fmt"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
)

type Book struct {
	isbn   string
	title  string
	author string
	price  float32
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://qndpukrm:OZr1ZmDFnatdCEdnnUrkw4ZQf7Fr3iDN@baasu.db.elephantsql.com:5432/qndpukrm")
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := httprouter.New()

	r.GET("/", homeIndex)

	r.GET("/books", booksIndex)
	r.POST("/books", booksCreate)

	// Posts singular
	r.GET("/books/:isbn", booksShow)

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}

func homeIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	resp := `
		<h1> Home </h1>
		<a href="/books">See index of Books</a>
	`
	fmt.Fprintln(w, resp)
}

func booksCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	isbn := r.FormValue("isbn")
	title := r.FormValue("title")
	author := r.FormValue("author")

	if isbn == "" || title == "" || author == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	price, err := strconv.ParseFloat(r.FormValue("price"), 32)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	createQuery := `INSERT INTO books
				   VALUES($1, $2, $3, $4)
	
	`

	result, err := db.Exec(createQuery, isbn, title, author, price)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Fprintf(w, "Book %s created successfully (%d row affected)\n", isbn, rowsAffected)
}

func booksShow(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	isbn := p.ByName("isbn")

	if isbn == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	showQuery := `SELECT *
					FROM books
					WHERE isbn = $1
	`

	row := db.QueryRow(showQuery, isbn)

	bk := new(Book)
	err := row.Scan(&bk.isbn, &bk.title, &bk.author, &bk.price)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Fprintf(w, "%s, %s, %s, £%.2f\n", bk.isbn, bk.title, bk.author, bk.price)

}

func booksIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	indexQuery := `SELECT *
					FROM books
	`

	rows, err := db.Query(indexQuery)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	bks := make([]*Book, 0)

	for rows.Next() {
		bk := new(Book)
		err := rows.Scan(&bk.isbn, &bk.title, &bk.author, &bk.price)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		bks = append(bks, bk)

	}

	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	for _, bk := range bks {
		fmt.Fprintf(w, "%s, %s, %s, £%.2f\n", bk.isbn, bk.title, bk.author, bk.price)
	}
}
