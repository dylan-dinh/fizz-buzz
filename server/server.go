package server

import (
	rds "github.com/dylan-dinh/fizz-buzz/redis"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"net/http"
)

type AppServer struct {
	Conn *redis.Conn
}

func StartServer(handler http.Handler) {
	server := &http.Server{
		Addr:      "127.0.0.1:2000",
		Handler:   handler,
		TLSConfig: nil,
	}

	// always return non nil error
	server.ListenAndServe()

}

func GetRouter(rds *rds.RedisService) http.Handler {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/fizz_buzz", http.HandlerFunc(rds.Simple).ServeHTTP).Methods(http.MethodGet)

	return r
}
