package server

import (
	"fmt"
	"log"
	"context"
	"encoding/json"
	"time"
	"math"
	"slices"

	"github.com/redis/go-redis/v9"
	"github.com/google/uuid"
	"github.com/apache/cassandra-gocql-driver/v2"

	"BlackHole/internal/poruka"
)

const (
	slobodnoPolje = 0
	zauzetoPoljePrviIgrac = 1
	zauzetoPoljeDrugiIgrac = 2
	krajZauzetoPoljePrviIgrac = 51
	krajZauzetoPoljeDrugiIgrac = 52
	crnaRupaPolje = 100
)

type Soba struct {
	UUID string
	Kod string
	Igraci []*Igrac
	Tabla []Polje
	IgracNaRedu string
	BrojNaRedu int
}

type Polje struct {
	Indeks int `json:"indeks"`
	Stanje int `json:"stanje"`
	Vrednost int `json:"vrednost"`
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
			Tabla []Polje `json:"tabla_polja"`
			IgracNaRedu string `json:"igrac_na_redu"`
			BrojNaRedu int `json:"broj_na_redu"`
		}
		if err := json.Unmarshal([]byte(vrednost), &sobaPodaci); err != nil {
			log.Printf("Greška prilikom konvertovanja hes-a u sobu: %v\n", err)
			continue
		}

		soba := &Soba {
			UUID: sobaPodaci.UUID,
			Kod: sobaPodaci.Kod,
			Igraci: make([]*Igrac, 0),
			Tabla: sobaPodaci.Tabla,
			IgracNaRedu: sobaPodaci.IgracNaRedu,
			BrojNaRedu: sobaPodaci.BrojNaRedu,
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
	sobaKod := kod
	if kod == "" {
		sobaKod = uuid.NewString()[0:4]
	}

	novaSoba := Soba {
		UUID: uuid.NewString(),
		Kod: sobaKod,
		Igraci: []*Igrac { igrac },
		Tabla: make([]Polje, 21),
		BrojNaRedu: 1,
	}

	for i := range novaSoba.Tabla {
		novaSoba.Tabla[i].Indeks = i
	}

	return &novaSoba
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
		Tabla []Polje `json:"tabla_polja"`
		IgracNaRedu string `json:"igrac_na_redu"`
		BrojNaRedu int `json:"broj_na_redu"`
	}{
		UUID: soba.UUID,
		Kod: soba.Kod,
		IgraciUsernames: igraciUsernames,
		Tabla: soba.Tabla,
		IgracNaRedu: soba.IgracNaRedu,
		BrojNaRedu: soba.BrojNaRedu,
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
		Tabla []Polje `json:"tabla_polja"`
		IgracNaRedu string `json:"igrac_na_redu"`
		BrojNaRedu int `json:"broj_na_redu"`
	}

	if err := json.Unmarshal([]byte(sobaJSON), &sobaPodaci); err != nil {
		log.Printf("UcitajSobuRedisDB() Unmarshal greška: %v\n", err)
		return nil
	}

	soba := &Soba {
		UUID: sobaPodaci.UUID,
		Kod: sobaPodaci.Kod,
		Igraci: make([]*Igrac, 0),
		Tabla: sobaPodaci.Tabla,
		IgracNaRedu: sobaPodaci.IgracNaRedu,
		BrojNaRedu: sobaPodaci.BrojNaRedu,
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

func obrisiSobuRedisDB(sobaUUID string, ctx context.Context, rdb *redis.Client) {
	sobaKljuc := fmt.Sprintf("soba:%s", sobaUUID)
	poteziKljuc := fmt.Sprintf("soba:%s:svi-potezi", sobaUUID)
	porukeKljuc := fmt.Sprintf("soba:%s:sve-poruke", sobaUUID)

	if err := rdb.HDel(ctx, "sve-sobe", sobaKljuc).Err(); err != nil {
		log.Printf("Greška pirlikom brisanja sobe iz redis baze podataka: %v\n", err)
	}
	if err := rdb.Del(ctx, poteziKljuc).Err(); err != nil {
		log.Printf("Greška pirlikom brisanja poteza iz redis baze podataka: %v\n", err)
	}
	if err := rdb.Del(ctx, porukeKljuc).Err(); err != nil {
		log.Printf("Greška pirlikom brisanja poruka iz redis baze podataka: %v\n", err)
	}
}

func OdigrajPotez(sobaUUID string, indeksPolja int, igracUsername string, ctx context.Context, rdb *redis.Client, s *gocql.Session) {
	soba := UcitajSobuRedisDB(sobaUUID, ctx, rdb)
	if soba == nil {
		return
	}

	if igracUsername != soba.IgracNaRedu || soba.Tabla[indeksPolja].Stanje != slobodnoPolje {
		return
	}

	if soba.Igraci[0].Username == igracUsername {
		soba.Tabla[indeksPolja].Stanje = zauzetoPoljePrviIgrac
	} else {
		soba.Tabla[indeksPolja].Stanje = zauzetoPoljeDrugiIgrac
	}

	soba.Tabla[indeksPolja].Vrednost = soba.BrojNaRedu

	if soba.IgracNaRedu == soba.Igraci[0].Username {
		soba.IgracNaRedu = soba.Igraci[1].Username
	} else {
		soba.IgracNaRedu = soba.Igraci[0].Username
		soba.BrojNaRedu++
	}

	potezPodaci := Potez {
		Vreme: time.Now(),
		IgracUsername: igracUsername,
		IndeksPolja: indeksPolja,
	}

	potezJSON, err := json.Marshal(potezPodaci)
	if err != nil {
		log.Printf("Greška prilikom marshalovanja poteza: %v\n", err)
		return
	}
	hesKljuc := fmt.Sprintf("soba:%s:svi-potezi", sobaUUID)
	potezKljuc := fmt.Sprintf("potez:%s", uuid.NewString())
	if err = rdb.HSet(ctx, hesKljuc, potezKljuc, potezJSON).Err(); err != nil {
		log.Printf("Greška prilikom dodavanja poteza u redis bazu podataka: %v\n", err)
		return
	}

	brojSlobodnihPolja, crnaRupaIndeks := 0, 0
	for i, polje := range soba.Tabla {
		if polje.Stanje == slobodnoPolje {
			brojSlobodnihPolja++
			crnaRupaIndeks = i
		}
	}

	if brojSlobodnihPolja == 1 {
		soba.Tabla[crnaRupaIndeks].Stanje = crnaRupaPolje
		soba.IgracNaRedu = ""
	}

	tablaJSON, err := json.Marshal(soba.Tabla)
	if err != nil {
		log.Printf("Greška prilikom marshalovanja table: %v\n", err)
		return
	}

	if err := sacuvajSobuRedisDB(soba, ctx, rdb); err != nil {
		log.Printf("OdigrajPotez, greška prilikom čuvanja sobe u redis bazu: %v\n", err)
		return
	}

	redniBrojIgraca := 0
    if soba.IgracNaRedu == soba.Igraci[0].Username {
            redniBrojIgraca = 1
    } else if soba.IgracNaRedu == soba.Igraci[1].Username {
            redniBrojIgraca = 2
    }

    sobaPodaciPoruka := poruka.SobaPodaci(soba.Kod, soba.IgracNaRedu, redniBrojIgraca)

	soba.Broadcast("Potez", string(tablaJSON))
	soba.Broadcast(sobaPodaciPoruka.Tip, sobaPodaciPoruka.Sadrzaj)

	if brojSlobodnihPolja == 1 {
		go krajIgre(sobaUUID, crnaRupaIndeks, ctx, rdb, s)
	}
}

func odrediRedPolja(indeksPolja int) int {
	return int(math.Floor((math.Sqrt(float64(8 * indeksPolja + 1)) - 1) / 2) + 1)
}

func krajIgre(sobaUUID string, crnaRupaIndeks int, ctx context.Context, rdb *redis.Client, s *gocql.Session) {
	soba := UcitajSobuRedisDB(sobaUUID, ctx, rdb)
	if soba == nil {
		return
	}
	brojPolja := len(soba.Tabla)
	crnaRupaRed := odrediRedPolja(crnaRupaIndeks)
	poljaOkoCrneRupe := []int{ crnaRupaIndeks }

	kandidati := []struct {
		indeks int
		ocekivaniRed int
	}{
		{ crnaRupaIndeks - crnaRupaRed, crnaRupaRed - 1 },
		{ crnaRupaIndeks - crnaRupaRed + 1, crnaRupaRed - 1 },
		{ crnaRupaIndeks - 1, crnaRupaRed },
		{ crnaRupaIndeks + 1, crnaRupaRed },
		{ crnaRupaIndeks + crnaRupaRed, crnaRupaRed + 1 },
		{ crnaRupaIndeks + crnaRupaRed + 1, crnaRupaRed + 1 },
	}

	for _, kandidat := range kandidati {
		if kandidat.indeks >= 0 && kandidat.indeks < brojPolja && odrediRedPolja(kandidat.indeks) == kandidat.ocekivaniRed {
			poljaOkoCrneRupe = append(poljaOkoCrneRupe, kandidat.indeks)
		}
	}

	zbirPrviIgrac, zbirDrugiIgrac := 0, 0
	pobednik := ""

	for i := range soba.Tabla {
		if !slices.Contains(poljaOkoCrneRupe, soba.Tabla[i].Indeks) {
			if soba.Tabla[i].Stanje == zauzetoPoljePrviIgrac {
				soba.Tabla[i].Stanje = krajZauzetoPoljePrviIgrac
			} else {
				soba.Tabla[i].Stanje = krajZauzetoPoljeDrugiIgrac
			}
		} else {
			if soba.Tabla[i].Stanje == zauzetoPoljePrviIgrac {
				zbirPrviIgrac += soba.Tabla[i].Vrednost
			} else if soba.Tabla[i].Stanje == zauzetoPoljeDrugiIgrac {
				zbirDrugiIgrac += soba.Tabla[i].Vrednost
			}
		}
	}

	if zbirPrviIgrac < zbirDrugiIgrac {
		pobednik = soba.Igraci[0].Username
	} else if zbirPrviIgrac > zbirDrugiIgrac {
		pobednik = soba.Igraci[1].Username
	}

	rezultatIgre := struct {
		Pobednik string `json:"pobednik"`
		Tabla []Polje `json:"tabla"`
	}{
		Pobednik: pobednik,
		Tabla: soba.Tabla,
	}

	rezultatJSON, err := json.Marshal(rezultatIgre)
	if err != nil {
		log.Printf("Greška prilikom marshalovanja rezultata igre: %v\n", err)
		return
	}

	soba.Broadcast("Kraj_Igre", string(rezultatJSON))

	if err := sacuvajSobuCassandraDB(soba, pobednik, s, ctx, rdb); err != nil {
		log.Print(err)
		return
	}

	obrisiSobuRedisDB(sobaUUID, ctx, rdb)

	log.Printf("Soba %v, kraj igre, pobednik: %v\n", sobaUUID, pobednik)
}

func sacuvajSobuCassandraDB(soba *Soba, pobednik string, s *gocql.Session, ctx context.Context, rdb *redis.Client) error {
	cqlSobaUUID, err := gocql.ParseUUID(soba.UUID)
	if err != nil {
		return fmt.Errorf("Greška prilikom parsovanja UUID-a sobe: %v\n", err)
	}

	queryDodajSobu := "INSERT INTO sobe (uuid, kod, pobednik, vreme) VALUES (?, ?, ?, ?) IF NOT EXISTS"
	dest := make(map[string]interface{})
	uspesnoDodata, err := s.Query(queryDodajSobu, cqlSobaUUID, soba.Kod, pobednik, time.Now()).MapScanCAS(dest)
	if err != nil {
		return fmt.Errorf("Greška prilikom čuvanja sobe u cassandra bazu podataka: %v\n", err)
	}
	if !uspesnoDodata {
		return fmt.Errorf("Greška prilikom čuvanja sobe u cassadra bazu podataka, soba sa uuid: %v reć postoji\n", soba.UUID)
	}

	for i := range soba.Igraci {
		queryDodajIgracSobe := "INSERT INTO igrac_sobe (username, soba_uuid) VALUES (?, ?) IF NOT EXISTS"
		uspesnoDodata, err = s.Query(queryDodajIgracSobe, soba.Igraci[i].Username, cqlSobaUUID).MapScanCAS(dest)
		if err != nil {
			return fmt.Errorf("Greška prilikom čuvanja sobe igrača %v u cassandra bazu podataka: %v\n", soba.Igraci[i].Username, err)
		}
		if !uspesnoDodata {
			return fmt.Errorf("Greška prilikom čuvanja sobe igrača %v u cassadra bazu podataka, soba sa uuid: %v reć postoji\n", soba.Igraci[i].Username, soba.UUID)
		}
	}

	poteziKljuc := fmt.Sprintf("soba:%s:svi-potezi", soba.UUID)
	porukeKljuc := fmt.Sprintf("soba:%s:sve-poruke", soba.UUID)

	poteziHes, err := rdb.HGetAll(ctx, poteziKljuc).Result()
	for _, potez := range poteziHes {
		var potezPodaci Potez

		if err := json.Unmarshal([]byte(potez), &potezPodaci); err != nil {
			return fmt.Errorf("Greška prilikom konvertovanja hes-a u potez: %v'n", err)
		}

		queryDodajPotez := "INSERT INTO potezi (soba_uuid, vreme, username, indeks_polja) VALUES (?, ?, ?, ?) IF NOT EXISTS"
		dest := make(map[string]interface{})
		uspesnoDodat, err := s.Query(queryDodajPotez, cqlSobaUUID, potezPodaci.Vreme, potezPodaci.IgracUsername, potezPodaci.IndeksPolja).MapScanCAS(dest)
		if err != nil {
			return fmt.Errorf("Greška prilikom čuvanja poteza u cassandra bazu podataka: %v\n", err)
		}
		if !uspesnoDodat {
			return fmt.Errorf("Greška prilikom čuvanja poteza u cassadra bazu podataka, soba sa uuid: %v reć postoji\n")
		}
	}

	porukeHes, err := rdb.HGetAll(ctx, porukeKljuc).Result()
	for _, cetPoruka := range porukeHes {
		var porukaPodaci CetPoruka

		if err := json.Unmarshal([]byte(cetPoruka), &porukaPodaci); err != nil {
			log.Printf("Greška prilikom konvertovanja hes-a u čet poruku: %v\n", err)
			continue
		}

		queryDodajPoruku := "INSERT INTO cet_poruke (soba_uuid, vreme, username, sadrzaj) VALUES (?, ?, ?, ?) IF NOT EXISTS"
		dest := make(map[string]interface{})
		uspesnoDodata, err := s.Query(queryDodajPoruku, cqlSobaUUID, porukaPodaci.Vreme, porukaPodaci.IgracUsername, porukaPodaci.Sadrzaj).MapScanCAS(dest)
		if err != nil {
			log.Printf("Greška prilikom čuvanja čet poruke u cassandra bazu podataka: %v\n", err)
			continue
		}
		if !uspesnoDodata {
			log.Printf("Greška prilikom čuvanja čet poruke u cassadra bazu podataka, soba sa uuid: %v reć postoji\n")
			continue
		}
	}

	return nil
}

func sveSobeIgraca(username string, s *gocql.Session) []byte {
	type potezPodaci struct {
		Vreme time.Time `json:"vreme"`
		Username string `json:"username"`
		IndeksPolja int `json:"indeks_polja"`
	}
	type sobaPodaci struct {
		UUID string `json:"uuid"`
		Kod string `json:"kod"`
		Pobednik string `json:"pobednik"`
		Potezi []potezPodaci  `json:"potezi"`
		Vreme time.Time `json:"vreme"`
	}

	sobe := make([]sobaPodaci, 0)

	guery := "SELECT * FROM igrac_sobe WHERE username = ?"

	scanner := s.Query(guery, username).Iter().Scanner()
	for scanner.Next() {
		var (
			igracUsername string
			sobaUUID gocql.UUID
		)
		if err := scanner.Scan(&igracUsername, &sobaUUID); err != nil {
			log.Printf("Greška prilikom skeniranja podataka sobe igrača iz cassandra baze podataka: %v\n", err)
			continue
		}

		var soba sobaPodaci

		querySoba := "SELECT * FROM sobe WHERE uuid = ?"
		ucitanaSoba := make(map[string]interface{})
		if err := s.Query(querySoba, sobaUUID).MapScan(ucitanaSoba); err != nil {
			log.Printf("Greška prilikom pribavljanja sobe sa UUID: %v\n", sobaUUID.String())
			continue
		}

		soba.UUID = ucitanaSoba["uuid"].(gocql.UUID).String()
		soba.Kod = ucitanaSoba["kod"].(string)
		soba.Pobednik = ucitanaSoba["pobednik"].(string)
		soba.Vreme = ucitanaSoba["vreme"].(time.Time)

		queryPotez := "SELECT vreme, username, indeks_polja FROM potezi WHERE soba_uuid = ?"
		scannerPotez := s.Query(queryPotez, sobaUUID).Iter().Scanner()
		for scannerPotez.Next() {
			var potez potezPodaci

			if err := scannerPotez.Scan(&potez.Vreme, &potez.Username, &potez.IndeksPolja); err != nil {
				log.Printf("Greška prilikom pribavljanja poteza sobe %v: %v\n", sobaUUID.String(), err)
				continue
			}

			soba.Potezi = append(soba.Potezi, potez)
		}

		sobe = append(sobe, soba)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Greška prilikom pribavljanja svih soba igrača iz cassandra baze podataka: %v\n", err)
		return []byte("[]")
	}

	sveSobeJSON, err := json.Marshal(sobe)
	if err != nil {
		log.Print("Greška prilikom marshalovanja svih soba igrača: %v\n", username)
		return []byte("[]")
	}

	return sveSobeJSON
}

func (soba *Soba) Broadcast(tipPoruke string, sadrzajPoruke string) {
	for i := range soba.Igraci {
		soba.Igraci[i].PosaljiOdgovorWS(poruka.NovaPoruka(tipPoruke, sadrzajPoruke).Marshal())
	}
}

func (soba *Soba) Start(ctx context.Context, rdb *redis.Client) {
	soba.IgracNaRedu = soba.Igraci[0].Username
	if err := sacuvajSobuRedisDB(soba, ctx, rdb); err != nil {
		log.Printf("soba.Start, greška prilikom čuvanja sobe u redis bazu: %v\n", err)
		return
	}

	for vreme := range 3 {
		soba.Broadcast("Pocetak_Igre", fmt.Sprintf("Igra počinje za %vs.", 3 - vreme))
		time.Sleep(time.Second)
	}

	redniBrojIgraca := 0
	for i, igrac := range soba.Igraci {
		if igrac.Username == soba.IgracNaRedu {
			redniBrojIgraca = i + 1
		}
	}
	sobaPodaciPoruka := poruka.SobaPodaci(soba.Kod, soba.IgracNaRedu, redniBrojIgraca)
	soba.Broadcast(sobaPodaciPoruka.Tip, sobaPodaciPoruka.Sadrzaj)
	soba.Broadcast("Start", fmt.Sprintf("Igra je počela. Prvi na potezu je igrač: %v", soba.IgracNaRedu))
}
