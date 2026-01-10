package server

import (
	"fmt"
	"context"
	"time"
	"log"
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/google/uuid"

	"BlackHole/internal/poruka"
)

type CetPoruka struct {
	Vreme time.Time
	IgracUsername string
	Sadrzaj string
}

type Potez struct {
	Vreme time.Time
	IgracUsername string
	IndeksPolja int
}

type Igrac struct {
	Username string `json:"username"`
	DatumRodjenja time.Time `json:"datumRodjenja"`
	Conn *websocket.Conn `json:"-"`
	cetPubSub *redis.PubSub
	sobaUUID string
}

func (igrac *Igrac) CitajWSPoruke(ctx context.Context, rdb *redis.Client) {
	defer func() {
        if err := igrac.Conn.Close(); err != nil {
            log.Printf("CitajWSPoruke greška: %v\n", err)
        }
        log.Printf("Debug: prekinuta konekcija igrac: %v\n", igrac.Username)
        DiskonektujIgraca(igrac.Username)
        igrac.cetPubSub.Close()
    }()

	igrac.Conn.SetReadLimit(maxMessageSize)
    igrac.Conn.SetPongHandler(func(string) error {
        igrac.Conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

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
				soba, err := DodajUSobu(dobijenaPoruka.Sadrzaj, igrac, ctx, rdb)
				if err != nil {
					igrac.PosaljiOdgovorWS(poruka.Greska(fmt.Sprintf("Greška prilikom dodavanja u sobu: %v", err)).Marshal())
					return
				}

				igracPodaciPoruka := poruka.IgracPodaci(igrac.Username, igrac.DatumRodjenja)
				igrac.PosaljiOdgovorWS(igracPodaciPoruka.Marshal())
				if igracPodaciPoruka.Tip == "Greska" {
					return
				}

				igrac.cetPubSub = rdb.Subscribe(ctx, fmt.Sprintf("soba:%s:cet-pub-sub", soba.UUID))
				igrac.sobaUUID = soba.UUID
				go igrac.handleCetPoruke()

				if len(soba.Igraci) == 2 {
					go soba.Start()
				} else {
					igrac.PosaljiOdgovorWS(poruka.NovaPoruka("Cekanje", "Nema dovoljno igrača za početak igre.").Marshal())
				}

			case "Igrac_Podaci":
				var igracPodaci struct {
					Username string
					DatumRodjenja time.Time
				}

				if err := json.Unmarshal([]byte(dobijenaPoruka.Sadrzaj), &igracPodaci); err != nil {
					log.Printf("Igrac_Podaci Unmarshal greška: %v\n", err)
					return
				}

				igrac.Username = igracPodaci.Username
				igrac.DatumRodjenja = igracPodaci.DatumRodjenja

			case "Cet_Poruka":
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
				

			case "Kraj_Igre":
				soba := UcitajSobuRedisDB(igrac.sobaUUID, ctx, rdb)
				if soba == nil {
					log.Fatal("Kraj igre greška.")
				}

				soba.Broadcast("Kraj_Igre", "")
				return

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

	for cetPoruka := range porukaChan {
		igrac.PosaljiOdgovorWS([]byte(cetPoruka.Payload))
	}
}
