package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/cd-x/distdb/config"
	"github.com/cd-x/distdb/db"
	"github.com/cd-x/distdb/web"
)

var (
	db_location   = flag.String("db_location", "", "Database location")
	http_location = flag.String("http-address", "localhost:8080", "Address to connect to database")
	configFile    = flag.String("config-file", "", "Static sharding config file")
	shardName     = flag.String("shard", "", "name of the shard")
)

func parseFlags() {
	flag.Parse()
	if *db_location == "" {
		log.Fatalf("missing db_location")
	}
	if *configFile == "" {
		log.Fatalf("Required config.toml file not provided.")
	}
}

func main() {

	parseFlags()
	// Unmarshall shard configs
	var data config.Config
	config.GetConfigs(&data, *configFile)
	idx, shardCount, shardMap := config.ParseShards(&data, *shardName)

	db, close, err := db.NewDatabase(*db_location)

	if err != nil {
		log.Fatalf("NewDatabase[%q]: %v", *db_location, err)
	}
	defer close()

	srv := web.NewServer(db, idx, shardCount, shardMap)

	http.HandleFunc("/get", srv.GetHandler)
	http.HandleFunc("/set", srv.SetHandler)

	log.Fatal(http.ListenAndServe(*http_location, nil))
}
