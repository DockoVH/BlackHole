export class StartPodesavanja
{
    constructor(igra)
    {
        this.igra = igra
    }

    init(container)
    {   
        const unesiKodSpan = document.getElementById('unesi-kod-span')
        const nasumicnaIgra = document.getElementById('nasumicna-igra')
        const unesiKodDialog = document.getElementById('unesi-kod-dialog-wrapper')
        const zatvoriUnosKodaDugme = document.getElementById('zatvori-unos-koda-dugme')
        const kodInputs = unesiKodDialog.querySelectorAll('input')
        const zapocniIgruKodDugme = document.getElementById('zapocni-igru-kod-dugme')
        const istorijaIgaraContainer = document.getElementById('istorija-igara-prikaz')
        const istorijaIgaraIzlazDugme = document.getElementById('istorija-igara-izlaz')
        const istorijaIgaraDugme = document.getElementById('istorija-igara-dugme')
        const prijaviSe = document.getElementById('start-prijavi-se')
        const registrujSe = document.getElementById('start-registruj-se')
        const odjaviSe = document.getElementById('start-odjavi-se')
        const igracUsername = document.getElementById('start-username-label')
        const korisnikPodaci = document.getElementById('korisnik-podaci')
        const igraPodaci = document.getElementById('igra-podaci')
        const prikaziIgruIzlaz = document.getElementById('prikazi-igru-izlaz')
        const prikaziIgruContainer = document.getElementById('prikazi-igru')

        const loginPrikaz = document.getElementById('login-prikaz')

        zatvoriUnosKodaDugme.onclick = () => {
            container.classList.toggle('hidden', false)
            unesiKodDialog.classList.toggle('hidden', true)
        }

        kodInputs.forEach((p, idx) => {
            p.oninput = (e) => {
                const vrednost = e.target.value
                e.target.value = vrednost.replace(/[^a-z0-9]/g, '')

                if (e.target.value)
                {
                    if (idx < kodInputs.length - 1)
                    {
                        kodInputs[idx + 1].focus()
                    }
                }

                ispravanKod()
            }

            p.onkeydown = (e) => {
                if (e.key === 'Backspace' && !e.target.value && idx > 0)
                {
                    kodInputs[idx - 1].focus()
                    kodInputs[idx - 1].value = ''
                    ispravanKod()
                }
            }

            const ispravanKod = () => {
                const popunjenaPolja = Array.from(kodInputs).every(p => p.value !== '')
                zapocniIgruKodDugme.disabled = !popunjenaPolja
            }
        })

        zapocniIgruKodDugme.onclick = () => {
            const kod = Array.from(kodInputs).map(p => p.value).join('')

            this.igra.kod = kod
            korisnikPodaci.classList.toggle('hidden', true)
            igraPodaci.classList.toggle('hidden', false)
            unesiKodDialog.classList.toggle('hidden', true)
            this.igra.socketInit()
        }

        unesiKodSpan.onclick = () => {
            container.classList.toggle('hidden', true)
            unesiKodDialog.classList.toggle('hidden', false)
        }

        nasumicnaIgra.onclick = () => {
            this.igra.kod = ''
            korisnikPodaci.classList.toggle('hidden', true)
            igraPodaci.classList.toggle('hidden', false)
            this.igra.socketInit()
        }

        istorijaIgaraDugme.onclick = () => {
            if (this.igra.igrac.username === '')
            {
                loginPrikaz.classList.toggle('hidden', false)
                return
            }


            container.classList.toggle('hidden', true)
            istorijaIgaraContainer.classList.toggle('hidden', false)

            const statistika = istorijaIgaraContainer.querySelector('#istorija-igara-statistika')

            sveSobe(this.igra.igrac.username)
            .then(res => {
                const sobe = JSON.parse(res.sobe)
                sobe.sort((a, b) => new Date(b.vreme) - new Date(a.vreme))

                const ukunoSoba = sobe.length
                const pobede = sobe.filter(p => p.pobednik === this.igra.igrac.username).length
                const neresene = sobe.filter(p => p.pobednik === '').length

                statistika.innerHTML = `Ukupno: ${ukunoSoba} | <span class="text-green-500">Pobede: ${pobede}</span> | <span class="text-red-600">Porazi: ${ukunoSoba - (pobede + neresene)}</span> | <span class="text-black">Nerešene: ${neresene}</span>`
                
                sobe.forEach(s => {
                    prikaziSobu(s)
                })
            })
            .catch(err => {
                console.error(err)
            })
        }

        prikaziIgruIzlaz.onclick = () => {
            istorijaIgaraContainer.classList.toggle('hidden', false)
            prikaziIgruContainer.classList.toggle('hidden', true)
            igraPodaci.innerHTML = ''
        }

        const prikaziSobu = (s) => {
            const sobaCSSKlase = 'border-1 border-gray-400 rounded-lg mx-3 my-1 text-left ps-15 pe-20 py-2 cursor-pointer bg-gray-200 hover:bg-gray-300 flex justify-between'

            const soba = document.createElement('span')
            soba.className = sobaCSSKlase

            const kod = document.createElement('span')
            kod.innerText = s.kod
            const pobednik = document.createElement('span')
            if (s.pobednik !== '')
            {
                pobednik.innerText = s.pobednik
                if (this.igra.igrac.username === s.pobednik)
                    pobednik.classList.add('text-green-400')
                else
                    pobednik.classList.add('text-red-600')
            }
            else
                pobednik.innerText = 'NEREŠENO'

            soba.appendChild(kod)
            soba.appendChild(pobednik)

            let idxPoteza = 0
            const tabla = this.igra.napraviTablu()

            soba.onclick = () => {
                istorijaIgaraContainer.classList.toggle('hidden', true)
                prikaziIgruContainer.classList.toggle('hidden', false)
                korisnikPodaci.classList.toggle('hidden', true)
                igraPodaci.classList.toggle('hidden', false)
                igraPodaci.innerHTML = ''

                const igraci = document.createElement('div')
                igraci.className = 'flex flex-col p-2 text-gray-500 text-[150%]'
                igraci.innerText = 'Igrači:'
                
                const igrac1 = document.createElement('div')
                igrac1.className = 'text-left text-red-600 ps-5'
                igrac1.innerText = `${s.potezi[0].username}`

                const igrac2 = document.createElement('div')
                igrac2.className = 'text-left text-blue-600 ps-5'
                igrac2.innerText = `${s.potezi[1].username}`

                igraci.appendChild(igrac1)
                igraci.appendChild(igrac2)
                igraPodaci.appendChild(igraci)

                crtajPrikazTable(prikaziIgruContainer, tabla)

                const levo = document.getElementById('prikazi-igru-prethodni-potez')
                const desno = document.getElementById('prikazi-igru-sledeci-potez')

                levo.onclick = () => {
                    if (idxPoteza <= 0)
                        return
                    else if (idxPoteza === s.potezi.length)
                    {
                        tabla.forEach((p, i) => {
                            if (p.stanje === 100)
                            {
                                tabla[i].stanje = 0
                                tabla[i].vrednost = 0
                            }
                            else if (p.stanje > 50)
                                tabla[i].stanje -= 50
                        })

                        igraPodaci.children[igraPodaci.children.length - 1].remove()
                    }

                    idxPoteza--
                    tabla[s.potezi[idxPoteza].indeks_polja].stanje = 0
                    tabla[s.potezi[idxPoteza].indeks_polja].vrednost = 0

                    crtajPrikazTable(prikaziIgruContainer, tabla)
                }

                desno.onclick = () => {
                    if (idxPoteza === s.potezi.length)
                        return
                    
                    tabla[s.potezi[idxPoteza].indeks_polja].stanje = (idxPoteza % 2) + 1
                    tabla[s.potezi[idxPoteza].indeks_polja].vrednost = Math.trunc(idxPoteza / 2) + 1
                    idxPoteza++

                    if (idxPoteza === s.potezi.length)
                    {
                        const odrediRedPolja = (idx) => Math.floor((Math.sqrt(8 * idx + 1) - 1) / 2) + 1

                        let idxCrnaRupa = 0
                        tabla.forEach((p, i) => {
                            if (p.stanje === 0)
                            {
                                tabla[i].stanje = 100
                                idxCrnaRupa = i
                            }
                            else if (p.stanje === 1 || p.stanje === 2)
                                tabla[i].stanje += 50

                        })

                        const brojPolja = tabla.length
                        const crnaRupaRed = odrediRedPolja(idxCrnaRupa)
                        const kandidati = [
                            { idx: idxCrnaRupa - crnaRupaRed, red: crnaRupaRed - 1 },
                            { idx: idxCrnaRupa - crnaRupaRed + 1, red: crnaRupaRed - 1 },
                            { idx: idxCrnaRupa - 1, red: crnaRupaRed },
                            { idx: idxCrnaRupa + 1, red: crnaRupaRed },
                            { idx: idxCrnaRupa + crnaRupaRed, red: crnaRupaRed + 1 },
                            { idx: idxCrnaRupa + crnaRupaRed + 1, red: crnaRupaRed + 1 },
                        ]

                        for (const kandidat of kandidati)
                        {
                            if (kandidat.idx >= 0 && kandidat.idx < brojPolja && kandidat.red === odrediRedPolja(kandidat.idx))
                                tabla[kandidat.idx].stanje -= 50
                        }

                        const pobednikPrikaz = document.createElement('div')
                        pobednikPrikaz.className = 'flex flex-col justify-center items-center w-full p-3 text-gray-500 text-[150%]'
                        pobednikPrikaz.innerText = s.pobednik ? 'Pobednik:' : 'NEREŠENO' 

                        if (s.pobednik)
                        {
                            const imePobednika = document.createElement('div')
                            imePobednika.classList.add('text-[150%]', s.pobednik === s.potezi[0].username ? 'text-red-600' : 'text-blue-600')
                            imePobednika.innerText = s.pobednik
                            
                            pobednikPrikaz.appendChild(imePobednika)
                        }
                        igraPodaci.appendChild(pobednikPrikaz)
                    }
                    crtajPrikazTable(prikaziIgruContainer, tabla)
                }
            }

            istorijaIgaraContainer.appendChild(soba)
        }

        const crtajPrikazTable = (container, tabla) => {
            let idxPolja = 0

            const redCSSKlase = 'p-3 h-full flex flex-row flex-nowrap justify-center '
            const poljeCSSKlaseBase = 'mx-3 border-1 rounded-full aspect-square h-full text-center flex justify-center items-center'
            const border_black = 'border-black'
            const border_red = 'border-red-600'
            const border_blue = 'border-blue-600'
            const bg_black = 'bg-black'
            const bg_red = 'bg-red-100'
            const bg_blue = 'bg-blue-100'

            const redovi = Array.from(container.children)
            redovi.forEach(r => {
                if (!r.hasAttribute('ostaje'))
                    r.remove()
            })

            const poslednjiRed = document.getElementById('prikazi-igru-prethodni-potez').parentElement.parentElement

            for (let i = 1; i <= 6; i++)
            {
                const red = document.createElement('div')
                red.className = redCSSKlase

                for (let j = 0; j < i; j++)
                {
                    const polje = document.createElement('div')
                    
                    const broj = document.createElement('span')
                    if (tabla[idxPolja].vrednost != 0)
                    {
                        broj.innerText = tabla[idxPolja].vrednost
                        broj.className = 'text-[200%] select-none'
                    }

                    switch (tabla[idxPolja].stanje) {
                        case 0:
                            polje.className = poljeCSSKlaseBase
                            polje.classList.add(border_black)
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
                    red.appendChild(polje)
                }

                container.insertBefore(red, poslednjiRed)
            }        
        }

        const sveSobe = async (username) => {
            try {
                const response = await fetch(`http://localhost:8080/api/igrac/sveSobe?username=${username}`, {
                    method: 'GET'
                })

                if (!response.ok)
                {
                    throw new Error(`Greška: ${response.status}`)
                }

                return await response.json()
            }
            catch(error)
            {
                throw error
            }
        }

        istorijaIgaraIzlazDugme.onclick = () => {
            container.classList.toggle('hidden', false)
            istorijaIgaraContainer.classList.toggle('hidden', true)
            igraPodaci.classList.toggle('hidden', true)
            korisnikPodaci.classList.toggle('hidden', false)
            const redovi = Array.from(istorijaIgaraContainer.children)
            redovi.forEach(r => {
                if (!r.hasAttribute('ostaje'))
                    r.remove()
            })
        }

        prijaviSe.onclick = () => {
            const login = document.getElementById('login-prikaz')
            login.classList.toggle('hidden', false)
        }

        registrujSe.onclick = () => {
            const signup = document.getElementById('signup-prikaz')
            signup.classList.toggle('hidden', false)
        }

        odjaviSe.onclick = () => {
            this.igra.igrac = {
                username: "",
                datumRodjenja: Date.now()
            }
            prijaviSe.classList.toggle('hidden', false)
            registrujSe.classList.toggle('hidden', false)
            igracUsername.classList.toggle('hidden', true)
            odjaviSe.classList.toggle('hidden', true)
        }
    }
}