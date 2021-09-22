package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/mux"
	"github.com/reactivex/rxgo/v2"
	"github.com/stevehebert/k2edge/internal/persistence"
	"github.com/stevehebert/k2edge/observable"

	"net/http"
)

type Server struct {
	StorageFacade persistence.StorageFacade

	Server *http.Server
	Log    *log.Logger
}

func mapToStorage(storageFacade persistence.StorageFacade) func(context.Context, interface{}) (interface{}, error) {
	return func(_ context.Context, item interface{}) (interface{}, error) {
		msg, convert := item.(*kafka.Message)
		if !convert {
			panic("bad read from kafka")
		}
		storageFacade.Set(msg.Key, msg.Value)
		return item, nil
	}
}

func New(storageFacade persistence.StorageFacade, address string) *Server {
	c, err := observable.KafkaConnect(&kafka.ConfigMap{
		"bootstrap.servers": "127.0.0.1:29092",
		"group.id":          "aa",
		"auto.offset.reset": "earliest",
	}, []string{"incoming"})

	if err != nil {
		return nil
	}

	go func(consumer *kafka.Consumer) {
		fmt.Println("starting")
		obs := rxgo.Defer([]rxgo.Producer{observable.KafkaObservable(consumer)}).
			Map(mapToStorage(storageFacade))

		for o := range obs.Observe() {
			if o.Error() {
				fmt.Print("e")
			}
		}
	}(c)

	context := Server{
		StorageFacade: storageFacade,
	}

	m := mux.NewRouter()
	//m.Methods("GET")
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
	result, err := c.StorageFacade.Get(key)

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
	ready := c.StorageFacade.Ready(r.Context())

	if !ready {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Server) GetMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := c.StorageFacade.GetMetrics()

	responseBytes, _ := json.Marshal(metrics)

	w.Header().Set("Content-Type", "application/json")

	_, _ = w.Write(responseBytes)
	w.WriteHeader(http.StatusOK)
}

func (c *Server) Start(ctx context.Context) error {
	//fmt.Println("listening")
	return http.ListenAndServe(":8080", c.Server.Handler)
	//return c.Server.ListenAndServe()

	/*r := mux.NewRouter()

	r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
	})

	return http.ListenAndServe(":80", r)*/
}
