package server

import (
	"fmt"
	"log"
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/google/uuid"

	"BlackHole/internal/poruka"
)

type Soba struct {
	UUID string
	Kod string
	Igraci []*Igrac
	//Potezi []Potez
}

func DodajUSobu(kod string, igrac *Igrac, ctx context.Context, rdb *redis.Client) (*Soba, error) {
	hes, err := rdb.HGetAll(ctx, "sve-sobe").Result()

	if err != nil {
		log.Printf("Greška prilikom pribavljanja svih soba: %v\n", err)
		return nil, fmt.Errorf("Nemoguće dodavanje igrača u sobu")
	}

	sveSobe := sveSobeFromHes(hes)
	if sveSobe == nil {
		log.Printf("DodajUSobu greška prilikom pribavljanja svih soba!\n")
		return nil, fmt.Errorf("Nemoguće dodavanje igrača u sobu")
	}

	if len(sveSobe) == 0 {
		novaSoba := napraviSobu(kod, igrac)
		if err := sacuvajSobuURedis(novaSoba, ctx, rdb); err != nil {
			log.Printf("DodajUSobu len(sveSobe) = 0, greška prilikom čuvanja sobe u redis bazu: %v\n", err)
			return nil, fmt.Errorf("Nemoguće dodavanje igrača u sobu.")
		}

		log.Printf("Igrač sa uuid: %v dodat u sobu sa uuid: %v, kod: %v\n", igrac.UUID, novaSoba.UUID, kod)
		return novaSoba, nil
	}

	if kod != "" {
		for i := range sveSobe {
			if sveSobe[i].Kod == kod {
				if len(sveSobe[i].Igraci) == 2 {
					return nil, fmt.Errorf("Soba sa zadatim kodom je puna.")
				}
				sveSobe[i].Igraci = append(sveSobe[i].Igraci, igrac)
				if err := sacuvajSobuURedis(sveSobe[i], ctx, rdb); err != nil {
					log.Printf("DodajUSobu kod != \"\", greška prilikom čuvanja sobe u redis bazu: %v\n", err)
					return nil, fmt.Errorf("Nemoguće dodavanje igrača u sobu.")
				}

				log.Printf("Igrač sa uuid: %v dodat u sobu sa uuid: %v, kod: %v\n", igrac.UUID, sveSobe[i].UUID, kod)
				return sveSobe[i], nil
			}
		}

		novaSoba := napraviSobu(kod, igrac)
		if err := sacuvajSobuURedis(novaSoba, ctx, rdb); err != nil {
			log.Printf("DodajUSobu kod != \"\", nova soba, greška prilikom čuvanja sobe u redis bazu: %v\n", err)
			return nil, fmt.Errorf("Nemoguće dodavanje igrača u sobu.")
		}

		log.Printf("Igrač sa uuid: %v dodat u sobu sa uuid: %v, kod: %v\n", igrac.UUID, novaSoba.UUID, kod)
		return novaSoba, nil
	}

	for i := range sveSobe {
		if len(sveSobe[i].Igraci) < 2 {
			sveSobe[i].Igraci = append(sveSobe[i].Igraci, igrac)
			if err := sacuvajSobuURedis(sveSobe[i], ctx, rdb); err != nil {
				log.Printf("DodajUSobu: greška prilikom čuvanja sobe u redis bazu: %v\n", err)
				return nil, fmt.Errorf("Nemoguće dodavanje igrača u sobu.")
			}

			log.Printf("Igrač sa uuid: %v dodat u sobu sa uuid: %v, kod: %v\n", igrac.UUID, sveSobe[i].UUID, kod)
			return sveSobe[i], nil
		}
	}

	novaSoba := napraviSobu(kod, igrac)
	if err := sacuvajSobuURedis(novaSoba, ctx, rdb); err != nil {
		log.Printf("DodajUSobu sve sobe pune, greška prilikom čuvanja sobe u redis bazu: %v\n", err)
		return nil, fmt.Errorf("Nemoguće dodavanje igrača u sobu.")
	}

	log.Printf("Igrač sa uuid: %v dodat u sobu sa uuid: %v, kod: %v\n", igrac.UUID, novaSoba.UUID, kod)
	return novaSoba, nil
}

func sveSobeFromHes(hes map[string]string) []*Soba {
	sveSobe := make([]*Soba, 0)

	for _, v := range hes {
		var sobaPodaci struct {
			UUID string `json:"uuid"`
			Kod string	`json:"kod"`
			IgraciUUID []string `json:igraci_uuid`
		}
		if err := json.Unmarshal([]byte(v), &sobaPodaci); err != nil {
			log.Printf("Greška prilikom konvertovanja hes-a u sobu: %v\n", err)
			continue
		}

		soba := &Soba {
			UUID: sobaPodaci.UUID,
			Kod: sobaPodaci.Kod,
			Igraci: make([]*Igrac, 0),
		}

		for _, igracUUID := range sobaPodaci.IgraciUUID {
			igrac := NadjiAktivnogIgraca(igracUUID)
			if igrac != nil {
				soba.Igraci = append(soba.Igraci, igrac)
			} else {
				log.Printf("Greška soba.uuid %v:, igrac sa uuid: %d nije aktivan!\n", soba.UUID, igracUUID)
			}
		}

		sveSobe = append(sveSobe, soba)
	}

	return sveSobe
}

func napraviSobu(kod string, igrac *Igrac) *Soba {
	return &Soba {
		UUID: uuid.NewString(),
		Kod: kod,
		Igraci: []*Igrac { igrac },
	}
}

func sacuvajSobuURedis(soba *Soba, ctx context.Context, rdb *redis.Client) error {
	igraciUUID := make([]string, len(soba.Igraci))
	for i, igrac := range soba.Igraci {
		igraciUUID[i] = igrac.UUID
	}

	sobaPodaci := struct {
			UUID string `json:"uuid"`
			Kod string	`json:"kod"`
			IgraciUUID []string `json:igraci_uuid`
		}{
			UUID: soba.UUID,
			Kod: soba.Kod,
			IgraciUUID: igraciUUID,
		}

	sobaJSON, err := json.Marshal(sobaPodaci)
	if err != nil {
		return err
	}

	return rdb.HSet(ctx, "sve-sobe", fmt.Sprintf("soba:%s", soba.UUID), sobaJSON).Err()
}

func (soba *Soba) Start() {
	for i := range soba.Igraci {
		soba.Igraci[i].PosaljiOdgovorWS(poruka.NovaPoruka("Test", "soba.Start() test poruka").Marshal())
	}
}
