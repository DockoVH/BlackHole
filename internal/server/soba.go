package server

import (
	"fmt"
	"log"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/google/uuid"

	"BlackHole/internal/poruka"
)

type Soba struct {
	UUID string
	Kod string
	Igraci []*Igrac
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
		if err := sacuvajSobuRedisDB(novaSoba, ctx, rdb); err != nil {
			log.Printf("DodajUSobu len(sveSobe) = 0, greška prilikom čuvanja sobe u redis bazu: %v\n", err)
			return nil, fmt.Errorf("Nemoguće dodavanje igrača u sobu.")
		}

		if kod != "" {
			log.Printf("Igrač %v dodat u sobu sa uuid: %v, kod: %v\n", igrac.Username, novaSoba.UUID, kod)
		} else {
			log.Printf("Igrač %v dodat u sobu sa uuid: %v\n", igrac.Username, novaSoba.UUID)
		}
		return novaSoba, nil
	}

	if kod != "" {
		for i := range sveSobe {
			if sveSobe[i].Kod == kod {
				if len(sveSobe[i].Igraci) == 2 {
					return nil, fmt.Errorf("Soba sa zadatim kodom je puna.")
				}
				sveSobe[i].Igraci = append(sveSobe[i].Igraci, igrac)
				if err := sacuvajSobuRedisDB(sveSobe[i], ctx, rdb); err != nil {
					log.Printf("DodajUSobu kod != \"\", greška prilikom čuvanja sobe u redis bazu: %v\n", err)
					return nil, fmt.Errorf("Nemoguće dodavanje igrača u sobu.")
				}

				log.Printf("Igrač %v dodat u sobu sa uuid: %v, kod: %v\n", igrac.Username, sveSobe[i].UUID, kod)
				return sveSobe[i], nil
			}
		}

		novaSoba := napraviSobu(kod, igrac)
		if err := sacuvajSobuRedisDB(novaSoba, ctx, rdb); err != nil {
			log.Printf("DodajUSobu kod != \"\", nova soba, greška prilikom čuvanja sobe u redis bazu: %v\n", err)
			return nil, fmt.Errorf("Nemoguće dodavanje igrača u sobu.")
		}

		log.Printf("Igrač %v dodat u sobu sa uuid: %v, kod: %v\n", igrac.Username, novaSoba.UUID, kod)
		return novaSoba, nil
	}

	for i := range sveSobe {
		if len(sveSobe[i].Igraci) < 2 {
			sveSobe[i].Igraci = append(sveSobe[i].Igraci, igrac)
			if err := sacuvajSobuRedisDB(sveSobe[i], ctx, rdb); err != nil {
				log.Printf("DodajUSobu: greška prilikom čuvanja sobe u redis bazu: %v\n", err)
				return nil, fmt.Errorf("Nemoguće dodavanje igrača u sobu.")
			}

			log.Printf("Igrač %v dodat u sobu sa uuid: %v\n", igrac.Username, sveSobe[i].UUID)
			return sveSobe[i], nil
		}
	}

	novaSoba := napraviSobu(kod, igrac)
	if err := sacuvajSobuRedisDB(novaSoba, ctx, rdb); err != nil {
		log.Printf("DodajUSobu sve sobe pune, greška prilikom čuvanja sobe u redis bazu: %v\n", err)
		return nil, fmt.Errorf("Nemoguće dodavanje igrača u sobu.")
	}

	log.Printf("Igrač %v dodat u sobu sa uuid: %v\n", igrac.Username, novaSoba.UUID)
	return novaSoba, nil
}

func sveSobeFromHes(hes map[string]string) []*Soba {
	sveSobe := make([]*Soba, 0)

	for _, vrednost := range hes {
		var sobaPodaci struct {
			UUID string `json:"uuid"`
			Kod string	`json:"kod"`
			IgraciUsernames []string `json:"igraci_usernames"`
		}
		if err := json.Unmarshal([]byte(vrednost), &sobaPodaci); err != nil {
			log.Printf("Greška prilikom konvertovanja hes-a u sobu: %v\n", err)
			continue
		}

		soba := &Soba {
			UUID: sobaPodaci.UUID,
			Kod: sobaPodaci.Kod,
			Igraci: make([]*Igrac, 0),
		}

		for _, igracUsername := range sobaPodaci.IgraciUsernames {
			igrac := NadjiAktivnogIgraca(igracUsername)
			if igrac != nil {
				soba.Igraci = append(soba.Igraci, igrac)
			} else {
				log.Printf("Greška soba.uuid %v:, igrač %v nije aktivan!\n", soba.UUID, igracUsername)
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

func sacuvajSobuRedisDB(soba *Soba, ctx context.Context, rdb *redis.Client) error {
	igraciUsernames := make([]string, len(soba.Igraci))
	for i, igrac := range soba.Igraci {
		igraciUsernames[i] = igrac.Username
	}

	sobaPodaci := struct {
		UUID string `json:"uuid"`
		Kod string	`json:"kod"`
		IgraciUsernames []string `json:"igraci_usernames"`
	}{
		UUID: soba.UUID,
		Kod: soba.Kod,
		IgraciUsernames: igraciUsernames,
	}

	sobaJSON, err := json.Marshal(sobaPodaci)
	if err != nil {
		return err
	}

	return rdb.HSet(ctx, "sve-sobe", fmt.Sprintf("soba:%s", soba.UUID), sobaJSON).Err()
}

func UcitajSobuRedisDB(sobaUUID string, ctx context.Context, rdb *redis.Client) *Soba {
	sobaJSON, err := rdb.HGet(ctx, "sve-sobe", fmt.Sprintf("soba:%s", sobaUUID)).Result()
	if err != nil {
		log.Printf("Greška prilikom učitavanja soba iz redis baze podataka: %v\n", err)
		return nil
	}

	var sobaPodaci struct {
		UUID string `json:"uuid"`
		Kod string `json:"kod"`
		IgraciUsernames []string `json:"igraci_usernames"`
	}

	if err := json.Unmarshal([]byte(sobaJSON), &sobaPodaci); err != nil {
		log.Printf("UcitajSobuRedisDB() Unmarshal greška: %v\n", err)
		return nil
	}

	soba := &Soba {
		UUID: sobaPodaci.UUID,
		Kod: sobaPodaci.Kod,
		Igraci: make([]*Igrac, 0),
	}

	for i := range sobaPodaci.IgraciUsernames {
		igrac := NadjiAktivnogIgraca(sobaPodaci.IgraciUsernames[i])
		if igrac != nil {
			soba.Igraci = append(soba.Igraci, igrac)
		} else {
			log.Printf("UcitajSobuRedisDB() greška: igrač %v nije aktivan!\n", sobaPodaci.IgraciUsernames[i])
		}
	}

	return soba
}

func (soba *Soba) Broadcast(tipPoruke string, sadrzajPoruke string) {
	for i := range soba.Igraci {
		soba.Igraci[i].PosaljiOdgovorWS(poruka.NovaPoruka(tipPoruke, sadrzajPoruke).Marshal())
	}
}

func (soba *Soba) Start() {
	for vreme := range 3 {
		soba.Broadcast("Pocetak_Igre", fmt.Sprintf("Igra počinje za %vs.", 3 - vreme))
		time.Sleep(time.Second)
	}
	soba.Broadcast("Start", "Igra je počela.")
}
