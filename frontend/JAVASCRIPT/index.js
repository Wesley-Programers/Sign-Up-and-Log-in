function alreadyHaveAccount() {
    setTimeout(() => {
        window.location.href = '../HTML/logIn.html'
    }, 250);
};


document.addEventListener("DOMContentLoaded", () => {

    let form = document.getElementById("form-sign-up");


    form.addEventListener("submit", async (a) => {
        a.preventDefault();

        const formData = new FormData(a.target)
        let thisNameAlreadyExists = document.getElementById('nameAlreadyExits');
        let thisEmailAlreadyExists = document.getElementById('emailAlreadyExits');
        
        try {

            const fetchAqui = await fetch("", {
                method: "POST",
                body: formData
            })

            const emailData = document.getElementById('inputEmail').value
            let passwordValue = document.getElementById('inputPassword').value
            let incorretEmail = document.getElementById('emailIncorrect')
            let shortPassword = document.getElementById('shortPassword')

            const status = fetchAqui.status
            const mensagem = await fetchAqui.text()
            alert(`Status: ${status} Mensagem: ${mensagem}`)
            
            let passwordLength = passwordValue.length >= 8;
            // let teste = emailData.includes("@email.com");

            if (teste && passwordLength && status === 201 && mensagem === "Dados validos") {
                console.log("EVERYTHING IS ALRIGHT");
                setTimeout(() => {
                    window.location.href = '../HTML/mainAccount.html'
                }, 500);
            } else {
                a.preventDefault();
                console.log("something is wrong");

                if (!teste) {
                    console.log("email nao valido");
                    incorretEmail.style.display = 'block';
                } 

                if (!passwordLength) {
                    console.log("senha curta");
                    shortPassword.style.display = 'block';
                }

                if (status === 409 && mensagem === "Nome já existe") {
                    console.log("nome ja existe")
                    thisNameAlreadyExists.style.display = 'block';  
                }

                if (status === 409 && mensagem === "Email já existe") {
                    console.log("email ja existe")
                    thisEmailAlreadyExists.style.display = 'block';
                }

            }

        } catch ( error ) {
            console.error("ERROR: ", error);
            alert("SOME ERROR");
        };

    });

});


function newLanguage() {
    
    let languagesOption = document.getElementById('languages');
    let languages = languagesOption.value;

    let emailNotValid = document.getElementById('emailIncorrect');
    let shortPassword = document.getElementById('shortPassword');
    let nameAlreadyExits = document.getElementById('nameAlreadyExits');
    let emailAlreadyExits = document.getElementById('emailAlreadyExits');
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
        title.innerHTML = 'CADASTRE-SE';
        withoutAccount.innerHTML = 'Já tem uma conta?';
        send.innerHTML = 'Cadastrar-se';
        portugueseLanguage.innerHTML = 'Português';
        englishLanguage.innerHTML = 'Inglês';
        name.placeholder = 'Nome';
        email.placeholder = 'Email';
        password.placeholder = 'Senha';
        nameAlreadyExits.innerHTML = 'Esse nome já está em uso';
        emailAlreadyExits.innerHTML = 'Esse email já esta em uso';
        shortPassword.innerHTML = 'Sua senha é muito fraca';
        emailNotValid.innerHTML = 'Esse email não é valido';

    } else if (languages === "english") {

        welcome.innerHTML = 'WELCOME';
        title.innerHTML = 'SIGN UP';
        withoutAccount.innerHTML = "Already have an account?";
        send.innerHTML = 'Sign up';
        portugueseLanguage.innerHTML = 'Portuguese';
        englishLanguage.innerHTML = 'English';
        name.placeholder = 'Name';
        email.placeholder = 'Email';
        password.placeholder = 'Password';
        nameAlreadyExits.innerHTML = 'This name already exists';
        emailAlreadyExits.innerHTML = 'This email already exists';
        shortPassword.innerHTML = 'Your password is so weak';
        emailNotValid.innerHTML = 'Enter a valid email adress';

    };
};
