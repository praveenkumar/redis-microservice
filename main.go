package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	host   = getEnv("DB_HOST", "localhost")
	port   = getEnv("DB_PORT", "6379")
	dbname = getEnv("REDIS_DATABASE", "0")
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func keys(w http.ResponseWriter, r *http.Request, db *redis.Client) {
	// Use the Redis database to get a value
	keys, err := db.Keys(r.Context(), "*").Result()
	if err != nil {
		fmt.Println("Error getting keys from Redis:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write the value to the response
	fmt.Fprintf(w, "Value from Redis: %s", keys)
}

func value(w http.ResponseWriter, r *http.Request, db *redis.Client) {
	// Get the key parameter from the query string
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing key parameter", http.StatusBadRequest)
		return
	}

	// Use the Redis database to get the value for the key
	val, err := db.Get(r.Context(), key).Result()
	if err == redis.Nil {
		fmt.Fprintf(w, "Key %s not found", key)
		return
	} else if err != nil {
		fmt.Println("Error getting value from Redis:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write the value to the response
	fmt.Fprintf(w, "Value for key %s: %s", key, val)
}

func main() {
	var ctx = context.Background()

	dbIndex, err := strconv.Atoi(dbname)
	if err != nil {
		log.Fatalln(err)
	}

	// Connect to Redis database
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "",      // no password set
		DB:       dbIndex, // use default DB
	})

	// Ping Redis to check connection
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalln("Error connecting to Redis:", err)
	}
	log.Println("Connected to Redis:", pong)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		keys(w, r, client)
	})
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		value(w, r, client)
	})

	err = http.ListenAndServe(":8080", nil)
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("server closed")
	} else if err != nil {
		log.Fatalf("error starting server: %s\n", err)
	}
}
