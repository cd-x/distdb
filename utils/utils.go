package utils

import "hash/fnv"

func GetKeyHash(key string) uint64 {
	hash := fnv.New64()
	hash.Write([]byte(key))
	return hash.Sum64()
}
