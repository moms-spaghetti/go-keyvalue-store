package protocols

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	_ "net/http/pprof"
	"task1/internal/logger"
	"task1/internal/metrics"
	"task1/internal/store"
	"time"
)

const (
	httpaddr    = ":8080"
	httptimeout = 5 * time.Second
)

type HTTPServer struct {
	http    *http.Server
	logger  *logger.Logger
	storage *store.Storage
	metrics *metrics.Metrics
}

func NewHTTP(
	logger *logger.Logger,
	storage *store.Storage,
	metrics *metrics.Metrics,
) *HTTPServer {

	return &HTTPServer{
		http: &http.Server{
			Addr: httpaddr,
		},
		logger:  logger,
		storage: storage,
		metrics: metrics,
	}
}

func (hs HTTPServer) Start() {
	http.HandleFunc("/", hs.rootHandler)

	go func() {
		log.Printf("http listning on %s", httpaddr)
		if err := hs.http.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}
	}()
}

func (hs HTTPServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), httptimeout)
	defer cancel()

	if err := hs.http.Shutdown(ctx); err != nil {
		panic(err)
	}
	log.Print("HTTP shutdown ok")
}

func (hs *HTTPServer) rootHandler(w http.ResponseWriter, r *http.Request) {
	var (
		req       jsonRequest
		err       error
		storeData interface{}
	)

	if err == nil {
		err = json.NewDecoder(r.Body).Decode(&req)
	}

	if err == nil {
		hs.metrics.LogMetrics(r.Method)
		switch r.Method {
		case http.MethodGet:
			hs.logger.Log("HTTP GET request")
			storeData, err = hs.storage.Get(req.Query)
		case http.MethodPost:
			hs.logger.Log("HTTP POST request")
			err = hs.storage.Post(req.Payload)
		case http.MethodDelete:
			hs.logger.Log("HTTP DELETE request")
			err = hs.storage.Delete(req.Query)
		default:
			err = ErrRouteForbidden
		}
	}

	status, out := BuildJsonResponse(err, storeData, hs.logger)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(out)
}
