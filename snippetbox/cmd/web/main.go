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
	addr := flag.String("addr", ":4000", "HTTP network port")
	flag.Parse()

	f, err := os.OpenFile("./tmp/info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal(err)
	}

	infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime) // here | operator is a bitwise operator

	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errLog, infoLog,
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
