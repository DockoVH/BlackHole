export class HowToPlay
{
    constructor(brojSlika)
    {
        this.slikaIdx = 0
        this.brojSlika = brojSlika
        this.uputstvaTekst = [
            'Blackhole je igra čiju tablu predstavljaju polja raspoređena u trougao. Na tabli se nalazi ukupno 21 polje.',
            'Svaki igrač ima svoju boju. Igrači naizmenično zauzimaju polja klikom na njih i označavaju ih svojom bojom.',
            'Prvo zauzeto polje igrača dobija vrednost 1. Svako sledeće zauzeto polje dobija vrednost za 1 veću od prethodnog.',
            'Prvo zauzeto polje igrača dobija vrednost 1. Svako sledeće zauzeto polje dobija vrednost za 1 veću od prethodnog.',
            'Kada na tabli ostane jedno slobodno polje ono postaje crna rupa. Pobednik je onaj igrač koji ima najmanji zbir vrednosti polja oko crne rupe.'
        ]
    }

    init(container)
    {
        const prethodnaSlikaDugme = document.getElementById('how-to-play-slika-levo')
        const sledecaSlikaDugme = document.getElementById('how-to-play-slika-desno')
        const howToPlayIzlazDugme = document.getElementById('how-to-play-izlaz')
        const howToPlaySlika = document.getElementById('how-to-play-slika')
        const uputstvo = document.getElementById('how-to-play-uputstvo')
        uputstvo.innerText = this.uputstvaTekst[0]

        prethodnaSlikaDugme.onclick = () => {
            if (this.slikaIdx === 0)
                return

            this.slikaIdx--
            howToPlaySlika.src = `static/slike/how_to_play_${this.slikaIdx + 1}.png`
            uputstvo.innerText = this.uputstvaTekst[this.slikaIdx]
            sledecaSlikaDugme.classList.toggle('invisible', false)

            if (this.slikaIdx === 0)
                prethodnaSlikaDugme.classList.toggle('invisible', true)
        }

        sledecaSlikaDugme.onclick = () => {
            if (this.slikaIdx === this.brojSlika - 1)
                return
            this.slikaIdx++
            howToPlaySlika.src = `static/slike/how_to_play_${this.slikaIdx + 1}.png`
            uputstvo.innerText = this.uputstvaTekst[this.slikaIdx]
            prethodnaSlikaDugme.classList.toggle('invisible', false)

            if (this.slikaIdx === this.brojSlika - 1)
                sledecaSlikaDugme.classList.toggle('invisible', true)
        }

        howToPlayIzlazDugme.onclick = () => {
            container.classList.toggle('hidden', true)
            this.slikaIdx = 0
            howToPlaySlika.src = `static/slike/how_to_play_${this.slikaIdx + 1}.png`
        }
    }
}