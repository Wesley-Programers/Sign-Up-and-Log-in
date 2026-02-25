document.addEventListener("DOMContentLoaded", () => {

    let wrongPassword = document.getElementById("wrongPassword");
    let theSamePassword = document.getElementById("theSamePassword");
    let shortPassword = document.getElementById("shortPassword");
    let passwordConfirmationIsWrong = document.getElementById("passwordConfirmationIsWrong");
    let passwordResetSuccessfully = document.getElementById("passwordResetSuccessfully");

    let leaveResetPassword = document.getElementById("leaveResetPassword");
    let resetForm = document.getElementById("resetForm");

    leaveResetPassword.addEventListener("click", () => {
        window.location.href = './logIn.html';
    });

    resetForm.addEventListener("submit", async (form) => {

        form.preventDefault();
        const formData = new FormData(form.target)

        try {

            const fetchValidToken = await fetch("http://127.0.0.1:8000/valid", {
                method: "POST",
                body: formData,
            })
            const status = fetchValidToken.status
            const message = await fetchValidToken.text()
            alert(`Status: ${status} Message ${message}`)

            if (status === 200 && message === "VALID") {

                const fetchResetPassword = await fetch("http://127.0.0.1:8000/reset/password", {
                    method: "POST",
                    body: formData,
                })

                const status = fetchResetPassword.status
                const message = await fetchResetPassword.text()
                alert(`Status: ${status} Message: ${message}`);

                if (status === 200 && message === "VALID") {
                    passwordResetSuccessfully.style.display = 'block';

                } else if (status === 400 && message === "INCORRECT PASSWORD") {
                    wrongPassword.style.display = 'block';

                } else if (status === 400 && message === "THE PASSWORD ARE THE SAME") {
                    theSamePassword.style.display = 'block';

                } else if (status === 400 && message === "SHORT PASSWORD") {
                    shortPassword.style.display = 'block';

                } else if (status === 400 && message === "PASSWORD CONFIRMATION IS WRONG") {
                    passwordConfirmationIsWrong.style.display = 'block';

                }

            } else {
                alert("SOME UNEXPECTED ERROR OCCURRED");
            }

        } catch (error) {
            console.log("ERROR: ", error);
            alert("SOME ERROR");
        };
    });

});