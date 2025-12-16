package server

import (
	"fmt"
	"time"
	"log"
	"net/http"
	"strings"
	"context"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	upg = websocket.Upgrader {
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true //OVO TREBA DA SE IZMENI!!!!
		},
	}
	aktivniIgraci = make(map[string]*Igrac)
	aktivniIgraciMux sync.RWMutex
)

type Server struct {
	port int
}

func (server *Server) HandlerInit() http.Handler {
	mux := http.NewServeMux()
	redisDB := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		Password: "",
		DB: 0,
	})
	ctx := context.Background()
	fs := http.FileServer(http.Dir("static"))

	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/static")
		r.URL.Path = path
		fs.ServeHTTP(w, r)
	})

	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWS(w, r, redisDB, ctx)
	})

	return mux
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	log.Print(r.URL)

	if r.URL.Path != "/" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "index.html")
}

func handleWS(w http.ResponseWriter, r *http.Request, rdb *redis.Client, ctx context.Context) {
	conn, err := upg.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Handler upgrader greška: %v\n", err)
		return
	}

	igrac := NoviIgrac(conn)

	aktivniIgraciMux.Lock()
	aktivniIgraci[igrac.UUID] = igrac
	aktivniIgraciMux.Unlock()

	go igrac.CitajWSPoruke(ctx, rdb)
}

func NoviServer() (*http.Server, int) {
	noviServer := Server {
		port: 8080,
	}

	return &http.Server {
		Addr: fmt.Sprintf(":%d", noviServer.port),
		Handler: noviServer.HandlerInit(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}, noviServer.port
}

func NadjiAktivnogIgraca(uuid string) *Igrac {
	aktivniIgraciMux.RLock()
	defer aktivniIgraciMux.RUnlock()
	return aktivniIgraci[uuid]
}

func DiskonektujIgraca(igracUUID string) {
	aktivniIgraciMux.Lock()
	delete(aktivniIgraci, igracUUID)
	aktivniIgraciMux.Unlock()
	log.Printf("Igrac sa uuid %v diskonektovan\n", igracUUID)
}
