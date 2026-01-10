package poruka

import (
	"encoding/json"
	"log"
	"time"
)

type Poruka struct {
	Tip string `json:"tip"`
	Sadrzaj string `json:"sadrzaj"`
}

func (p Poruka) Marshal() []byte {
	poruka, err := json.Marshal(p)
	if err != nil {
		log.Printf("Poruka.Marshal() greška: %v\n", err)
		return []byte {}
	}
	return poruka
}

func Greska(sadrzaj string) Poruka {
	return Poruka {
		Tip: "Greska",
		Sadrzaj: sadrzaj,
	}
}

func CetPoruka(imeIgraca string, text string) Poruka {
	novaPoruka := struct {
		ImeIgraca string
		Text string
	} {
		ImeIgraca: imeIgraca,
		Text: text,
	}

	sadrzaj, err := json.Marshal(&novaPoruka)
	if err != nil {
		log.Printf("CetPoruka: Marshal() greška: %v\n", err)
		return Greska("Greška prilikom slanja poruke.")
	}

	return Poruka {
		Tip: "Cet_Poruka",
		Sadrzaj: string(sadrzaj),
	}
}

func NovaPoruka(tip string, sadrzaj string) Poruka {
	return Poruka {
		Tip: tip,
		Sadrzaj: sadrzaj,
	}
}

func IgracPodaci(igracUsername string, igracDatumRodjenja time.Time) Poruka {
	podaci := struct {
		Username string
		DatumRodjenja time.Time
	}{
		Username: igracUsername,
		DatumRodjenja: igracDatumRodjenja,
	}

	sadrzaj, err := json.Marshal(&podaci)
	if err != nil {
		log.Printf("IgracPodaci: Marshal() greška: %v\n", err)
		return Greska("Greška prilikom slanja poruke.")
	}

	return Poruka {
		Tip: "Igrac_Podaci",
		Sadrzaj: string(sadrzaj),
	}
}
