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
        const registrujSeGreska = document.getElementById('signup-registruj-se-greska')
        
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
            registrujSeGreska.classList.toggle('hidden', true)
        }

        usernameLabel.onclick = () => {
            usernameInput.focus()
        }

        usernameInput.oninput = (e) => {
            usernameInput.value = usernameInput.value.toLowerCase().replace(/[^a-z0-9]/g, '')
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
            let greska = false

            if (usernameInput.value.length == 0)
            {
                greskaUsernamePoruka.classList.toggle('hidden', false)
                usernameInput.classList.toggle('border-red-600', true)
                greska = true
            }
            else
            {
                greskaUsernamePoruka.classList.toggle('hidden', true)
                usernameInput.classList.toggle('border-red-600', false)
            }
            
            if (!validnaLozinka(passwordInput.value))
            {
                greskaPasswordPoruka.classList.toggle('hidden', false)
                passwordInput.classList.toggle('border-red-600', true)
                greska = true
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
                greska = true
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
                greska = true
            }
            else
            {
                greskaDatumRodjenja.classList.toggle('hidden', true)
                datumRodjenja.classList.toggle('border-red-600', false)

            }

            if (!greska)
            {
                const podaci = {
                    username: usernameInput.value,
                    password: passwordInput.value,
                    datumRodjenja: datumRodjenja.value
                }
                registrujSe(podaci)
                .then(res => {
                    const uspesnoRegistrovanPoruka = document.getElementById('signup-uspesno-poruka')
                    const progressBar = uspesnoRegistrovanPoruka.querySelector('#signup-uspesno-progress-bar')

                    greskaUsernamePoruka.classList.toggle('hidden', true)
                    usernameInput.classList.toggle('border-red-600', false)
                    usernameInput.value = ''
                    greskaPasswordPoruka.classList.toggle('hidden', true)
                    passwordInput.classList.toggle('border-red-600', false)
                    passwordInput.value = ''
                    ponoviGreskaPasswordPoruka.classList.toggle('hidden', true)
                    ponoviPasswordInput.classList.toggle('border-red-600', false)
                    ponoviPasswordInput.value = ''
                    greskaDatumRodjenja.classList.toggle('hidden', true)
                    datumRodjenja.classList.toggle('border-red-600', false)
                    datumRodjenja.value = ''
                    
                    registrujSeGreska.classList.toggle('hidden', false)
                    registrujSeGreska.innerText = ''

                    uspesnoRegistrovanPoruka.classList.toggle('hidden', false)
                    setTimeout(() => {
                        uspesnoRegistrovanPoruka.classList.toggle('opacity-60', true)
                        progressBar.classList.toggle('w-full', true)
                    }, 100)

                    setTimeout(() => {
                        uspesnoRegistrovanPoruka.classList.toggle('hidden', true)
                        uspesnoRegistrovanPoruka.classList.toggle('opacity-60', false)
                        progressBar.classList.toggle('w-full', false)
                    }, 4500)
                })
                .catch(err => {
                    registrujSeGreska.classList.toggle('hidden', false)
                    switch (err.message) {
                        case 'LOZINKA':
                            passwordInput.classList.toggle('border-red-600', true)
                            registrujSeGreska.innerText = 'Nevalidna lozinka!'
                            break;
                        case 'USERNAME':
                            usernameInput.classList.toggle('border-red-600', true)
                            registrujSeGreska.innerText = 'Korisnik sa datim korisničkim imenom već postoji!'
                            break;
                        default:
                            registrujSeGreska.innerText = 'Internal server error!'
                            break;
                    }
                })
            }
        }

        const validnaLozinka = (lozinka) => {
            if (lozinka.length < 8)
                return false
            if (!/[a-z]/.test(lozinka))
                return false
            if (!/[A-Z]/.test(lozinka))
                return false
            if (!/[0-9]/.test(lozinka))
                return false
            if (!/[!@#$%&]/.test(lozinka))
                return false

            return true
        }

        const registrujSe = async (podaci) => {
            try {
                const response = await fetch('http://localhost:8080/api/register', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(podaci)
                })

                const res = await response.json()
                
                if (!response.ok)
                {
                    throw new Error(res.greska)
                }

                return res
            }
            catch(error)
            {
                throw error
            }
        }

        prijaviSeDugme.onclick = () => {
            this.container.classList.toggle('hidden', true)
            loginPrikaz.classList.toggle('hidden', false)

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
            registrujSeGreska.classList.toggle('hidden', true)
        }
    }
}