export class StartPodesavanja
{
    constructor(igra)
    {
        this.igra = igra
    }

    init(container)
    {   
        const unesiKodSpan = container.querySelector('.unesi-kod-span')
        const nasumicnaIgra = container.querySelector('.nasumicna-igra')
        const unesiKodDialog = document.getElementById('unesi-kod-dialog-wrapper')
        const zatvoriUnosKodaDugme = unesiKodDialog.querySelector('.zatvori-unos-koda-dugme')
        const kodInputs = unesiKodDialog.querySelectorAll('input')
        const zapocniIgruKodDugme = document.getElementById('zapocni-igru-kod-dugme')

        zatvoriUnosKodaDugme.onclick = () => {
            container.classList.add('vidljiv')
            unesiKodDialog.classList.remove('vidljiv')
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
            container.classList.remove('vidljiv')
            unesiKodDialog.classList.add('vidljiv')
        }

        nasumicnaIgra.onclick = () => {
            this.igra.kod = ''
            this.igra.socketInit()
        }
    }
}