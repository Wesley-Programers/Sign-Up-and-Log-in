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
