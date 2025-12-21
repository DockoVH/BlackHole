export class Signup
{
    constructor(container)
    {
        this.container = container
    }

    init()
    {
        const usernameInput = document.getElementById('signup-input-username')
        const usernameLabel = document.getElementById('signup-label-username')
        const greskaUsernamePoruka = document.getElementById('signup-greska-username')
        const passwordInput = document.getElementById('signup-input-password')
        const passwordLabel = document.getElementById('signup-label-password')
        const prikaziLozinku = document.getElementById('signup-prikazi-lozinku')
        const greskaPasswordPoruka = document.getElementById('signup-greska-password')
        const ponoviPasswordInput = document.getElementById('signup-input-password-ponovi')
        const ponoviPasswordLabel = document.getElementById('signup-label-password-ponovi')
        const ponoviPrikaziLozinku = document.getElementById('signup-prikazi-lozinku-ponovi')
        const ponoviGreskaPasswordPoruka = document.getElementById('signup-greska-password-ponovi')
        const datumRodjenja = document.getElementById('signup-datum-rodjenja')
        const greskaDatumRodjenja = document.getElementById('signup-greska-datum-rodjenja')
        
        const registrujSeDugme = document.getElementById('signup-registruj-se-dugme')
        const izadjiDugme = document.getElementById('signup-izadji')

        const prijaviSeDugme = document.getElementById('signup-prijavi-se-dugme')
        const loginPrikaz = document.getElementById('login-prikaz')
        
        izadjiDugme.onclick = () => {
            this.container.classList.toggle('hidden', true)
            usernameInput.value = ''
            usernameInput.classList.toggle('border-red-600', false)
            greskaUsernamePoruka.classList.toggle('hidden', true)
            passwordInput.value = ''
            passwordInput.classList.toggle('border-red-600', false)
            greskaPasswordPoruka.classList.toggle('hidden', true)
            prikaziLozinku.checked = false
            ponoviPasswordInput.value = ''
            ponoviGreskaPasswordPoruka.classList.toggle('border-red-600', false)
            ponoviGreskaPasswordPoruka.classList.toggle('hidden', true)
            ponoviPrikaziLozinku.checked = false
            datumRodjenja.value = ''
            datumRodjenja.classList.toggle('border-red-600', false)
            greskaDatumRodjenja.classList.toggle('hidden', true)
        }

        usernameLabel.onclick = () => {
            usernameInput.focus()
        }

        passwordLabel.onclick = () => {
            passwordInput.focus()
        }

        ponoviPasswordLabel.onclick = () => {
            ponoviPasswordInput.focus()
        }

        prikaziLozinku.onclick = () => {
            if (prikaziLozinku.checked)
            {
                passwordInput.type = 'text'
            }
            else
            {
                passwordInput.type = 'password'
            }
        }

        ponoviPrikaziLozinku.onclick = () => {
            if (ponoviPrikaziLozinku.checked)
            {
                ponoviPasswordInput.type = 'text'
            }
            else
            {
                ponoviPasswordInput.type = 'password'
            }
        }

        registrujSeDugme.onclick = () => {
            if (usernameInput.value.length == 0)
            {
                greskaUsernamePoruka.classList.toggle('hidden', false)
                usernameInput.classList.toggle('border-red-600', true)
            }
            else
            {
                greskaUsernamePoruka.classList.toggle('hidden', true)
                usernameInput.classList.toggle('border-red-600', false)
            }
            
            if (passwordInput.value.length == 0)
            {
                greskaPasswordPoruka.classList.toggle('hidden', false)
                passwordInput.classList.toggle('border-red-600', true)
            }
            else
            {
                greskaPasswordPoruka.classList.toggle('hidden', true)
                passwordInput.classList.toggle('border-red-600', false)
            }

            if (ponoviPasswordInput.value !== passwordInput.value)
            {
                ponoviGreskaPasswordPoruka.classList.toggle('hidden', false)
                ponoviPasswordInput.classList.toggle('border-red-600', true)
            }
            else
            {
                ponoviGreskaPasswordPoruka.classList.toggle('hidden', true)
                ponoviPasswordInput.classList.toggle('border-red-600', false)
            }

            if (datumRodjenja.value === '')
            {
                greskaDatumRodjenja.classList.toggle('hidden', false)
                datumRodjenja.classList.toggle('border-red-600', true)
            }
            else
            {
                greskaDatumRodjenja.classList.toggle('hidden', true)
                datumRodjenja.classList.toggle('border-red-600', false)

            }
        }

        prijaviSeDugme.onclick = () => {
            this.container.classList.toggle('hidden', true)
            loginPrikaz.classList.toggle('hidden', false)
        }
    }
}