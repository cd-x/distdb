package web

import (
	"fmt"
	"io"
	"net/http"

	"github.com/cd-x/distdb/db"
	"github.com/cd-x/distdb/utils"
)

type Server struct {
	db         *db.Database
	shardIdx   int
	shardCount int
	shardMap   map[int]string
}

func NewServer(db *db.Database, shardIdx int, shardCount int, addressMap map[int]string) *Server {
	return &Server{
		db:         db,
		shardIdx:   shardIdx,
		shardCount: shardCount,
		shardMap:   addressMap,
	}
}

func (server *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	shardIdx := server.getShardIndex(key)
	if shardIdx != server.shardIdx {
		server.route(w, r, shardIdx)
		return
	}
	value, err := server.db.GetKey(key)
	fmt.Fprintf(w, "Value=%q, error = %v\n", value, err)
}

func (server *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")
	// redirection
	shardIdx := server.getShardIndex(key)
	if shardIdx != server.shardIdx {
		server.route(w, r, shardIdx)
		return
	}
	err := server.db.SetKey(key, []byte(value))
	fmt.Fprintf(w, "Error = %v, shardIdx = %d \n", err, shardIdx)
}

func (server *Server) getShardIndex(key string) int {
	hashVal := utils.GetKeyHash(key)
	shardIdx := int(hashVal % uint64(server.shardCount))
	return shardIdx
}

func (server *Server) route(w http.ResponseWriter, r *http.Request, shardIdx int) {
	routedAddress := "http://" + server.shardMap[shardIdx] + r.RequestURI
	fmt.Fprintf(w, "redirecting to shard = %d from shard = %d, address = (%q)\n", shardIdx, server.shardIdx, routedAddress)
	resp, err := http.Get(routedAddress)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error occured redirecting the request: %v\n", err)
		return
	}
	io.Copy(w, resp.Body)
}
