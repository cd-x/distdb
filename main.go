package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	db_location   = flag.String("db_location", "my.db", "Database location")
	http_location = flag.String("http_address", "localhost:8080", "Address to connect to database")
)

func parseFlags() {
	flag.Parse()
	if *db_location == "" {
		log.Fatalf("missing db_location")
	}
}

func main() {

	parseFlags()

	getHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Get database\n")
	}

	setHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Set database\n")
	}

	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/set", setHandler)

	log.Fatal(http.ListenAndServe(*http_location, nil))
}
