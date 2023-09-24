package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
)

var client = redis.NewClient(&redis.Options{
	Addr:     "redis:6379",
	Password: "",
	DB:       0,
})

func main() {
	ctx := context.Background()
	client.Set(ctx, "counter", 0, 0)
	client.Set(ctx, "inc", 0, 0)
	http.HandleFunc("/set", handleSet)
	http.HandleFunc("/get", handleGet)
	http.HandleFunc("/total", getTotal)
	http.HandleFunc("/inc", handleInc)
	log.Println("Listening on port 80")
	log.Panic(http.ListenAndServe(":80", nil))
}

func incCounter() {
	ctx := context.Background()
	client.Incr(ctx, "counter")
}

func inc() {
	ctx := context.Background()
	client.Incr(ctx, "inc")
}

func handleInc(w http.ResponseWriter, r *http.Request) {
	inc()
	val, err := getDB("inc")
	json.NewEncoder(w).Encode(Response{Result: val, Error: err.Error()})
}

func setDB(key string, value interface{}) error {
	ctx := context.Background()

	err := client.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func getDB(key string) (value string, err error) {
	ctx := context.Background()

	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func handleSet(w http.ResponseWriter, r *http.Request) {
	// r.ParseForm()
	incCounter()
	err := validateRequestParameters(r, []string{"key", "value"})
	if err != nil {
		json.NewEncoder(w).Encode(Response{Result: "", Error: err.Error()})
		return
	}
	key, value := r.URL.Query().Get("key"), r.URL.Query().Get("value")
	err = setDB(key, value)
	log.Println(key, value)
	if err != nil {
		json.NewEncoder(w).Encode(Response{Result: "", Error: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(Response{Result: "OK", Error: ""})
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	incCounter()
	err := validateRequestParameters(r, []string{"key"})
	if err != nil {
		json.NewEncoder(w).Encode(Response{Result: "", Error: err.Error()})
		return
	}
	key := r.URL.Query().Get("key")
	val, err := getDB(key)
	if err != nil {
		json.NewEncoder(w).Encode(Response{Result: "", Error: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(Response{Result: val, Error: ""})
}

func getTotal(w http.ResponseWriter, r *http.Request) {
	val, err := getDB("counter")
	if err != nil {
		json.NewEncoder(w).Encode(Response{Result: "", Error: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(Response{Result: val, Error: ""})
}

type Response struct {
	Result string
	Error  string
}

func validateRequestParameters(r *http.Request, keys []string) error {
	r.ParseForm()
	for _, key := range keys {
		if _, ok := r.Form[key]; !ok {
			return fmt.Errorf("Key `%s` required, but not provided!", key)
		}
	}
	return nil
}
