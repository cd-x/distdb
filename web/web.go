package web

import (
	"fmt"
	"net/http"

	"github.com/cd-x/distdb/db"
)

type Server struct {
	db *db.Database
}

func NewServer(db *db.Database) *Server {
	return &Server{
		db: db,
	}
}

func (server *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value, err := server.db.GetKey(key)
	fmt.Fprintf(w, "Value=%q, error = %v", value, err)
}

func (server *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")
	err := server.db.SetKey(key, []byte(value))
	fmt.Fprintf(w, "Error = %v", err)
}
