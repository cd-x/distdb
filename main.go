package main

import (
	"flag"
	"log"
	"net/http"
	"os"
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
	if *configFile == "" {
		log.Fatalf("Required config.toml file not provided.")
	}
}

func main() {

	parseFlags()

	content, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("os.ReadFile(%q): %v", *configFile, err)
	}
	var data config.Config
	if err := toml.Unmarshal([]byte(content), &data); err != nil {
		log.Fatalf("toml.Unmarshal(%q):%v", *configFile, err)
	}
	log.Printf("ConfigFile(%q): %v\n", *configFile, data.Shards)

	// find shard index
	idx := slices.IndexFunc(data.Shards, func(s config.Shard) bool { return s.Name == *shardName })
	shardCount := len(data.Shards)
	if idx == -1 || shardCount < 1 {
		log.Fatalf("shard %v not found", *shardName)
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
