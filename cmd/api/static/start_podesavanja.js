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
        const istorijaIgaraDugme = document.getElementById('istorija-igara-dugme')
        const prijaviSe = document.getElementById('start-prijavi-se')
        const registrujSe = document.getElementById('start-registruj-se')

        const loginPrikaz = document.getElementById('login-prikaz')

        zatvoriUnosKodaDugme.onclick = () => {
            container.classList.toggle('hidden', false)
            unesiKodDialog.classList.toggle('hidden', true)
        }

        kodInputs.forEach((p, idx) => {
            p.oninput = (e) => {
                const vrednost = e.target.value
                e.target.value = vrednost.replace(/[^0-9]/g, '')

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
            this.igra.socketInit()
        }

        unesiKodSpan.onclick = () => {
            container.classList.toggle('hidden', true)
            unesiKodDialog.classList.toggle('hidden', false)
        }

        nasumicnaIgra.onclick = () => {
            this.igra.kod = ''
            this.igra.socketInit()
        }

        istorijaIgaraDugme.onclick = () => {
            loginPrikaz.classList.toggle('hidden', false)
        }

        prijaviSe.onclick = () => {
            const login = document.getElementById('login-prikaz')
            login.classList.toggle('hidden', false)
        }

        registrujSe.onclick = () => {
            const signup = document.getElementById('signup-prikaz')
            signup.classList.toggle('hidden', false)
        }
    }
}