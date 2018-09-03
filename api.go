package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var playerRedis *redis.Client

type Player struct {
	Id  string
	Status	string
	Gameserver	string
}

func main() {
	corsObj := handlers.AllowedOrigins([]string{"*"})
	router := mux.NewRouter()
	router.HandleFunc("/player/{token}", PlayerHandler).Methods("GET")
	playerRedis = connectToRedis("redis-players:6379")
	log.Fatal(http.ListenAndServe(":6002", handlers.CORS(corsObj)(router)))
}

func connectToRedis(addr string) *redis.Client {
	var client *redis.Client
	for {
		client = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "",
			DB:       0,
		})
		_, err := client.Ping().Result()
		if err != nil {
			fmt.Println("Could not connect to redis")
			fmt.Println(err)
		} else {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("Connected to redis")
	return client
}

func PlayerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	status, err := playerRedis.HGet(token, "status").Result()
	if err!=nil || status == ""  {
		status = "invalid"
	}
	json.NewEncoder(w).Encode(status)
}
