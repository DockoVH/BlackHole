export class Igra
{
    constructor(container)
    {
        this.polja = [[1, 2, 3, 4, 5, 6], 10]
        this.container = container
    }

    crtajPolja()
    {
        this.container.innerHTML = ''
        let idxPolja = 0

        this.polja[0].forEach(p => {
            const red = document.createElement('div')
            red.classList.add('red')

            for (let i = 0; i < p; i++)
            {
                const polje = document.createElement('div')
                polje.classList.add('polje')
                polje.id = `polje-${idxPolja++}`
                polje.onclick = (e) => {
                    poljeIdx = e.target.id
                }

                red.appendChild(polje)
            }

            this.container.appendChild(red)
        })
    }
}