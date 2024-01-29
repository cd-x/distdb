package config

type Shard struct {
	Name    string `toml:"name"`
	Id      int    `toml:"id"`
	Address string `toml:address`
}

type Config struct {
	Shards []Shard `toml:"shard"`
}
