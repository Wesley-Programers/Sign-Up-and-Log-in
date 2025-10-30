document.querySelectorAll('input[name="themes"]').forEach((input) => {
    input.addEventListener("change", (event) => {
        const theme = event.target.value;
        console.log("Escolhido: ", theme);

        body.classList.remove("light", "dark");
        body.classList.add(theme);

        localStorage.setItem("currentTheme", theme);
    })
});


const getTheme = localStorage.getItem("currentTheme");
const body = document.body;

window.addEventListener("DOMContentLoaded", () => {

    if (getTheme) {
        body.classList.add(getTheme);
        const currentTheme = document.querySelector(`input[name="theme"][value="${theme}"]`)
        if (currentTheme) {
            currentTheme.checked = true;
        } else {
            body.classList.add("light");
        }
    };

});


window.addEventListener("DOMContentLoaded", async () => {

    try {
        const res = await fetch("", {
            credentials: "include"
        })

        if (!res.ok) {
            console.log("");
        }

        const data = await res.json();
        console.log(data);
        console.log(data.name);
        console.log(data.email);
        welcome.innerHTML = ``

    } catch (error) {
        console.error("", error)
    }
});
