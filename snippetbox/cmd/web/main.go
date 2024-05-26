package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"isurucuma.golearning.letsgo/internal/models"

	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql" // this is aliazed with _ because there is no place in this file we use the driver specific thing
	// but still we want to run the init method in the driver so that it can register itself with the databse/sql package
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

const (
	username = "root"
	password = "password"
	hostname = "localhost:3306"
	dbName   = "snippetbox"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network port")
	defaultDsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, hostname, dbName)

	dsn := flag.String("dsn", defaultDsn, "MySQL data source name")

	flag.Parse()

	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	f, err := os.OpenFile("./tmp/info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		errLog.Fatal(err)
	}
	infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime) // here | operator is a bitwise operator

	db, err := openDB(*dsn)
	if err != nil {
		errLog.Fatal(err)
	}
	infoLog.Print("Database connection successfull")
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	app := &application{
		errorLog:      errLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}
	infoLog.Print("Server starting on", *addr)

	server := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errLog,
	}
	err = server.ListenAndServe()
	errLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
