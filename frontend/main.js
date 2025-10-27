const themeOption = document.querySelectorAll('input[name="themes"]');
themeOption.forEach(themeOption => {
    themeOption.addEventListener("change", (event) => {
        const theme = event.target.value;
        console.log("Escolhido: ", theme);

        localStorage.setItem("currentTheme", theme);
    });
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

    } catch (error) {
        console.error("", error);
    }
});
