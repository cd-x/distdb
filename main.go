package main

import (
	"flag"
	"log"
	"net/http"
	"slices"

	"github.com/BurntSushi/toml"
	"github.com/cd-x/distdb/config"
	"github.com/cd-x/distdb/db"
	"github.com/cd-x/distdb/web"
)

var (
	db_location   = flag.String("db_location", "", "Database location")
	http_location = flag.String("http_address", "localhost:8080", "Address to connect to database")
	configFile    = flag.String("config-file", "", "Static sharding config file")
	shardName     = flag.String("shard", "", "name of the shard")
)

func parseFlags() {
	flag.Parse()
	if *db_location == "" {
		log.Fatalf("missing db_location")
	}
}

func main() {

	parseFlags()

	var c config.Config

	if _, err := toml.DecodeFile(*configFile, &c); err != nil {
		log.Fatalf("toml.DecodeFile(%q): %v", *configFile, err)
	}
	// find shard index
	idx := slices.IndexFunc(c.Shards, func(s config.Shard) bool { return s.Name == *shardName })
	shardCount := len(c.Shards)
	if idx == -1 || shardCount < 1 {
		log.Fatalf("shard %v not foung", *shardName)
	}

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
