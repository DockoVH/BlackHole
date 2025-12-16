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

type Igrac struct {
	UUID string `json:"uuid"`
	Ime string `json:"ime"`
	Conn *websocket.Conn `json:"-"`
	cetPubSub *redis.PubSub
	sobaUUID string
}

func NoviIgrac(c *websocket.Conn) *Igrac {
	return &Igrac {
		UUID: uuid.NewString(),
		Ime: "",
		Conn: c,
	}
}

func (igrac *Igrac) CitajWSPoruke(ctx context.Context, rdb *redis.Client) {
	defer func() {
        if err := igrac.Conn.Close(); err != nil {
            log.Printf("CitajWSPoruke greška: %v\n", err)
        }
        DiskonektujIgraca(igrac.UUID)
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
                    log.Printf("citajWSPoruke greška: conn: %v, err: %v\n", igrac.Conn, err)
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
					continue
				}
				igrac.cetPubSub = rdb.Subscribe(ctx, fmt.Sprintf("soba:%s:cet-pub-sub", soba.UUID))
				igrac.sobaUUID = soba.UUID
				go igrac.handleCetPoruke()

				if len(soba.Igraci) == 2 {
					igrac.PosaljiOdgovorWS(poruka.NovaPoruka("Start", "Igra je počela.").Marshal())
					go soba.Start()
				} else {
					igrac.PosaljiOdgovorWS(poruka.NovaPoruka("Cekanje", "Nema dovoljno igrača za početak igre.").Marshal())
				}
			case "Cet_Poruka":
				cetPoruka := poruka.CetPoruka(igrac.Ime , dobijenaPoruka.Sadrzaj)
				if cetPoruka.Tip == "Greska" {
					igrac.PosaljiOdgovorWS(cetPoruka.Marshal())
					continue
				}

				if err := rdb.Publish(ctx, fmt.Sprintf("soba:%s:cet-pub-sub", igrac.sobaUUID), cetPoruka.Marshal()).Err(); err != nil {
					log.Printf("Greška prilikom slanje poruke u kanal soba:%s:cet-pub-sub: %v\n", igrac.sobaUUID, err)
				}
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
