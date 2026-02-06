export class Igra
{
    constructor(container)
    {
        this.polja = [[1, 2, 3, 4, 5, 6], 10]
        this.igraContainer = container
        this.igraPrikazContainer = document.getElementById('igra-prikaz')
        this.igraPodaci = document.getElementById('igra-podaci')
        this.startPodesavanjaContainer = document.getElementById('start-podesavanja')
        this.cetContainer = document.getElementById('cet')
        this.sveCetPoruke = this.cetContainer.children[1]
        this.socket = null
        this.kod = ''
        this.igrac = {
            username: "",
            datumRodjenja: new Date()
        }
        this.tabla = this.napraviTablu()

        //this.crtajTablu(this.tabla)
    }

    napraviTablu()
    {
        const tabla = []

        for (let i = 0; i < 21; i++)
        {
            tabla.push({
                indeks: i,
                stanje: 0,
                vrednost: 0
            })
        }

        return tabla
    }

    crtajTablu(tabla, krajIgre = false)
    {
        let idxPolja = 0

        const redCSSKlase = 'p-3 h-full flex flex-row flex-nowrap justify-center '
        const poljeCSSKlaseBase = 'mx-3 border-1 rounded-full aspect-square h-full text-center flex justify-center items-center'
        const hover_border = 'hover:border-2'
        const border_black = 'border-black'
        const border_red = 'border-red-600'
        const border_blue = 'border-blue-600'
        const bg_black = 'bg-black'
        const bg_red = 'bg-red-100'
        const bg_blue = 'bg-blue-100'
        const cursor = 'cursor-pointer'

        this.igraPrikazContainer.innerHTML = ''

        for (let i = 1; i <= 6; i++)
        {
            const red = document.createElement('div')
            red.className = redCSSKlase

            for (let j = 0; j < i; j++)
            {
                const polje = document.createElement('div')
                polje.id = `polje-${idxPolja}`
                
                const broj = document.createElement('span')
                if (tabla[idxPolja].vrednost != 0)
                {
                    broj.innerText = tabla[idxPolja].vrednost
                    broj.className = 'text-[200%] select-none'
                }

                switch (tabla[idxPolja].stanje) {
                    case 0:
                        polje.className = poljeCSSKlaseBase
                        polje.classList.add(border_black, hover_border, cursor)
                        break;
                    case 1:
                        polje.className = poljeCSSKlaseBase
                        polje.classList.add(border_red, bg_red)
                        broj.classList.add('text-red-600')
                        break;
                    case 2:
                        polje.className = poljeCSSKlaseBase
                        polje.classList.add(border_blue, bg_blue)
                        broj.classList.add('text-blue-600')
                        break;
                    case 51:
                        polje.className = poljeCSSKlaseBase
                        polje.classList.add(border_black)
                        broj.classList.add('text-red-600')
                        break;
                    case 52:
                        polje.className = poljeCSSKlaseBase
                        polje.classList.add(border_black)
                        broj.classList.add('text-blue-600')
                        break;
                    case 100:
                        polje.className = poljeCSSKlaseBase
                        polje.classList.add(border_black, bg_black)
                        break;
                    default:
                        break;
                }
                
                if (tabla[idxPolja].stanje != 0)
                    polje.appendChild(broj)
                idxPolja++

                if (!krajIgre)
                {
                    polje.onclick = (e) => {
                        const potez = {
                            vreme: new Date().toISOString(),
                            igrac_username: this.igrac.username,
                            indeks_polja: +e.target.id.slice(6)
                        }
                        
                        this.socket.send(JSON.stringify({ tip: 'Potez', sadrzaj: JSON.stringify(potez) }))
                    }
                }
                red.appendChild(polje)
            }

            this.igraPrikazContainer.appendChild(red)
        }        
    }

    socketInit()
    {
        this.socket = new WebSocket('ws://localhost:8080/ws')
        this.tabla = this.napraviTablu()

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

            this.startPodesavanjaContainer.classList.toggle('hidden', true)
            this.igraPrikazContainer.classList.toggle('hidden', false)
            this.cetContainer.classList.toggle('invisible', false)

            const podaci = {
                kod: this.kod,
                username: this.igrac.username,
                datumRodjenja: this.igrac.datumRodjenja.toISOString()
            }

            this.socket.send(JSON.stringify({ tip: 'Dodaj_U_Sobu', sadrzaj: JSON.stringify(podaci) }))
        }

        this.socket.onmessage = async (e) => {
            const sadrzajPoruke = await e.data.text()
            if (sadrzajPoruke === '')
                return

            const poruka = JSON.parse(await e.data.text())

            const cekanjeIgraca = document.getElementById('cekanje-igraca')

            switch (poruka.tip) {
                case 'Igrac_Podaci':
                    if (this.igrac.username == '')
                    {
                        const podaci = JSON.parse(poruka.sadrzaj)
                        this.igrac.username = podaci.username
                        this.igrac.datumRodjenja = new Date(podaci.datumRodjenja)
                    }
                    break;
                case 'Soba_Podaci':
                    const sobaPodaci = JSON.parse(poruka.sadrzaj)

                    this.igraPodaci.innerHTML = ''

                    const sobaKod = document.createElement('div')
                    sobaKod.className = 'flex flex-row px-2 text-[130%]'
                    
                    const kodLabel = document.createElement('span')
                    kodLabel.innerText = 'Kod:'
                    kodLabel.className = 'text-gray-700 text-left w-full'
                    
                    const kodText = document.createElement('span')
                    kodText.className = 'text-gray-700'
                    kodText.innerText = sobaPodaci.kod

                    sobaKod.appendChild(kodLabel)
                    sobaKod.appendChild(kodText)
                    
                    const sobaPotez = document.createElement('div')
                    sobaPotez.className = 'flex flex-row px-2 text-[130%]'
                    
                    const potezLabel = document.createElement('span')
                    potezLabel.className = 'text-gray-700 text-left w-full'
                    potezLabel.innerText = 'Potez:'
                    
                    const potezText = document.createElement('span')
                    if (sobaPodaci.redni_broj_igraca === 1)
                        potezText.className = 'text-red-600'
                    else if (sobaPodaci.redni_broj_igraca === 2)
                        potezText.className = 'text-blue-600'
                    else
                        potezText.className = 'text-gray-700'
                    potezText.innerText = sobaPodaci.potez

                    sobaPotez.appendChild(potezLabel)
                    sobaPotez.appendChild(potezText)

                    this.igraPodaci.appendChild(sobaKod)
                    this.igraPodaci.appendChild(sobaPotez)
                    break;
                case 'Cekanje':
                    cekanjeIgraca.classList.toggle('hidden', false)
                    break;

                case 'Start':
                    const containerStart = document.createElement('div')
                    containerStart.className = 'flex justify-center'
                    const textStart = document.createElement('div')
                    textStart.innerText = poruka.sadrzaj
                    textStart.className = 'border-1 border-gray-800 bg-gray-500 text-gray-800 text-sm mx-5 rounded-full w-50 mb-1'
                    containerStart.appendChild(textStart)
                    this.sveCetPoruke.appendChild(containerStart)
                    this.crtajTablu(this.tabla)
                    break;

                case 'Pocetak_Igre':
                    cekanjeIgraca.classList.toggle('hidden', true)
                    const containerPocetak = document.createElement('div')
                    containerPocetak.className = 'flex justify-center'
                    const textPocetak = document.createElement('div')
                    textPocetak.innerText = poruka.sadrzaj
                    textPocetak.className = 'border-1 border-gray-800 bg-gray-500 text-gray-800 text-sm mx-5 rounded-full w-50 mb-1'
                    containerPocetak.appendChild(textPocetak)
                    this.sveCetPoruke.appendChild(containerPocetak)
                    break;

                case 'Kraj_Igre':
                    const rezultat = JSON.parse(poruka.sadrzaj)
                    this.tabla = rezultat.tabla
                    zavrsiIgru(rezultat)
                    break;

                case 'Cet_Poruka':
                    const cetPoruka = JSON.parse(poruka.sadrzaj)

                    const novaPoruka = document.createElement('div')
                    novaPoruka.className = 'text-left text-wrap p-1 mx-2 mb-1 text-sm'
                    novaPoruka.innerText = `${cetPoruka.ime_igraca}:    ${cetPoruka.text}`
                    
                    this.sveCetPoruke.appendChild(novaPoruka)
                    break;

                case 'Potez':
                    this.tabla = JSON.parse(poruka.sadrzaj)
                    this.crtajTablu(this.tabla)
                    break;

                default:
                    console.log('default:', poruka)
                    break;
            }
        }

        const zavrsiIgru = (rezultat) => {
            const kraj = true
            this.crtajTablu(this.tabla, kraj)

            const pobednik = document.createElement('div')
            pobednik.className = 'flex justify-center items-center w-full p-3'
            const imePobednika = document.createElement('span')
            imePobednika.className = 'text-[200%]'
            if (rezultat.pobednik === '')
            {
                imePobednika.innerText = 'NEREŠENO!'
            }
            else
            {
                let zbirIgrac1 = 0
                let zbirIgrac2 = 0

                this.tabla.forEach(p => {
                    if (p.stanje === 1)
                        zbirIgrac1 += p.vrednost
                    else if (p.stanje === 2)
                        zbirIgrac2 += p.vrednost 
                })
                
                if (zbirIgrac1 < zbirIgrac2)
                    imePobednika.classList.add('text-red-600')
                else
                    imePobednika.classList.add('text-blue-600')
                imePobednika.innerText = `Pobednik:\n${rezultat.pobednik}`
            }

            pobednik.appendChild(imePobednika)
            this.igraPodaci.appendChild(pobednik)

            const pocetniMeniDugme = document.createElement('div')
            pocetniMeniDugme.innerText = 'Početni meni'
            pocetniMeniDugme.className = 'absolute top-1/12 right-1/12 overflow-hidden p-3 border-2 border-blue-600 rounded-lg text-[xl] text-blue-600 transition-all duration-300 cursor-pointer text-center bg-transparent after:content-[\'\'] after:absolute after:top-1/2 after:left-1/2 after:h-0 after:w-0 after:rounded-full after:bg-blue-600 after:-translate-x-1/2 after:-translate-y-1/2 after:transition-all after:duration-400 after:-z-10 hover:text-white hover:-translate-y-[2px] hover:after:w-screen hover:after:h-screen'
 
            pocetniMeniDugme.onclick = () => {
                this.igraPrikazContainer.classList.toggle('hidden', true)
                this.igraPrikazContainer.innerHTML = ''
                this.startPodesavanjaContainer.classList.toggle('hidden', false)
                const korisnikPodaci = document.getElementById('korisnik-podaci')
                const igraPodaci = document.getElementById('igra-podaci')
                korisnikPodaci.classList.toggle('hidden', false)
                igraPodaci.innerHTML = ''
                igraPodaci.classList.toggle('hidden', true)

                this.tabla = this.napraviTablu()
                this.socket.close()
            }
            
            this.igraPrikazContainer.appendChild(pocetniMeniDugme)
        }

        const cetInit = () => {
            const porukaInput = this.cetContainer.querySelector('input')
            const posaljiDugme = document.getElementById('posalji-poruku-dugme')
            this.sveCetPoruke.innerHTML = ''

            posaljiDugme.onclick = () => {
                if (porukaInput.value !== '')
                {
                    this.socket.send(JSON.stringify({ tip: 'Cet_Poruka', sadrzaj: porukaInput.value }))
                    porukaInput.value = ''
                }
            }

            porukaInput.onkeydown = (e) => {
                if (e.key == 'Enter' && porukaInput.value !== '')
                {
                    this.socket.send(JSON.stringify({ tip: 'Cet_Poruka', sadrzaj: porukaInput.value }))
                    porukaInput.value = ''
                }
                else if (e.key == 'Enter' && porukaInput.value === '')
                {
                    this.socket.send(JSON.stringify({ tip: 'Kraj_Igre', sadrzaj: porukaInput.value }))
                    porukaInput.value = ''
                }
            }
        }
    }
}