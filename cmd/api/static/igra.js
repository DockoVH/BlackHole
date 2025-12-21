export class Igra
{
    constructor(container)
    {
        this.polja = [[1, 2, 3, 4, 5, 6], 10]
        this.igraContainer = container
        this.igraPrikazContainer = document.getElementById('igra-prikaz')
        this.startPodesavanjaContainer = document.getElementById('start-podesavanja')
        this.cetContainer = document.getElementById('cet')
        this.socket = null
        this.kod = ''

        this.crtajPolja()
    }

    crtajPolja()
    {
        const redCSSKlase = 'p-3 h-full flex flex-row flex-nowrap justify-center '
        const poljeCSSKlase = 'mx-3 border-1 border-black rounded-full aspect-square h-full text-center cursor-pointer hover:border-2'

        let idxPolja = 0

        this.igraPrikazContainer.innerHTML = ''

        this.polja[0].forEach(p => {
            const red = document.createElement('div')
            red.className = redCSSKlase

            for (let i = 0; i < p; i++)
            {
                const polje = document.createElement('div')
                polje.className = poljeCSSKlase
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
            this.startPodesavanjaContainer.classList.toggle('hidden', false)
            this.igraPrikazContainer.classList.toggle('hidden', true)
        }

        this.socket.onclose = () => {
            this.startPodesavanjaContainer.classList.toggle('hidden', false)
            this.igraPrikazContainer.classList.toggle('hidden', true)
            this.cetContainer.classList.toggle('invisible', true)
        }

        this.socket.onopen = () => {
            cetInit()

            this.crtajPolja()

            this.startPodesavanjaContainer.classList.toggle('hidden', true)
            this.igraPrikazContainer.classList.toggle('hidden', false)
            this.cetContainer.classList.toggle('invisible', false)

            this.socket.send(JSON.stringify({Tip: 'Dodaj_U_Sobu', Sadrzaj: this.kod}))
        }

        this.socket.onmessage = async (e) => {
            const poruka = JSON.parse(await e.data.text())
            console.log(poruka)
        }

        const cetInit = () => {
            const porukaInput = this.cetContainer.querySelector('input')
            const posaljiDugme = document.getElementById('posalji-poruku-dugme')

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