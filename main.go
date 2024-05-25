package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// we need 3 things
// a handler
// a rounter
// a server

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello this is let's go"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "THis is the snippet with ID %d...", id)
}

func createSnippet(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		// w.WriteHeader(http.StatusMethodNotAllowed)
		// w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(201)
	w.Write([]byte("snippet created..."))
}

func main() {
	// use http.ServeMux to Multiplex the requests based on their request path and method
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", createSnippet)

	log.Print("Server starting on :4000")
	err := http.ListenAndServe(":4000", mux)

	log.Fatal(err)
}
