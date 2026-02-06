package server

import (
	"fmt"
	"time"
	"log"
	"net/http"
	"strings"
	"context"
	"sync"
	"encoding/json"

	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/google/uuid"
	"github.com/apache/cassandra-gocql-driver/v2"
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
			return true
		},
	}
	aktivniIgraci = make(map[string]*Igrac)
	aktivniIgraciMux sync.RWMutex
	nepozeljneReci = make([]string, 0)
	znakoviInterpunkcije = []byte { ',', '.', ';', '!', '?', '{', '}', ':', '"' }
)

type Server struct {
	port int
}

func (server *Server) HandlerInit(cassandraSession *gocql.Session) http.Handler {
	mux := http.NewServeMux()
	redisDB := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		Password: "",
		DB: 0,
	})
	ctx := context.Background()
	fs := http.FileServer(http.Dir("static"))

	nepozeljneReci = pribaviNepozeljneReci(cassandraSession)

	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/static")
		r.URL.Path = path
		fs.ServeHTTP(w, r)
	})

	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWS(w, r, redisDB, ctx, cassandraSession)
	})
	mux.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		handleRegistrujIgraca(w, r, cassandraSession)
	})
	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		handlePrijaviIgraca(w, r, cassandraSession)
	})
	mux.HandleFunc("/api/igrac/sveSobe", func(w http.ResponseWriter, r *http.Request) {
		handleSveSobeIgraca(w, r, cassandraSession)
	})

	return mux
}

func handleSveSobeIgraca(w http.ResponseWriter, r *http.Request, s *gocql.Session) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	sveSobeJSON := sveSobeIgraca(username, s)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string {
		"sobe": string(sveSobeJSON),
	})
}

func pribaviNepozeljneReci(cassandraSession *gocql.Session) []string {
	query := "SELECT * FROM nepozeljne_reci"

	scanner := cassandraSession.Query(query).Iter().Scanner()
	reci := make([]string, 0)

	for scanner.Next() {
		var rec string
		if err := scanner.Scan(&rec); err != nil {
			log.Printf("Greška prilikom skeniranja nepoželjne reči iz cassandra baze podataka: %v\n", err)
			continue
		}

		reci = append(reci, rec)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Greška prilikom pribavljanja nepoželjnih reči iz cassandra baze podataka: %v\n", err)
	}

	return reci
}

func handlePrijaviIgraca(w http.ResponseWriter, r *http.Request, s *gocql.Session) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var podaci map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&podaci); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	username := podaci["username"].(string)
	lozinka := podaci["password"].(string)

	query := "SELECT * FROM igraci WHERE username = ?"

	igrac := make(map[string]interface{})
	if err := s.Query(query, username).MapScan(igrac); err != nil {
		log.Printf("Greška prilikom pribavljanja igrača iz cassandra baze podataka: %v\n", err)
		if err == gocql.ErrNotFound {
			http.Error(w, fmt.Sprintf("Ne postoji igrač sa username: %v", username), http.StatusNotFound)
		} else {
			http.Error(w, "Greška", http.StatusInternalServerError)
		}
		return
	}

	lozinka_hes := []byte(igrac["lozinka_hes"].(string))

	if err := bcrypt.CompareHashAndPassword(lozinka_hes, []byte(lozinka)); err != nil {
		log.Printf("Greška prilikom prijavljivanja: %v\n", err)
		http.Error(w, "Lozinka nije ispravna", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{} {
		"status": "uspešna prijava",
		"igrac": map[string]interface{} {
			"username": username,
			"datumRodjenja": igrac["datum_rodjenja"],
		},
	})
	log.Printf("Igrač %v uspešno prijavljen\n", username)
}

func handleRegistrujIgraca(w http.ResponseWriter, r *http.Request, s *gocql.Session) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var podaci map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&podaci); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	username := podaci["username"].(string)
	if !validanUsername(username) {
		log.Printf("Nevalidan username: %v\n", username)
		http.Error(w, "{ \"greska\": \"USERNAME\" }", http.StatusBadRequest)
		return
	}

	lozinka := podaci["password"].(string)
	if !validnaLozinka(lozinka) {
		log.Printf("Nevalidna lozinka: %v\n", lozinka)
		http.Error(w, "{ \"greska\": \"LOZINKA\" }", http.StatusBadRequest)
		return
	}

	datumRodjenja, err := time.Parse(time.DateOnly, podaci["datumRodjenja"].(string))
	if err != nil {
		log.Printf("Greška prilikom parsovanja datuma rođenja: %v\n", err)
		http.Error(w, "Pogrešni podaci", http.StatusBadRequest)
		return
	}

	lozinkaHes, err := bcrypt.GenerateFromPassword([]byte(lozinka), 10)
	if err != nil {
		log.Printf("Greška prilikom generisanja heša lozinke: %v\n", err)
		http.Error(w, "{ \"greska\": \"LOZINKA\" }", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO igraci (username, lozinka_hes, datum_rodjenja) " +
			"VALUES (?, ?, ?) IF NOT EXISTS"

	dest := make(map[string]interface{})
	uspesnoRegistrovan, err := s.Query(query, username, string(lozinkaHes), datumRodjenja).MapScanCAS(dest)
	if err != nil {
		log.Printf("Greška prilikom registrovanja igrača u cassandra bazu podataka: %v\n", err)
		http.Error(w, "Greška.", http.StatusInternalServerError)
		return
	}

	if !uspesnoRegistrovan {
		log.Printf("Greška prilikom registrovanja, igrač %v već postoji\n", username)
		http.Error(w, "{ \"greska\": \"USERNAME\" }", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string {
		"status": "created",
	})
	log.Printf("Igrač %v uspešno registrovan\n", username)
}

func validanUsername(username string) bool {
	for i := range username {
		if username[i] >= 'a' && username[i] <= 'z' {
			continue
		}
		if username[i] >= '0' && username[i] <= '9' {
			continue
		}

		return false
	}

	return true
}

func validnaLozinka(lozinka string) bool {
	if len(lozinka) < 8 || len(lozinka) > 72 {
		return false
	}

	specKarakteri := []byte { '!', ' ', '@', '#', '$', '%', '&' }

	maloSlovo := false
	velikoSlovo := false
	broj := false
	specKarakter := false

	for i := range lozinka {
		if lozinka[i] >= 'a' && lozinka[i] <= 'z' {
			maloSlovo = true
		}
		if lozinka[i] >= 'A' && lozinka[i] <= 'Z' {
			velikoSlovo = true
		}
		if lozinka[i] >= '0' && lozinka[i] <= '9' {
			broj = true
		}
		for j := range specKarakteri {
			if lozinka[i] == specKarakteri[j] {
				specKarakter = true
				break
			}
		}

		if maloSlovo && velikoSlovo && broj && specKarakter {
			return true
		}
	}

	return false
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

func handleWS(w http.ResponseWriter, r *http.Request, rdb *redis.Client, ctx context.Context, s *gocql.Session) {
	conn, err := upg.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Handler upgrader greška: %v\n", err)
		return
	}

	igrac := &Igrac {
		Username: fmt.Sprintf("gost_%s", uuid.NewString()[0:13]),
		DatumRodjenja: time.Now(),
		Conn: conn,
	}

	aktivniIgraciMux.Lock()
	aktivniIgraci[igrac.Username] = igrac
	aktivniIgraciMux.Unlock()

	go igrac.CitajWSPoruke(ctx, rdb, s)
}

func NoviServer(cassandraSession *gocql.Session) (*http.Server, int) {
	noviServer := Server {
		port: 8080,
	}

	return &http.Server {
		Addr: fmt.Sprintf(":%d", noviServer.port),
		Handler: noviServer.HandlerInit(cassandraSession),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}, noviServer.port
}

func NadjiAktivnogIgraca(username string) *Igrac {
	aktivniIgraciMux.RLock()
	defer aktivniIgraciMux.RUnlock()
	return aktivniIgraci[username]
}

func DiskonektujIgraca(igracUsername string) {
	aktivniIgraciMux.Lock()
	delete(aktivniIgraci, igracUsername)
	aktivniIgraciMux.Unlock()
	log.Printf("Igrač %v diskonektovan\n", igracUsername)
}
