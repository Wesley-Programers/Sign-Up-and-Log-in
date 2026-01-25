const data = new Date();
const day = String(data.getDate());
const month = String(data.getMonth() + 1);
const year = data.getFullYear();

const hours = String(data.getHours());
const minutes = String(data.getMinutes());
const seconds = String(data.getSeconds());

let inputCurrentEmail = document.getElementById("currentEmail");
let inputNewEmail = document.getElementById("newEmail");
let inputConfirmNewEmail = document.getElementById("confirmNewEmail");
let inputCurrentPassword = document.getElementById("currentPassword");
let inputPasswordDeleteAccount = document.getElementById("passwordConfirm");
let inputEmailDeleteAccount = document.getElementById("emailConfirm");
let theSameName = document.getElementById("newNameIsTheSame");

let mainEmail = document.getElementById("main-email");
let mainName = document.getElementById("main-name");

let changeEmail = document.getElementById("change-email");
let changeName = document.getElementById("change-name");

let change = document.getElementById("div-change-name-window");
let changeEmailWindow = document.getElementById("div-change-email-window");

let leaveChangeName = document.getElementById("leaveChangeName");
let leaveChangeEmail = document.getElementById("leaveChangeEmail");
let leaveDeleteAccount = document.getElementById("leaveDeleteAccount");

let clear = document.getElementById("clear");
let LastLoginAttempt = document.getElementById("logLastLoginAttempt");
let LastEmailChange = document.getElementById("logLastEmailChange");
let lastNameChange = document.getElementById("logLastNameChange");
let lastLogin = document.getElementById("logLastLogin");
let logLastAttemptToChangeTheEmail = document.getElementById("logLastAttemptToChangeTheEmail");
let logLastAttemptToDeleteTheAccount = document.getElementById("logLastAttemptToDeleteTheAccount");

let deleteAccountForm = document.getElementById("deleteAccountForm");
let deleteSession = document.getElementById("deleteSession");
let deleteAccount = document.getElementById("deleteAccount");
let incorrectPasswordDeleteAccount = document.getElementById("incorrectPassword");
let incorrectEmailDeleteAccount = document.getElementById("incorrectEmail");
let someErrorDeleteAccount = document.getElementById("someErrorDeleteAccount");

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


const cookieGet = (name) => {
    const test = document.cookie.match(new RegExp('(^| )' + name + '=([^;]+)'))
    if (test) return test[2];
    return null;
};


changeName.addEventListener("click", () => {
    change.style.display = "block";
});


leaveDeleteAccount.addEventListener("click", () => {
    deleteSession.style.display = 'none';

    inputPasswordDeleteAccount.value = '';
    inputEmailDeleteAccount.value = '';
});


leaveChangeName.addEventListener("click", () => {

    let inputCurrentName = document.getElementById("currentName");
    let inputNewName = document.getElementById("newName");

    inputCurrentName.value = '';
    inputNewName.value = '';

    change.style.display = 'none';
});


changeEmail.addEventListener("click", () => {
    changeEmailWindow.style.display = 'block';
});

leaveChangeEmail.addEventListener("click", () => {

    inputCurrentEmail.value = '';
    inputNewEmail.value = '';
    inputConfirmNewEmail.value = '';
    inputCurrentPassword.value = '';

    changeEmailWindow.style.display = 'none';
});

clear.addEventListener("click", () => {      
    localStorage.removeItem("logs");
    localStorage.removeItem("lastEmail");
    localStorage.removeItem("lastName");
    localStorage.removeItem("lastLogin");
    localStorage.removeItem("lastAttemptToChangeTheEmail");
    localStorage.removeItem("lastAttemptToDeleteTheAccount");
    localStorage.removeItem("lastAttemptToChangeThePassword");

    window.location.reload();

});

deleteAccount.addEventListener("click", () => {
    deleteSession.style.display = 'block';
});


window.addEventListener("DOMContentLoaded", () => {

    console.log("COOKIES: ", document.cookie);

    const cookieBase64 = cookieGet("user_data");
    console.log("", cookieBase64)

    if (cookieBase64) {

        try {

            const idk = atob(cookieBase64)
            const dataWithout = decodeURIComponent(escape(idk));
            const [name, email] = dataWithout.split("|");

            document.getElementById("main-name").innerHTML = name
            document.getElementById("main-email").innerText = email

            console.log("USER NAME: ", name);
            console.log("USER EMAIL: ", email);

        } catch ( error ) {
            console.error("ERROR: ", error);
            // alert("SOME ERROR");
        };
    } else {
        console.log("NOTHING");
    };

});


document.addEventListener("DOMContentLoaded", () => {

    let formEmail = document.getElementById("formChangeEmail");
    let currentNameWrongMessage = document.getElementById("currentNameWrongMessage");

    let differentEmailsMessage = document.getElementById("diferentEmailsMessage");
    let incorrectPasswordMessage = document.getElementById("incorrectPasswordMessage");
    let someErrorEmailMessage = document.getElementById("someError");
    let theSameEmail = document.getElementById("newEmailIsTheSame");
    let currentEmailWrong = document.getElementById("currentEmailWrongMessage");
    let loglastAttemptToChangeThePassword = document.getElementById("loglastAttemptToChangeThePassword");

    let form = document.getElementById("formChangeName");
    let someErrorMessage = document.getElementById("someErrorMessage");

    let lastName = [];
    let lastEmail = [];
    let lastAttemptToChangeTheEmail = [];
    let lastAttemptToDeleteTheAccount = [];

    const testEmail = localStorage.getItem("lastEmail");
    LastEmailChange.innerHTML = "LAST EMAIL CHANGE: " + testEmail

    const testName = localStorage.getItem("lastName");
    lastNameChange.innerHTML = "LAST NAME CHANGE: " + testName

    const logs = localStorage.getItem("logs");
    LastLoginAttempt.innerHTML = 'LAST LOGIN ATTEMPT: ' + logs

    const lastLoginLocalStorage = localStorage.getItem("lastLogin");
    lastLogin.innerHTML = 'LAST LOGIN: ' + lastLoginLocalStorage

    const lastAttemptToChangeTheEmailLocalStorage = localStorage.getItem("lastAttemptToChangeTheEmail");
    logLastAttemptToChangeTheEmail.innerHTML = 'LAST ATTEMPT TO CHANGE THE EMAIL: ' + lastAttemptToChangeTheEmailLocalStorage

    const attemptToDeleteTheAccount = localStorage.getItem("lastAttemptToDeleteTheAccount");
    logLastAttemptToDeleteTheAccount.innerHTML = 'LAST ATTEMPT TO DELETE THE ACCOUNT: ' + attemptToDeleteTheAccount

    const attemptToChangeThePassword = localStorage.getItem("lastAttemptToChangeThePassword")
    loglastAttemptToChangeThePassword.innerHTML = 'LAST ATTEMPT TO CHANGE THE PASSWORD: ' + attemptToChangeThePassword


    form.addEventListener("submit", async (form) => {
        form.preventDefault();
        const formData = new FormData(form.target);

        try {

            const fetchChangeName = await fetch("", {
                method: "POST",
                body: formData,
            })

            const status = fetchChangeName.status
            const message = await fetchChangeName.text()
            alert(`STATUS: ${status} MESSAGE: ${message}`);

            if (status === 200 && message === "NAME VALID") {

                alert("EVERYTHING IS WORK");
                teste("lastName", lastName);
                window.location.reload();

            } else if (status === 400 && message === "INVALID NAME") {
                currentNameWrongMessage.style.display = 'block';
                
            } else if (status === 400 && message === "NAMES IS NOT DIFFERENT") {
                theSameName.style.display = 'block';

            } else {
                someErrorMessage.style.display = 'block';
            };

        } catch ( error ) {
            console.error("ERROR: ", error)
            alert("SOME ERROR");
        };

    });


    formEmail.addEventListener("submit", async (form) => {
        form.preventDefault();
        const formDataEmail = new FormData(form.target);

        try {
            const fetchChangeEmail = await fetch("", {
                method: "POST",
                body: formDataEmail,
            })

            const status = fetchChangeEmail.status
            const message = await fetchChangeEmail.text()
            alert(`STATUS: ${status} MESSAGE: ${message}`);

            if (status === 200 && message === "EVERYTHING VALID") {

                teste("lastEmail", lastEmail);
                alert("EVERYTHING OK");
                window.location.reload();

            } else if (status === 400 && message === "EMAIL IS NOT THE SAME") {
                currentEmailWrong.style.display = 'block';
                teste("lastAttemptToChangeTheEmail", lastAttemptToChangeTheEmail);

            } else if (status === 400 && message === "NEW EMAIL WRONG") {
                differentEmailsMessage.style.display = 'block';
                teste("lastAttemptToChangeTheEmail", lastAttemptToChangeTheEmail);

            } else if (status === 400 && message === "INCORRECT PASSWORD") {
                incorrectPasswordMessage.style.display = 'block';
                teste("lastAttemptToChangeTheEmail", lastAttemptToChangeTheEmail);

            } else if (status === 400 && message === "EMAIL IS NOT DIFFERENT") {
                theSameEmail.style.display = 'block';
                teste("lastAttemptToChangeTheEmail", lastAttemptToChangeTheEmail);

            } else {
                someErrorEmailMessage.style.display = 'block';
                teste("lastAttemptToChangeTheEmail", lastAttemptToChangeTheEmail);

            };

        } catch ( error ) {
            console.error("ERROR: ", error)
        };

    });

    
    deleteAccountForm.addEventListener("submit", async (form) => {
        form.preventDefault();
        const formData = new FormData(form.target);

        try {
            
            const deleteAccountFetch = await fetch("", {
                method: "POST",
                body: formData,
            })

            const status = deleteAccountFetch.status
            const message = await deleteAccountFetch.text()
            alert(`Status: ${status} Message: ${message}`);

            if (status === 200 && message === "EVERYTHING VALID") {
                alert("CONTA SERIA DELETADA AGORA");

            } else if (status === 400 && message === "INCORRECT PASSWORD") {
                incorrectPasswordDeleteAccount.style.display = 'block';
                teste("lastAttemptToDeleteTheAccount", lastAttemptToDeleteTheAccount);

            } else if (status === 400 && message === "INCORRECT EMAIL") {
                incorrectEmailDeleteAccount.style.display = 'block';
                teste("lastAttemptToDeleteTheAccount", lastAttemptToDeleteTheAccount);

            } else {
                someErrorDeleteAccount.style.display = 'block';
                teste("lastAttemptToDeleteTheAccount", lastAttemptToDeleteTheAccount);

            };

        } catch (error) {
            console.error("ERROR: ", error);
            alert("SOME ERROR");
        };
    });
});
