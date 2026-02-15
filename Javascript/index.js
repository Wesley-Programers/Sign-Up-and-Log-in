function alreadyHaveAccount() {
    setTimeout(() => {
        window.location.href = '../HTML/logIn.html'
    }, 250);
};


document.addEventListener("DOMContentLoaded", () => {

    let form = document.getElementById("form-sign-up");


    form.addEventListener("submit", async (form) => {
        form.preventDefault();

        const formData = new FormData(form.target)
        let thisNameAlreadyExists = document.getElementById('nameAlreadyExits');
        let thisEmailAlreadyExists = document.getElementById('emailAlreadyExits');
        
        try {

            const fetchAqui = await fetch("http://127.0.0.1:8000/sign", {
                method: "POST",
                body: formData,
                credentials: "include",
            })

            let incorretEmail = document.getElementById('emailIncorrect');
            let shortPassword = document.getElementById('shortPassword');

            const status = fetchAqui.status
            const message = await fetchAqui.text()
            alert(`Status: ${status} Message: ${message}`);

            if (status === 201 && message === "VALID DATA") {
                
                setTimeout(() => {
                    window.location.href = '../HTML/mainAccount.html'
                }, 500);
                
            } else {
                form.preventDefault();

                if (status === 400 && message === "NAME ALREADY EXISTS") {
                    thisNameAlreadyExists.style.display = 'block';  

                } else if (status === 400 && message === "EMAIL ALREADY EXISTS") {
                    thisEmailAlreadyExists.style.display = 'block';

                } else if (status === 400 && message === "SOME ERROR") {

                } else {
                    
                }

            }

        } catch ( error ) {
            console.error("ERROR: ", error);
            alert("SOME ERROR");
        };

    });

});
