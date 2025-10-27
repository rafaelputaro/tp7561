package main

import (
	"context"
	"log"
	"net/http"
	"time"
	"tp/common"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const MSG_ERROR_CLOSING_SERVER = "error closing server: %v"
const MSG_SERVER_ERROR = "server error: %v"
const MSG_CLOSING_SERVER = "closing server"
const MSG_SERVER_STARTED = "server started: %v"

// Aloja un buffer con las métricas del sistema las cuales puede ser recuperadas.
// Se interactúa con el mismo a través de una API-REST
type MetricsServer struct {
	config MetricsServerConfig
	server *http.Server
	gauge  prometheus.Gauge
}

func NewMetricsServer(config *MetricsServerConfig) *MetricsServer {
	// Crear router
	router := createRouter()
	// Crear metrics server
	server := MetricsServer{
		config: *config,
		server: &http.Server{
			Addr:    config.Url,
			Handler: router,
		},
		gauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "metrics",
			Name:      "metrics",
			Help:      "Metrics",
		}),
	}
	return &server
}

func (server *MetricsServer) Serve() {
	common.Log.Infof(MSG_SERVER_STARTED, server.config.Url)
	if err := server.server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf(MSG_SERVER_ERROR, err)
	}
}

func createRouter() *http.ServeMux {
	reg := prometheus.NewRegistry()
	router := http.NewServeMux()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	router.Handle("/metrics", promHandler)
	return router
}

// Cierra el servidor y sus recursos
func (server *MetricsServer) DisposeMetricsServer() {
	// Cerrar servidor
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.server.Shutdown(ctx); err != nil {
		common.Log.Fatalf(MSG_ERROR_CLOSING_SERVER, err)
		return
	}
	common.Log.Infof(MSG_CLOSING_SERVER)
}

/*

package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Book represents a book entity
type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// In-memory storage for books (for demonstration purposes)
var books []Book

func main() {
	// Initialize some sample data
	books = append(books, Book{ID: "1", Title: "The Hitchhiker's Guide to the Galaxy", Author: "Douglas Adams"})
	books = append(books, Book{ID: "2", Title: "1984", Author: "George Orwell"})

	// Create a new ServeMux
	router := http.NewServeMux()

	// Handle GET requests to /books
	router.HandleFunc("GET /books", getBooks)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", router))
}

// getBooks handles the GET /books request
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}


*/
