package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	var app *application = &application{
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	infoLog.Printf("Starting server on %s", *addr)

	srv := &http.Server{
		Addr:     *addr,
		Handler:  mux,
		ErrorLog: errorLog,
	}

	err := srv.ListenAndServe()
	errorLog.Fatal(err)

}
