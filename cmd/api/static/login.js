export class Login
{
    constructor(container, igra)
    {
        this.container = container
        this.igra = igra
    }

    init()
    {
        const usernameInput = document.getElementById('login-input-username')
        const usernameLabel = document.getElementById('login-label-username')
        const greskaUsernamePoruka = document.getElementById('login-greska-username')
        const passwordInput = document.getElementById('login-input-password')
        const passwordLabel = document.getElementById('login-label-password')
        const greskaPasswordPoruka = document.getElementById('login-greska-password')
        const prijaviSeDugme = document.getElementById('login-prijavi-se-dugme')
        const prikaziLozinku = document.getElementById('login-prikazi-lozinku')

        const registrujSeDugme = document.getElementById('login-registruj-se-dugme')
        const izadjiDugme = document.getElementById('login-izadji')

        const signupPrikaz = document.getElementById('signup-prikaz')

        izadjiDugme.onclick = () => {
            this.container.classList.toggle('hidden', true)
            usernameInput.value = ''
            usernameInput.classList.toggle('border-red-600', false)
            greskaUsernamePoruka.classList.toggle('hidden', true)
            passwordInput.value = ''
            passwordInput.classList.toggle('border-red-600', false)
            prikaziLozinku.checked = false
            greskaPasswordPoruka.classList.toggle('hidden', true)
        }

        usernameLabel.onclick = () => {
            usernameInput.focus()
        }

        passwordLabel.onclick = () => {
            passwordInput.focus()
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

        registrujSeDugme.onclick = () => {
            this.container.classList.toggle('hidden', true)
            signupPrikaz.classList.toggle('hidden', false)
        }

        prijaviSeDugme.onclick = () => {
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

            if (usernameInput.value.length == 0 || passwordInput.value.length == 0)
            {
                console.log('nema ses')
                return
            }

            usernameInput.value = ''
            passwordInput.value = ''

            this.container.classList.toggle('hidden', true)
        }
    }
}