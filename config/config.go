package config

type Shard struct {
	Name string `toml:"name"`
	Id   int    `toml:"id"`
}

type Config struct {
	Shards []Shard `toml:"shard"`
}

type Fuck struct {
	Wtf string
}
