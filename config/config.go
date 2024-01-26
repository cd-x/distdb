package config

type Shard struct {
	Name string
	Id   uint64
}

type Config struct {
	Shards []Shard
}
