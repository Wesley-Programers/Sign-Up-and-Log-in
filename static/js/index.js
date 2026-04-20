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

        const formData = new FormData(form.target)
        let thisNameAlreadyExists = document.getElementById('nameAlreadyExits');
        let thisEmailAlreadyExists = document.getElementById('emailAlreadyExits');
        
        try {

            const fetchAqui = await fetch("http://127.0.0.1:8000/register", {
                method: "POST",
                body: formData,
                credentials: "include",
            })

            const status = fetchAqui.status
            const message = await fetchAqui.text()
            alert(`Status: ${status} Message: ${message}`);

            if (status === 201 && message === "VALID") {
                setTimeout(() => {
                    window.location.href = '../html/mainAccount.html'
                }, 500);
                
            } else {

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

        } finally {
            button.disabled = false;
        };

    });
});
