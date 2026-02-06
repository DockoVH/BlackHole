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
        const neuspesnaPrijavaLabel = document.getElementById('login-neuspesna-prijava')

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

            usernameInput.value = ''
            usernameInput.classList.toggle('border-red-600', false)
            greskaUsernamePoruka.classList.toggle('hidden', true)
            passwordInput.value = ''
            passwordInput.classList.toggle('border-red-600', false)
            prikaziLozinku.checked = false
            greskaPasswordPoruka.classList.toggle('hidden', true)
            neuspesnaPrijavaLabel.classList.toggle('hidden', true)
        }

        usernameInput.oninput = (e) => {
            usernameInput.value = usernameInput.value.toLowerCase().replace(/[^a-z0-9]/g, '')
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
                return
            }

            const podaci = {
                username: usernameInput.value,
                password: passwordInput.value
            }
            prijaviSe(podaci)
            .then(res => {
                this.igra.igrac = {
                    username: res.igrac.username,
                    datumRodjenja: new Date(res.igrac.datumRodjenja)
                }

                usernameInput.value = ''
                passwordInput.value = ''
                this.container.classList.toggle('hidden', true)
                neuspesnaPrijavaLabel.classList.toggle('hidden', true)

                const startPrijaviSe = document.getElementById('start-prijavi-se')
                const startRegistrujSe = document.getElementById('start-registruj-se')
                const startIgracUsername = document.getElementById('start-username-label')
                const startOdjaviSe = document.getElementById('start-odjavi-se')

                startPrijaviSe.classList.toggle('hidden', true)
                startRegistrujSe.classList.toggle('hidden', true)
                startIgracUsername.classList.toggle('hidden', false)
                startIgracUsername.innerText = this.igra.igrac.username
                startOdjaviSe.classList.toggle('hidden', false)
            })
            .catch(err => {
                console.error(err)
                neuspesnaPrijavaLabel.classList.toggle('hidden', false)
            })
        }

        const prijaviSe = async (podaci) => {
            try {
                const response = await fetch('http://localhost:8080/api/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(podaci)
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
    }
}
