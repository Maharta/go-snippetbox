package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"snippetbox.maharta.dev/internal/models"

	_ "github.com/go-sql-driver/mysql" // Because this package isn't used in this code directly, but we only need for the driver to call
	// its init() function to initialize itself to the database/sql package.
)

type application struct {
	errorLog     *log.Logger
	infoLog      *log.Logger
	snippetModel *models.SnippetModel
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:gopass@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDBPool(*dsn)

	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	var app *application = &application{
		infoLog:      infoLog,
		errorLog:     errorLog,
		snippetModel: &models.SnippetModel{DB: db},
	}

	infoLog.Printf("Starting server on %s", *addr)

	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}

	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}

func openDBPool(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil

}
