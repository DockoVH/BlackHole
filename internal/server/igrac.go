package server

import (
	"fmt"
	"context"
	"time"
	"log"
	"encoding/json"
	"strings"
	"slices"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/google/uuid"
	"github.com/apache/cassandra-gocql-driver/v2"
	"github.com/texttheater/golang-levenshtein/levenshtein"

	"BlackHole/internal/poruka"
)

type CetPoruka struct {
	Vreme time.Time
	IgracUsername string
	Sadrzaj string
}

type Potez struct {
	Vreme time.Time `json:"vreme"`
	IgracUsername string `json:"igrac_username"`
	IndeksPolja int `json:"indeks_polja"`
}

type Igrac struct {
	Username string `json:"username"`
	DatumRodjenja time.Time `json:"datumRodjenja"`
	Conn *websocket.Conn `json:"-"`
	cetPubSub *redis.PubSub
	sobaUUID string
}

func (igrac *Igrac) CitajWSPoruke(ctx context.Context, rdb *redis.Client, s *gocql.Session) {
	defer func() {
        if err := igrac.Conn.Close(); err != nil {
            log.Printf("CitajWSPoruke greška: %v\n", err)
        }
        DiskonektujIgraca(igrac.Username)
        if igrac.cetPubSub != nil {
			igrac.cetPubSub.Close()
		}
    }()

	igrac.Conn.SetReadLimit(maxMessageSize)
    igrac.Conn.SetPongHandler(func(string) error {
        igrac.Conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    igraTraje := false

    for {
        _, primljenaPoruka, err := igrac.Conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                    log.Printf("citajWSPoruke igrač %v, konekcija zatvorena,\n", igrac.Username)
            }
            break
        }

        var dobijenaPoruka poruka.Poruka
        if err := json.Unmarshal(primljenaPoruka, &dobijenaPoruka); err != nil {
			log.Printf("Greška prilikom unmarshal-ovanja poruke: %v\n", err)
			igrac.PosaljiOdgovorWS(poruka.Greska(fmt.Sprintf("Greška prilikom unmarshal-ovanja poruke: %v\n", err)).Marshal())
			continue
        }

		switch dobijenaPoruka.Tip {
			case "Dodaj_U_Sobu":
				if igraTraje {
					continue
				}

				var podaci struct {
					Kod string `json:"kod"`
					Username string `json:"username"`
					DatumRodjenja string `json:"datumRodjenja"`
				}

				if err := json.Unmarshal([]byte(dobijenaPoruka.Sadrzaj), &podaci); err != nil {
					log.Printf("Dodaj_U_Sobu Unmarshal greška: %v\n", err)
					return
				}

				if podaci.Username != "" {
					aktivniIgraciMux.Lock()

					delete(aktivniIgraci, igrac.Username)
					igrac.Username = podaci.Username
					if datum, err := time.Parse(time.RFC3339, podaci.DatumRodjenja); err != nil {
						igrac.DatumRodjenja = time.Now()
					} else {
						igrac.DatumRodjenja = datum
					}
					aktivniIgraci[igrac.Username] = igrac

					aktivniIgraciMux.Unlock()
				}

				soba, err := DodajUSobu(podaci.Kod, igrac, ctx, rdb)
				if err != nil {
					igrac.PosaljiOdgovorWS(poruka.Greska(fmt.Sprintf("Greška prilikom dodavanja u sobu: %v", err)).Marshal())
					return
				}

				redniBrojIgraca := 0
				for i, igrac := range soba.Igraci {
					if igrac.Username == soba.IgracNaRedu {
						redniBrojIgraca = i + 1
					}
				}

				sobaPodaciPoruka := poruka.SobaPodaci(soba.Kod, soba.IgracNaRedu, redniBrojIgraca)
				if sobaPodaciPoruka.Tip == "Greska" {
					return
				}
				igrac.PosaljiOdgovorWS(sobaPodaciPoruka.Marshal())

				igracPodaciPoruka := poruka.IgracPodaci(igrac.Username, igrac.DatumRodjenja)
				if igracPodaciPoruka.Tip == "Greska" {
					return
				}
				igrac.PosaljiOdgovorWS(igracPodaciPoruka.Marshal())

				igrac.cetPubSub = rdb.Subscribe(ctx, fmt.Sprintf("soba:%s:cet-pub-sub", soba.UUID))
				igrac.sobaUUID = soba.UUID
				go igrac.handleCetPoruke()

				igraTraje = true

				if len(soba.Igraci) == 2 {
					go soba.Start(ctx, rdb)
				} else {
					igrac.PosaljiOdgovorWS(poruka.NovaPoruka("Cekanje", "Nema dovoljno igrača za početak igre.").Marshal())
				}

			case "Cet_Poruka":
				if !igraTraje {
					continue
				}

				cetPorukaSlanje := poruka.CetPoruka(igrac.Username , dobijenaPoruka.Sadrzaj)
				if cetPorukaSlanje.Tip == "Greska" {
					igrac.PosaljiOdgovorWS(cetPorukaSlanje.Marshal())
					continue
				}

				if err := rdb.Publish(ctx, fmt.Sprintf("soba:%s:cet-pub-sub", igrac.sobaUUID), cetPorukaSlanje.Marshal()).Err(); err != nil {
					log.Printf("Greška prilikom slanje poruke u kanal soba:%s:cet-pub-sub: %v\n", igrac.sobaUUID, err)
					continue
				}

				cetPorukaRedis := CetPoruka {
					Vreme: time.Now(),
					IgracUsername: igrac.Username,
					Sadrzaj: dobijenaPoruka.Sadrzaj,
				}

				cetPorukaRedisJSON, err := json.Marshal(cetPorukaRedis)
				if err != nil {
					log.Printf("Greška prilikom marshalovanja cetPorukaRedis: %v\n", err)
					continue
				}

				hesKljuc := fmt.Sprintf("soba:%s:sve-poruke", igrac.sobaUUID)
				porukaKljuc := fmt.Sprintf("poruka:%s", uuid.NewString())
				if err := rdb.HSet(ctx, hesKljuc, porukaKljuc, cetPorukaRedisJSON).Err(); err != nil {
					log.Printf("Greška prilikom dodavanja poruke u redis bazu podataka: %v\n", err)
				}

			case "Potez":
				var noviPotez Potez

				if err := json.Unmarshal([]byte(dobijenaPoruka.Sadrzaj), &noviPotez); err != nil {
					log.Printf("Potez Unmarshal greška: %v\n", err)
					return
				}

				OdigrajPotez(igrac.sobaUUID, noviPotez.IndeksPolja, noviPotez.IgracUsername, ctx, rdb, s)

			case "Kraj_Igre":
				if !igraTraje {
					continue
				}

				soba := UcitajSobuRedisDB(igrac.sobaUUID, ctx, rdb)
				if soba == nil {
					log.Fatal("Kraj igre greška.")
				}

				soba.Broadcast("Kraj_Igre", "")
				igraTraje = false

			default:
				igrac.PosaljiOdgovorWS(primljenaPoruka)
		}
    }
}

func (igrac *Igrac) PosaljiOdgovorWS(wsPoruka []byte) {
	igrac.Conn.SetWriteDeadline(time.Now().Add(writeWait))

	writer, err := igrac.Conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		log.Printf("PosaljiOdgovorWS igrac.Conn.NextWriter greška: %v\n", err)
		return
	}

	writer.Write(wsPoruka)

	if err = writer.Close(); err != nil {
		log.Printf("PosaljiOdgovorWS writer.Close() greška: %v\n", err)
	}
}

func (igrac *Igrac) handleCetPoruke() {
	porukaChan := igrac.cetPubSub.Channel()

	for primljenaPoruka := range porukaChan {
		var porukaZaSlanje poruka.Poruka
        if err := json.Unmarshal([]byte(primljenaPoruka.Payload), &porukaZaSlanje); err != nil {
			log.Printf("Greška prilikom unmarshal-ovanja cet poruke: %v\n", err)
		}

		sada := time.Now()
		datum18Rodjendan := igrac.DatumRodjenja.AddDate(18, 0, 0)

		if sada.Before(datum18Rodjendan) {
			sadrzaj := strings.ReplaceAll(porukaZaSlanje.Sadrzaj, string(znakoviInterpunkcije[0]), " ")
			for i := range len(znakoviInterpunkcije) - 1 {
				sadrzaj = strings.ReplaceAll(sadrzaj,  string(znakoviInterpunkcije[i + 1]), " ")
			}

			reci := strings.Split(sadrzaj, " ")

			for _, rec := range nepozeljneReci {
				for i := range reci {
					if levenshtein.DistanceForStrings([]rune(strings.ToLower(reci[i])), []rune(rec), levenshtein.DefaultOptions) < 3 {
						reci[i] = strings.Repeat("*", len(reci[i]))
					}
				}
			}

			noviSadrzaj := []byte(strings.Join(reci, " "))
			for i := range porukaZaSlanje.Sadrzaj {
				if slices.Contains(znakoviInterpunkcije, porukaZaSlanje.Sadrzaj[i]) {
					noviSadrzaj[i] = porukaZaSlanje.Sadrzaj[i]
				}
			}

			porukaZaSlanje.Sadrzaj = string(noviSadrzaj)
		}

		igrac.PosaljiOdgovorWS(porukaZaSlanje.Marshal())
	}
}
