package config

import (
	"log"
	"os"
	"slices"

	"github.com/BurntSushi/toml"
)

type Shard struct {
	Name    string `toml:"name"`
	Id      int    `toml:"id"`
	Address string `toml:address`
}

type Config struct {
	Shards []Shard `toml:"shard"`
}

func GetConfigs(data *Config, configFile string) error {
	content, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("os.ReadFile(%q): %v", configFile, err)
	}
	if err := toml.Unmarshal([]byte(content), &data); err != nil {
		log.Fatalf("toml.Unmarshal(%q):%v", configFile, err)
		return err
	}
	log.Printf("ConfigFile(%q): %v\n", configFile, data.Shards)
	return nil
}

func ParseShards(data *Config, shardName string) (int, int, map[int]string) {
	// find shard index
	idx := slices.IndexFunc(data.Shards, func(s Shard) bool { return s.Name == shardName })
	shardMap := make(map[int]string)
	for _, shard := range data.Shards {
		shardMap[shard.Id] = shard.Address
	}
	shardCount := len(data.Shards)
	if idx == -1 || shardCount < 1 {
		log.Fatalf("shard %v not found", shardName)
	}
	return idx, shardCount, shardMap
}
