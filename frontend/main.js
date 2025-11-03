let settings = document.getElementById("setting");
let welcome = document.getElementById("welcome");
let account = document.getElementById("account");
let theme = document.getElementById("theme");
let security = document.getElementById("security");
let language = document.getElementById("language");

let darkTheme = document.getElementById("darkTheme");
let lightTheme = document.getElementById("lightTheme");
let dark = document.getElementById("dark");
let light = document.getElementById("light");

let portugueseLanguage = document.getElementById("portugueseLanguage");
let englishLanguage = document.getElementById("englishLanguage");
let portuguese = document.getElementById("portuguese");
let english = document.getElementById("english");


document.querySelectorAll('input[name="languages"]').forEach((input) => {
    
    input.addEventListener("change", (event) => {
        const language = event.target.value;
        console.log("Escolhido: ", language);

        body.classList.remove("english", "portuguese");
        body.classList.add(language);

        localStorage.setItem("currentLanguage", language);
        localStorage.setItem("currentLanguage", language);
        setTimeout(() => {
            window.location.href = window.location.href;
        }, 300);
    });
});


document.querySelectorAll('input[name="themes"]').forEach((input) => {
    
    input.addEventListener("change", (event) => {
        const theme = event.target.value;
        console.log("Escolhido: ", theme);

        body.classList.remove("light", "dark");
        body.classList.add(theme);

        localStorage.setItem("currentTheme", theme);
        localStorage.setItem("currentLanguage", language);
        setTimeout(() => {
            window.location.href = window.location.href;
        }, 300);
    });
});


const getTheme = localStorage.getItem("currentTheme");
const getLanguage = localStorage.getItem("currentLanguage");
const body = document.body;


window.addEventListener("DOMContentLoaded", () => {
    
    console.log("classes atuais: ", body.classList);

    if (getTheme) {
        body.classList.add(getTheme);
        const currentTheme = document.querySelector(`input[name="theme"][value="${getTheme}"]`);
        if (currentTheme) {
            currentTheme.checked = true;
        } else {
            body.classList.add("light");
        };
    };

    if (getLanguage) {
        body.classList.add(getLanguage);
        const currentLanguage = document.querySelector(`input[name="languages"][value="${getLanguage}"]`);
        if (currentLanguage) {
            currentLanguage.checked = true;
        } else {
            body.classList.add("portuguese")
            console.log("portugues adicionado");
        };
    };

    if (body.classList.contains("portuguese")) {
        console.log("usuario prefere portugues");

        settings.innerHTML = '';
        welcome.innerHTML =  '';
        account.innerHTML = '';
        theme.innerHTML = '';
        language.innerHTML = '';
        security.innerHTML = '';
        dark.innerHTML = '';
        light.innerHTML = '';
        english.innerHTML = '';
        portuguese.innerHTML = '';
    } else if (body.classList.contains("english")) {
        console.log("usuario prefere ingles");

        settings.innerHTML = '';
        welcome.innerHTML =  '';
        account.innerHTML = '';
        theme.innerHTML = '';
        language.innerHTML = '';
        security.innerHTML = '';
        dark.innerHTML = '';
        light.innerHTML = '';
        english.innerHTML = '';
        portuguese.innerHTML = '';
    };

});


window.addEventListener("DOMContentLoaded", async () => {

    try {
        const res = await fetch("", {
            credentials: "include"
        })

        if (!res.ok) {
            console.log("nao logado");
        }

        const data = await res.json();
        console.log(data);
        console.log(data.name);
        console.log(data.email);
        welcome.innerHTML = `HELLO, HOW'S IT GOING ${data.name}?`

    } catch (error) {
        console.error("Erro ao carregar usuario: ", error)
    }
})


settings.addEventListener("click", () => {
    settings.style.display = 'none';

    welcome.style.display = 'none';
    account.style.display = 'block';
    theme.style.display = 'block';
    language.style.display = 'block';
    security.style.display = 'block';
});


theme.addEventListener("click", () => {
    settings.style.display = 'none';

    account.style.display = 'none';
    theme.style.display = 'none';
    language.style.display = 'none';
    security.style.display = 'none';


    dark.style.display = 'block';
    light.style.display = 'block';
    darkTheme.style.display = 'block';
    lightTheme.style.display = 'block';
});


account.addEventListener("click", () => {
    account.style.display = 'none';
    theme.style.display = 'none';
    language.style.display = 'none';
    security.style.display = 'none';
});

language.addEventListener("click", () => {

    account.style.display = 'none';
    theme.style.display = 'none';
    language.style.display = 'none';
    security.style.display = 'none';

    portuguese.style.display = 'block';
    english.style.display = 'block';
    portugueseLanguage.style.display = 'block';
    englishLanguage.style.display = 'block';
});
