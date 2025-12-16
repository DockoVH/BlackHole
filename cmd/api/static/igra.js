export class Igra
{
    constructor(container)
    {
        this.polja = [[1, 2, 3, 4, 5, 6], 10]
        this.igraContainer = container
        this.igraPrikazContainer = container.querySelector('.igra-prikaz')
        this.startPodesavanjaContainer = container.querySelector('.start-podesavanja')
        this.cetContainer = document.getElementById('cet')
        this.socket = null
        this.kod = ''
    }

    crtajPolja()
    {
        let idxPolja = 0

        this.igraPrikazContainer.classList.add('vidljiv')
        this.igraPrikazContainer.innerHTML = ''

        this.polja[0].forEach(p => {
            const red = document.createElement('div')
            red.classList.add('red')

            for (let i = 0; i < p; i++)
            {
                const polje = document.createElement('div')
                polje.classList.add('polje')
                polje.id = `polje-${idxPolja++}`
                polje.onclick = (e) => {
                    const poljeIdx = e.target.id
                    this.socket.send(`Kliknuto polje sa indeksom ${poljeIdx}`)
                }

                red.appendChild(polje)
            }

            this.igraPrikazContainer.appendChild(red)
        })
    }

    socketInit()
    {
        this.socket = new WebSocket('ws://localhost:8080/ws')

        this.socket.onerror = () => {
            this.startPodesavanjaContainer.classList.add('vidljiv')
            this.igraPrikazContainer.classList.remove('vidljiv')
        }

        this.socket.onclose = () => {
            this.startPodesavanjaContainer.classList.add('vidljiv')
            this.igraPrikazContainer.classList.remove('vidljiv')
            this.cetContainer.classList.remove('vidljiv')
        }

        this.socket.onopen = () => {
            cetInit()

            this.crtajPolja()

            this.startPodesavanjaContainer.classList.remove('vidljiv')
            this.igraPrikazContainer.classList.add('vidljiv')
            this.cetContainer.classList.add('vidljiv')

            this.socket.send(JSON.stringify({Tip: 'Dodaj_U_Sobu', Sadrzaj: this.kod}))
        }

        this.socket.onmessage = async (e) => {
            const poruka = JSON.parse(await e.data.text())
            console.log(poruka)
        }

        const cetInit = () => {
            const porukaInput = this.cetContainer.querySelector('input')
            const posaljiDugme = this.cetContainer.querySelector('.posalji-poruku-dugme')

            posaljiDugme.onclick = () => {
                if (porukaInput.value !== '')
                {
                    this.socket.send(JSON.stringify({Tip: 'Cet_Poruka', Sadrzaj: porukaInput.value}))
                    porukaInput.value = ''
                }
            }

            porukaInput.onkeydown = (e) => {
                if (e.key == 'Enter' && porukaInput.value !== '')
                {
                    this.socket.send(JSON.stringify({Tip: 'Cet_Poruka', Sadrzaj: porukaInput.value}))
                    porukaInput.value = ''
                }
            }
        }
    }
}