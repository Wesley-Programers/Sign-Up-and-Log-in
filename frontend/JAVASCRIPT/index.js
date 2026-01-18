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
                body: formData,
                credentials: "include",
            })

            const emailData = document.getElementById('').value
            let passwordValue = document.getElementById('').value
            let incorretEmail = document.getElementById('')
            let shortPassword = document.getElementById('')

            const status = fetchAqui.status
            const mensagem = await fetchAqui.text()
            alert(`Status: ${status} Message: ${mensagem}`)
            
            let passwordLength = passwordValue.length >= 8;
            let teste = emailData.includes("teste");

            if (teste && passwordLength && status === 201 && mensagem === "") {
                setTimeout(() => {
                    window.location.href = '../HTML/mainAccount.html'
                }, 500);
                
            } else {
                a.preventDefault();
                console.log("something is wrong");

                if (!teste) {
                    incorretEmail.style.display = 'block';
                } 

                if (!passwordLength) {
                    shortPassword.style.display = 'block';
                }

                if (status === 409 && mensagem === "") {
                    thisNameAlreadyExists.style.display = 'block';  
                }

                if (status === 409 && mensagem === "") {
                    thisEmailAlreadyExists.style.display = 'block';
                }

            }

        } catch ( error ) {
            console.error("ERROR: ", error);
            alert("SOME ERROR");
        };

    });

});
