document.addEventListener("DOMContentLoaded", () => {

    let form = document.getElementById("form-sign-up");
    let alreadyHaveAnAccount = document.getElementById("withoutAccount");

    alreadyHaveAnAccount.addEventListener("click", () => {
        setTimeout(() => {
            window.location.href = '../html/login.html'
        }, 150);
    });

    form.addEventListener("submit", async (form) => {
        let button = document.getElementById("send");

        form.preventDefault();
        button.disabled = true;

        const formData = new FormData(form.target);
        const data = Object.fromEntries(formData.entries());

        let thisNameAlreadyExists = document.getElementById('nameAlreadyExits');
        let thisEmailAlreadyExists = document.getElementById('emailAlreadyExits');
        
        try {

            const fetchAqui = await fetch("http://127.0.0.1:8000/register", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                credentials: "include",
                body: JSON.stringify(data),
            })

            const status = fetchAqui.status
            const message = await fetchAqui.text()
            alert(`Status: ${status} Message: ${message}`);

            if (status === 201 && message === "success") {
                setTimeout(() => {
                    window.location.href = '../html/mainAccount.html'
                }, 150);
                
            } else {

            }

        } catch ( error ) {
            console.error("ERROR: ", error);
            alert("SOME ERROR");

        } finally {
            button.disabled = false;
        };

    });
});
