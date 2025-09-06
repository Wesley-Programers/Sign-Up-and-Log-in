function alreadyHaveAccount() {
    setTimeout(() => {
        window.location.href = './signIn.html'
    }, 150);
};

function sign() {

    let form = document.getElementById("form-sign-up");
        
    form.addEventListener('submit', async function(e) {
        const emailData = document.getElementById('inputEmail').value;


        let thisNameAlreadyExists = document.getElementById('nameAlreadyExits');
        let thisEmailAlreadyExists = document.getElementById('emailAlreadyExits');
        let incorretEmail = document.getElementById('emailIncorrect');
        let shortPassword = document.getElementById('shortPassword');

        let passwordValue = document.getElementById('inputPassword').value;


        if (emailData.includes("@email.com")) {
            // setTimeout(() => {
            //     window.location.href = './signIn.html'
            // }, 600);
        } else {
            incorretEmail.style.display = 'block';
            e.preventDefault();
        };

        let passwordLength = passwordValue.length;
        if (passwordLength < 8) {
        //     setTimeout(() => {
        //         window.location.href = './signIn.html'
        //     }, 600);
        // } else {
            shortPassword.style.display = 'block';
            e.preventDefault();
        };

    });

};



function newLanguage() {
    
    let languagesOption = document.getElementById('languages');
    let languages = languagesOption.value;

    let englishLanguage = document.getElementById('englishLanguage');
    let portugueseLanguage = document.getElementById('portugueseLanguage');
    let welcome = document.getElementById('welcome');
    let title = document.getElementById('sign-up-title');
    let name = document.getElementById('inputName');
    let email = document.getElementById('inputEmail');
    let password = document.getElementById('inputPassword');
    let withoutAccount = document.getElementById('withoutAccount');
    let send = document.getElementById('send');

    if (languages === "portuguese") {

        welcome.innerHTML = 'BEM-VINDO';
        title.innerHTML = 'CADASTRAR-SE';
        withoutAccount.innerHTML = 'Já tem uma conta?';
        send.innerHTML = 'Cadastrar-se';
        portugueseLanguage.innerHTML = 'Português';
        englishLanguage.innerHTML = 'Inglês';
        name.placeholder = 'Nome';
        email.placeholder = 'Email';
        password.placeholder = 'Senha';

    }else if (languages === "english") {

        welcome.innerHTML = 'WELCOME';
        title.innerHTML = 'SIGN UP';
        withoutAccount.innerHTML = "Already have an account?";
        send.innerHTML = 'Sign up';
        portugueseLanguage.innerHTML = 'Portuguese';
        englishLanguage.innerHTML = 'English';
        name.placeholder = 'Name';
        email.placeholder = 'Email';
        password.placeholder = 'Password';

    };
};

</script>
