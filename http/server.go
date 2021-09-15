package http

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/mux"

	"net/http"
)

type Server struct {
	Interaction Interaction

	Server *http.Server
	Log    *log.Logger
}

type Interaction interface {
	Get(key string) (value []byte, err error)
	Set(key []byte, value []byte) error
	GetMetrics() interface{}
	Ready(cxt context.Context) bool
}

func New(interaction Interaction, address string) *Server {
	context := Server{
		Interaction: interaction,
	}
	m := mux.NewRouter()
	m.Methods("GET")
	m.HandleFunc("/key/{key}", context.GetByKey)
	m.HandleFunc("/metrics", context.GetMetrics)

	m.HandleFunc("/health/ready", context.Ready)
	m.HandleFunc("/health/live", context.Live)

	srv := &http.Server{
		Handler:      m,
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	context.Server = srv
	return &context
}

func (c *Server) GetByKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	key := vars["key"]
	result, err := c.Interaction.Get(key)

	if err == nil && len(result) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		//log.Error(r.Context(), "Error retrieving from cache", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(result)
}

func (c *Server) Live(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (c *Server) Ready(w http.ResponseWriter, r *http.Request) {
	ready := c.Interaction.Ready(r.Context())

	if !ready {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Server) GetMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := c.Interaction.GetMetrics()

	responseBytes, _ := json.Marshal(metrics)

	w.Header().Set("Content-Type", "application/json")

	_, _ = w.Write(responseBytes)
	w.WriteHeader(http.StatusOK)
}

func (c *Server) Start(ctx context.Context) error {
	return c.Server.ListenAndServe()
}
