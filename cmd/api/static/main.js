import { HowToPlay } from "./how_to_play.js"
import { Igra } from "./igra.js"
import { Login } from "./login.js"
import { Signup } from "./signup.js"
import { StartPodesavanja } from "./start_podesavanja.js"

const igraContainer = document.getElementById('igra')
const howToPlayDugme = document.getElementById('how-to-play')
const howToPlayPrikaz = document.getElementById('how-to-play-prikaz')
const startPodesavanjaContainer = document.getElementById('start-podesavanja')
const loginPrikaz = document.getElementById('login-prikaz')
const signupPrikaz = document.getElementById('signup-prikaz')

let howToPlay = new HowToPlay(5)
let igra = new Igra(igraContainer)
let startPodesavanja = new StartPodesavanja(igra)
let login = new Login(loginPrikaz, igra)
let signup = new Signup(signupPrikaz)

howToPlay.init(howToPlayPrikaz)
startPodesavanja.init(startPodesavanjaContainer)
login.init()
signup.init()

howToPlayDugme.onclick = () => {
    howToPlayPrikaz.classList.toggle('hidden', false)
}
