import { HowToPlay } from "./how_to_play.js"
import { Igra } from "./igra.js"
import { StartPodesavanja } from "./start_podesavanja.js"

const igraContainer = document.getElementById('igra')
const howToPlayDugme = document.getElementById('how-to-play')
const howToPlayPrikaz = document.getElementById('how-to-play-prikaz')
const startPodesavanjaContainer = document.querySelector('.start-podesavanja')

let howToPlay = new HowToPlay(4)
let igra = new Igra(igraContainer)
let startPodesavanja = new StartPodesavanja(igra)

howToPlay.init(howToPlayPrikaz)
startPodesavanja.init(startPodesavanjaContainer)

howToPlayDugme.onclick = () => {
    howToPlayPrikaz.style.visibility = 'visible'
}
