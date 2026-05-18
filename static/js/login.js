document.addEventListener("DOMContentLoaded", () => {

    let form = document.getElementById("formLogIn");

    let incorrectName = document.getElementById("incorrectName");
    let incorrectEmail = document.getElementById("incorrectEmail");
    let incorrectPassword = document.getElementById("incorrectPassword");

    let forgotPassword = document.getElementById("forgotPassword");
    let reset = document.getElementById("reset-section");
    let emailForResetPassword = document.getElementById("formEmailForResetPassword");
    let leave = document.getElementById("leave");

    let dontHaveAccount = document.getElementById("dontHaveAccount");

    dontHaveAccount.addEventListener("click", () => {
        window.location.href = '../html/index.html'
    });

    forgotPassword.addEventListener("click", () => {
        reset.style.display = 'block';
    });

    leave.addEventListener("click", () => {
        reset.style.display = 'none';
    });
    

    form.addEventListener("submit", async (e) => {
        let button = document.getElementById("signInButton");

        e.preventDefault();
        button.disabled = true;
        const formData = new FormData(e.target)
        const data = Object.fromEntries(formData.entries());

        try {

            const fetchLogin = await fetch("http://127.0.0.1:8000/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                credentials: "include",
                body: JSON.stringify(data),
            })

            if (!fetchLogin.ok) {
                const errorText = await fetchLogin.text();
                console.log("error: ", errorText);
                return;
            }

            const status = fetchLogin.status
            const message = await fetchLogin.json();
            alert(`Status: ${status} Message: ${message}`);

            if (status === 200) {
                incorrectName.style.display = 'none';
                incorrectEmail.style.display = 'none';
                incorrectPassword.style.display = 'none';

                localStorage.setItem("jwt_key", message.token);

                setTimeout(() => {
                    window.location.href = '../html/mainAccount.html'
                }, 150);

            } else if (status != 200) {

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

        } finally {
            button.disabled = false;
        };
    });

    
    emailForResetPassword.addEventListener("submit", async (form) => {
        let button = document.getElementById("button");
        let link = document.getElementById("link");

        form.preventDefault();
        const formData = new FormData(form.target);
        const data = Object.fromEntries(formData.entries());

        // button.disabled = true;
        try {

            const resetFetch = await fetch("http://127.0.0.1:8000/reset", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                credentials: "include",
                body: JSON.stringify(data),
            })

            const status = resetFetch.status
            const message = await resetFetch.text()
            alert(`Status: ${status} Message: ${message}`);
            window.location.href = message.redirect

        } catch (error) {
            console.error("ERROR: ", error);
            alert("SOME ERROR");

        } finally {
            // button.disabled = false;
        };
    });
});
