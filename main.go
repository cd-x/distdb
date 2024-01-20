package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/cd-x/distdb/db"
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

	db, close, err := db.NewDatabase(*db_location)

	if err != nil {
		log.Fatalf("NewDatabase[%q]: %v", *db_location, err)
	}
	defer close()
	getHandler := func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		key := r.Form.Get("key")
		value, err := db.GetKey(key)
		fmt.Fprintf(w, "Value=%q, error = %v", value, err)
	}

	setHandler := func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		key := r.Form.Get("key")
		value := r.Form.Get("value")
		err := db.SetKey(key, []byte(value))
		fmt.Fprintf(w, "Error = %v", err)
	}

	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/set", setHandler)

	log.Fatal(http.ListenAndServe(*http_location, nil))
}
