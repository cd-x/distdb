package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/cd-x/distdb/db"
	"github.com/cd-x/distdb/web"
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

	srv := web.NewServer(db)

	http.HandleFunc("/get", srv.GetHandler)
	http.HandleFunc("/set", srv.SetHandler)

	log.Fatal(http.ListenAndServe(*http_location, nil))
}
