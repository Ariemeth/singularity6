package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var cache ttlcache.SimpleCache = ttlcache.NewCache()

func main() {
	log.Println("Initializing service")

	rand.Seed(time.Now().UnixNano())
	// set TTL on shorteners to 24 hours
	cache.SetTTL(time.Duration(24 * time.Hour))

	r := mux.NewRouter()
	r.HandleFunc("/", shortenHandler).Methods("POST")
	r.HandleFunc("/{shortendURL}", redirectHandler).Methods("GET")
	r.HandleFunc("/healthz", healthHandler).Methods("GET")

	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    ":9000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("healthy and running"))
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	Shortened string `json:"short_url_code"`
	URL string `json:"url"`
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var sr shortenRequest
	if err := decoder.Decode(&sr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request: body does not contain url"))
	}
	key := createKey()
	err := cache.Set(key, sr.URL)
	if err != nil {
		log.Panicln(err)
	}

	resp:= shortenResponse{Shortened: key, URL: sr.URL}
	json.NewEncoder(w).Encode(resp)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortend := vars["shortendURL"]
	if shortend == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request: no shortened url present"))
		return
	}

	log.Println(shortend)

	fullURL, err := cache.Get(shortend)
	if err == ttlcache.ErrNotFound || fullURL == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Bad request: no url found for this shortner %s", fullURL)))
		return
	}

	log.Printf("Redirecting to %s\n", fullURL)
	http.Redirect(w, r, fmt.Sprintf("%s", fullURL), http.StatusMovedPermanently)
}

func createKey() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}
