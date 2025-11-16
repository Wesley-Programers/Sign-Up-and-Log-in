document.addEventListener("DOMContentLoaded", () => {

});

function createAccount() {
    setTimeout(() => {
        location.href = './index.html';
    }, 250)
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

    };

};
