document.addEventListener("DOMContentLoaded", () => {

    function teste(setItem, array) {

        if (month < 10 && day > 10) {

            if (hours < 10 && minutes > 10 && seconds > 10) {
                let newHours = `0${hours}`;
                let newMonth = `0${month}`;
                const logMessage = ` ${day}/${newMonth}/${year} As ${newHours}:${minutes}:${seconds}`;
                array.push(logMessage);
                localStorage.setItem(`${setItem}`, array);

            } else if (hours > 10 && minutes < 10 && seconds > 10) {
                let newMinutes = `0${minutes}`;
                let newMonth = `0${month}`;
                const logMessage = ` ${day}/${newMonth}/${year} As ${hours}:${newMinutes}:${seconds}`;
                array.push(logMessage);
                localStorage.setItem(`${setItem}`, array);

            } else if (hours > 10 && minutes > 10 && seconds < 10) {
                let newSeconds = `0${seconds}`;
                let newMonth = `0${month}`;
                const logMessage = ` ${day}/${newMonth}/${year} As ${hours}:${minutes}:${newSeconds}`;
                array.push(logMessage);
                localStorage.setItem(`${setItem}`, array);

            } else if (hours < 10 && minutes < 10 && seconds > 10) {
                let newHours = `0${hours}`;
                let newMinutes = `0${minutes}`;
                let newMonth = `0${month}`;
                const logMessage = ` ${day}/${newMonth}/${year} As ${newHours}:${newMinutes}:${seconds}`;
                array.push(logMessage);
                localStorage.setItem(`${setItem}`, array);

            } else if (hours > 10 && minutes < 10 && seconds < 10) {
                let newSeconds = `0${seconds}`;
                let newMinutes = `0${minutes}`;
                let newMonth = `0${month}`;
                const logMessage = ` ${day}/${newMonth}/${year} As ${hours}:${newMinutes}:${newSeconds}`;
                array.push(logMessage);
                localStorage.setItem(`${setItem}`, array);

            } else if (hours < 10 && minutes > 10 && seconds < 10) {
                let newHours = `0${hours}`;
                let newSeconds = `0${seconds}`;
                let newMonth = `0${month}`;
                const logMessage = ` ${day}/${newMonth}/${year} As ${newHours}:${minutes}:${newSeconds}`;
                array.push(logMessage);
                localStorage.setItem(`${setItem}`, array);

            } else {
                let newMonth = `0${month}`;
                const logMessage = ` ${day}/${newMonth}/${year} As ${hours}:${minutes}:${seconds}`;
                array.push(logMessage);
                localStorage.setItem(`${setItem}`, array);

            }
            
        } else if (day < 10 && month > 10) {
            
            if (hours < 10 && minutes > 10 && seconds > 10) {
                let newHours = `0${hours}`;
                let newDay = `0${day}`
                const logMessage = ` ${newDay}/${month}/${year} As ${newHours}:${minutes}:${seconds}`;
                array.push(logMessage);
                localStorage.setItem(`${setItem}`, array);

            } else if (hours > 10 && minutes < 10 && seconds > 10) {
                let newMinutes = `0${minutes}`;
                let newDay = `0${day}`
                const logMessage = ` ${newDay}/${month}/${year} As ${hours}:${newMinutes}:${seconds}`;
                array.push(logMessage);
                localStorage.setItem(`${setItem}`, array);

            } else if (hours > 10 && minutes > 10 && seconds < 10) {
                let newSeconds = `0${seconds}`;
                let newDay = `0${day}`
                const logMessage = ` ${newDay}/${month}/${year} As ${hours}:${minutes}:${newSeconds}`;
                array.push(logMessage);
                localStorage.setItem(`${setItem}`, array);

            } else if (hours < 10 && minutes < 10 && seconds > 10) {
                let newHours = `0${hours}`;
                let newMinutes = `0${minutes}`;
                let newDay = `0${day}`
                const logMessage = ` ${newDay}/${month}/${year} As ${newHours}:${newMinutes}:${seconds}`;
                array.push(logMessage);
                localStorage.setItem(`${setItem}`, array);

            } else if (hours > 10 && minutes < 10 && seconds < 10) {
                let newSeconds = `0${seconds}`;
                let newMinutes = `0${minutes}`;
                let newDay = `0${day}`
                const logMessage = ` ${newDay}/${month}/${year} As ${hours}:${newMinutes}:${newSeconds}`;
                array.push(logMessage);
                localStorage.setItem(`${setItem}`, array);

            } else if (hours < 10 && minutes > 10 && seconds < 10) {
                let newHours = `0${hours}`;
                let newSeconds = `0${seconds}`;
                let newDay = `0${day}`
                const logMessage = ` ${newDay}/${month}/${year} As ${newHours}:${minutes}:${newSeconds}`;
                array.push(logMessage);
                localStorage.setItem(`${setItem}`, array);

            } else {
                let newDay = `0${day}`
                const logMessage = ` ${newDay}/${month}/${year} As ${hours}:${minutes}:${seconds}`;
                array.push(logMessage);
                localStorage.setItem(`${setItem}`, array);

            }

        } else if (day > 10 && month > 10) {

            if (hours < 10 && minutes > 10 && seconds > 10) {
                let newHours = `0${hours}`;
                const logMessage = ` ${day}/${month}/${year} As ${newHours}:${minutes}:${seconds}`;
                array.push(logMessage);
                localStorage.setItem(`${setItem}`, array);

            } else if (hours > 10 && minutes < 10 && seconds > 10) {
                let newMinutes = `0${minutes}`;
                const logMessage = ` ${day}/${month}/${year} As ${hours}:${newMinutes}:${seconds}`;
                array.push(logMessage);
                localStorage.setItem(`${setItem}`, array);

            } else if (hours > 10 && minutes > 10 && seconds < 10) {
                let newSeconds = `0${seconds}`;
                const logMessage = ` ${day}/${month}/${year} As ${hours}:${minutes}:${newSeconds}`;
                array.push(logMessage);
                localStorage.setItem(`${setItem}`, setarrayItem);

            } else if (hours < 10 && minutes < 10 && seconds > 10) {
                let newHours = `0${hours}`;
                let newMinutes = `0${minutes}`;
                const logMessage = ` ${day}/${month}/${year} As ${newHours}:${newMinutes}:${seconds}`;
                array.push(logMessage);
                localStorage.setItem(`${setItem}`, array);

            } else if (hours > 10 && minutes < 10 && seconds < 10) {
                let newSeconds = `0${seconds}`;
                let newMinutes = `0${minutes}`;
                const logMessage = ` ${day}/${month}/${year} As ${hours}:${newMinutes}:${newSeconds}`;
                array.push(logMessage);
                localStorage.setItem(`${setItem}`, array);

            } else if (hours < 10 && minutes > 10 && seconds < 10) {
                let newHours = `0${hours}`;
                let newSeconds = `0${seconds}`;
                const logMessage = ` ${day}/${month}/${year} As ${newHours}:${minutes}:${newSeconds}`;
                array.push(logMessage);
                localStorage.setItem(`${setItem}`, array);

            } else {
                const logMessage = ` ${day}/${month}/${year} As ${hours}:${minutes}:${seconds}`;
                array.push(logMessage);
                localStorage.setItem(`${setItem}`, array);
            };

        };
    };

    let logs = [];
    let lastLogin = [];

    const data = new Date();
    const day = String(data.getDate());
    const month = String(data.getMonth() + 1);
    const year = data.getFullYear();

    const hours = String(data.getHours());
    const minutes = String(data.getMinutes());
    const seconds = String(data.getSeconds());


    let form = document.getElementById("");
    let button = document.getElementById("");

    let incorrectName = document.getElementById("");
    let incorrectEmail = document.getElementById("");
    let incorrectPassword = document.getElementById("");

    let dontHaveAccount = document.getElementById("");

    dontHaveAccount.addEventListener("click", () => {
        window.location.href = '../HTML/index.html'
    });
    

    form.addEventListener("submit", async (e) => {
        e.preventDefault();
        const formData = new FormData(e.target)

        try {

            const fetchLogin = await fetch("", {
                method: "POST",
                credentials: "include",
                body: formData
            })

            const status = fetchLogin.status
            const mensagem = await fetchLogin.text()
            alert(`Status: ${status} Mensagem: ${mensagem}`);

            if (status === 200 && mensagem === "") {

                teste("lastLogin", lastLogin);

                console.log("Is everything alright");

                incorrectName.style.display = 'none';
                incorrectEmail.style.display = 'none';
                incorrectPassword.style.display = 'none';

                setTimeout(() => {
                    window.location.href = '../HTML/mainAccount.html'
                }, 500);

            } else {
                e.preventDefault();

                if (status === 409 && mensagem === "") {
                    console.log("");
                    incorrectName.style.display = 'block';

                    teste("logs", logs);

                    incorrectEmail.style.block = 'none';
                    incorrectPassword.style.display = 'none';
                } else if (status === 409 && mensagem === "") {
                    console.log("");
                    incorrectPassword.style.display = 'block';

                    teste("logs", logs);

                    incorrectName.style.display = 'none';
                    incorrectEmail.style.display = 'none';
                } else if (status === 409 && mensagem === "") {
                    console.log("");
                    incorrectEmail.style.display = 'block';

                    teste("logs", logs);
                    
                    incorrectName.style.display = 'none';
                    incorrectPassword.style.display = 'none';
                };
            };

        } catch(error) {
            console.error("ERROR: ", error);
            alert("SOME ERROR");
        };
    });

    
    resetPassword.addEventListener("submit", async (form) => {
        form.preventDefault();
        const formData = new FormData(form.target);

        try {

            const resetFetch = await fetch("", {
                method: "POST",
                body: formData,
                credentials: "include",
            })

            const status = resetFetch.status
            const message = await resetFetch.text()
            alert(`Status: ${status} Message: ${message}`)

            if (status === 200 && message === "") {
                alert("");

            } else if (status === 401 && message === "") {
                alert("");
            }

        } catch (error) {
            console.error("ERROR: ", error);
            alert("SOME ERROR");
        };
        
    });
});
