document.addEventListener("DOMContentLoaded", () => {
    let form = document.getElementById("formLogIn");
    let button = document.getElementById("signInButton");

    let incorrectName = document.getElementById("incorrectName");
    let incorrectEmail = document.getElementById("incorrectEmail");
    let incorrectPassword = document.getElementById("incorrectPassword");
    

    form.addEventListener("submit", async (e) => {
        e.preventDefault();
        const formData = new FormData(e.target)

        fetch("/main")
            .then(res => res.json())
            .then(data => {
                alert(data.email)
            });

        try {

            const fetchLogin = await fetch("", {
                method: "POST",
                body: formData
            })

            const status = fetchLogin.status
            const mensagem = await fetchLogin.text()
            alert(`Status: ${status} Mensagem: ${mensagem}`);

            if (status === 200 && mensagem === "Dados de login validos") {

                console.log("Is everything alright");

                incorrectName.style.display = 'none';
                incorrectEmail.style.block = 'none';
                incorrectPassword.style.display = 'none';
                setTimeout(() => {
                    window.location.href = '../HTML/mainAccount.html'
                }, 500);

            } else {
                e.preventDefault();

                if (status === 409 && mensagem === "Nome incorreto") {
                    console.log("nome incorreto");
                    incorrectName.style.display = 'block';

                    incorrectEmail.style.block = 'none';
                    incorrectPassword.style.display = 'none';
                } else if (status === 409 && mensagem === "Senha incorreta") {
                    console.log("senha incorreta");
                    incorrectPassword.style.display = 'block';

                    incorrectName.style.display = 'none';
                    incorrectEmail.style.display = 'none';
                } else if (status === 409 && mensagem === "Email incorreto") {
                    console.log("email incorreto");
                    incorrectEmail.style.display = 'block';

                    incorrectName.style.display = 'none';
                    incorrectPassword.style.display = 'none';
                };
            };

        } catch(error) {
            console.error("ERROR: ", error);
            alert("SOME ERROR");
        };
    });
});


function createAccount() {
    setTimeout(() => {
        location.href = '../HTML/index.html';
    }, 350);
};


function newLanguage() {

    let languages = document.getElementById('languages');
    let language = languages.value;

    let button = document.getElementById('signInButton');
    let portugueseLanguage = document.getElementById('portugueseLanguage');
    let englishLanguage = document.getElementById('englishLanguage');
    let welcomeBack = document.getElementById('welcomeBack');
    let title = document.getElementById('signInTitle');
    let dontHaveAccount = document.getElementById('dontHaveAccount');
    let forgotPassword = document.getElementById('forgotPassword');
    let inputNameEmail = document.getElementById('inputNameEmail');
    let passwordInput = document.getElementById('inputPassword');
    let incorrectName = document.getElementById("incorrectName");
    let incorrectEmail = document.getElementById("incorrectEmail");
    let incorrectPassword = document.getElementById("incorrectPassword");


    if (language === "portuguese") {

        button.innerHTML = 'Entrar';
        portugueseLanguage.innerHTML = 'Português';
        englishLanguage.innerHTML = 'Inglês';
        welcomeBack.innerHTML = 'BEM VINDO DE VOLTA';
        title.innerHTML = 'LOGIN';
        forgotPassword.innerHTML = 'Esqueceu a senha?';
        inputNameEmail.placeholder = 'Nome ou Email';
        passwordInput.placeholder = 'Senha';
        dontHaveAccount.innerHTML = 'Não tem uma conta?';
        incorrectName.innerHTML = 'Nome incorreto';
        incorrectEmail.innerHTML = 'Email incorreto';
        incorrectPassword.innerHTML = 'Senha incorreta';

    } else if (language === "english") {

        button.innerHTML = 'Log in';
        portugueseLanguage.innerHTML = 'Portuguese';
        englishLanguage.innerHTML = 'English';
        welcomeBack.innerHTML = 'WELCOME BACK';
        title.innerHTML = 'LOG IN';
        forgotPassword.innerHTML = 'Forgot the password?';
        inputNameEmail.placeholder = 'Name or Email';
        passwordInput.placeholder = 'Password';
        dontHaveAccount.innerHTML = "Don't have an account?"
        incorrectName.innerHTML = 'Incorrect Name';
        incorrectEmail.innerHTML = 'Incorrect Email';
        incorrectPassword.innerHTML = 'Incorrect Password';

    };

};
