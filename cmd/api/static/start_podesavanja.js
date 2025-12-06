export class StartPodesavanja
{
    constructor(igra)
    {
        this.igra = igra
    }

    init(container)
    {
        const unesiKodSpan = container.querySelector('.unesi-kod-span')
        const unesiKodInput = container.querySelector('input')
        const nasumicnaIgra = container.querySelector('.nasumicna-igra')

        unesiKodSpan.onclick = () => {
            if (unesiKodSpan.innerText === 'Unesite kod')
            {
                unesiKodSpan.innerText = 'Započni igru'
            }
            else
            {
                unesiKodSpan.innerText = 'Unesite kod'
            }
            unesiKodInput.classList.toggle('vidljiv')
        }

        nasumicnaIgra.onclick = () => {
            container.style.display = 'none'
            this.igra.crtajPolja()
        }
    }
}