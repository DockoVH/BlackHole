package main

import (
	"log"
	"context"

	"BlackHole/internal/server"

	"github.com/redis/go-redis/v9"
	"github.com/apache/cassandra-gocql-driver/v2"
)

func main() {
	redisDB := redis.NewClient(&redis.Options {
		Addr: "redis:6379",
		Password: "",
		DB: 0,
	})

	redisDB.Del(context.Background(), "sve-sobe")

	cluster := gocql.NewCluster("cassandra:9042")
	cluster.Keyspace = "blackhole"
	cassandraSession, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Greška prilikom pravljenja cassandra sesije: ", err)
	}

	defer func() {
		redisDB.Close()
		cassandraSession.Close()
	}()

	server, port := server.NoviServer(cassandraSession)
	log.Printf("Server osluškuje na adresi: %v\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe greška: ", err)
	}
}
