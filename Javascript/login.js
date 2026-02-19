document.addEventListener("DOMContentLoaded", () => {

    let form = document.getElementById("formLogIn");

    let incorrectName = document.getElementById("incorrectName");
    let incorrectEmail = document.getElementById("incorrectEmail");
    let incorrectPassword = document.getElementById("incorrectPassword");

    let forgotPassword = document.getElementById("forgotPassword");
    let inputEmailForResetPassword = document.getElementById("inputEmailForResetPassword");
    let reset = document.getElementById("reset");
    let formEmailForResetPassword = document.getElementById("formEmailForResetPassword");
    let emailForReset = document.getElementById("emailForReset");
    let leave = document.getElementById("leave");

    let invalidEmail = document.getElementById("invalidEmail");
    let link = document.getElementById("link");

    let dontHaveAccount = document.getElementById("dontHaveAccount");

    dontHaveAccount.addEventListener("click", () => {
        window.location.href = '../HTML/index.html'
    });

    forgotPassword.addEventListener("click", () => {
        reset.style.display = 'block';
        emailForReset.style.display = 'block';
        // passwordForReset.style.display = 'none';
    });

    leave.addEventListener("click", () => {
        reset.style.display = 'none';
        inputEmailForResetPassword.value = '';
    });
    

    form.addEventListener("submit", async (e) => {
        e.preventDefault();
        const formData = new FormData(e.target)

        try {

            const fetchLogin = await fetch("http://127.0.0.1:8000/login", {
                method: "POST",
                credentials: "include",
                body: formData
            })

            const status = fetchLogin.status
            const message = await fetchLogin.text()
            alert(`Status: ${status} Message: ${message}`);

            if (status === 200 && message === "VALID DATA") {
                incorrectName.style.display = 'none';
                incorrectEmail.style.display = 'none';
                incorrectPassword.style.display = 'none';

                setTimeout(() => {
                    window.location.href = '../HTML/mainAccount.html'
                }, 300);

            } else {
                e.preventDefault();

                if (status === 409 && message === "WRONG EMAIL OR NAME") {
                    incorrectName.style.display = 'block';
                    incorrectEmail.style.block = 'none';
                    incorrectPassword.style.display = 'none';

                } else if (status === 409 && message === "WRONG PASSWORD") {
                    incorrectPassword.style.display = 'block';
                    incorrectName.style.display = 'none';
                    incorrectEmail.style.display = 'none'

                } else if (status === 409 && message === "WRONG EMAIL") {
                    incorrectEmail.style.display = 'block';
                    incorrectName.style.display = 'none';
                    incorrectPassword.style.display = 'none';

                } else {
                    alert("RANDOM ERROR");
                }
            };

        } catch (error) {
            console.error("ERROR: ", error);
            alert("SOME ERROR");
        };

    });

    
    formEmailForResetPassword.addEventListener("submit", async (form) => {
        form.preventDefault();
        const formData = new FormData(form.target);

        try {

            const resetFetch = await fetch("http://127.0.0.1:8000/reset", {
                method: "POST",
                body: formData,
            })

            const status = resetFetch.status
            const message = await resetFetch.text()
            alert(`Status: ${status} Message: ${message}`);

            if (status === 200) {

            } else if (status === 400 && anotherTest === "INVALID EMAIL") {
                alert("INVALID EMAIL");
                invalidEmail.style.display = 'block';
            }

        } catch (error) {
            console.error("ERROR: ", error);
            alert("SOME ERROR");
        };

    });

});
