package main

import (
	"flag"
	"log"

	"github.com/boltdb/bolt"
)

var (
	db_location = flag.String("db_location", "", "Database location")
)

func parseFlags() {
	flag.Parse()
	if *db_location == "" {
		log.Fatalf("missing db_location")
	}
}

func main() {

	parseFlags()

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
