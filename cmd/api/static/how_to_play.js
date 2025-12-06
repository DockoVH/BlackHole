export class HowToPlay
{
    constructor(brojSlika)
    {
        this.slikaIdx = 0
        this.brojSlika = brojSlika
    }

    init(container)
    {
        const prethodnaSlikaDugme = container.querySelector('.how-to-play-slika-levo')
        const sledecaSlikaDugme = container.querySelector('.how-to-play-slika-desno')
        const howToPlayIzlazDugme = container.querySelector('.how-to-play-izlaz-dugme')
        const howToPlaySlika = container.querySelector('.how-to-play-slika')

        prethodnaSlikaDugme.onclick = () => {
            this.slikaIdx--
            if (this.slikaIdx < 0)
                this.slikaIdx = this.brojSlika - 1
            howToPlaySlika.src = `static/slike/how_to_play_${this.slikaIdx + 1}.png`
        }

        sledecaSlikaDugme.onclick = () => {
            this.slikaIdx = (this.slikaIdx + 1) % this.brojSlika
            howToPlaySlika.src = `static/slike/how_to_play_${this.slikaIdx + 1}.png`
        }

        howToPlayIzlazDugme.onclick = () => {
            container.style.visibility = 'hidden'
            this.slikaIdx = 0
            howToPlaySlika.src = `static/slike/how_to_play_${this.slikaIdx + 1}.png`
        }
    }
}