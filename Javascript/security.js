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

let deleteAccountForm = document.getElementById("deleteAccountForm");
let deleteSession = document.getElementById("deleteSession");
let deleteAccount = document.getElementById("deleteAccount");
let incorrectPasswordDeleteAccount = document.getElementById("incorrectPassword");
let incorrectEmailDeleteAccount = document.getElementById("incorrectEmail");
let someErrorDeleteAccount = document.getElementById("someErrorDeleteAccount");


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

deleteAccount.addEventListener("click", () => {
    deleteSession.style.display = 'block';
});


document.addEventListener("DOMContentLoaded", () => {

    let formEmail = document.getElementById("formChangeEmail");
    let currentNameWrongMessage = document.getElementById("currentNameWrongMessage");

    let differentEmailsMessage = document.getElementById("diferentEmailsMessage");
    let incorrectPasswordMessage = document.getElementById("incorrectPasswordMessage");
    let someErrorEmailMessage = document.getElementById("someError");
    let theSameEmail = document.getElementById("newEmailIsTheSame");
    let currentEmailWrong = document.getElementById("currentEmailWrongMessage");

    let form = document.getElementById("formChangeName");
    let someErrorMessage = document.getElementById("someErrorMessage");

    form.addEventListener("submit", async (e) => {
        e.preventDefault();
        const formData = new FormData(e.target);

        try {

            const fetchChangeName = await fetch("http://127.0.0.1:8000/change", {
                method: "POST",
                body: formData,
            })

            const status = fetchChangeName.status
            const message = await fetchChangeName.text()
            alert(`STATUS: ${status} MESSAGE: ${message}`);

            if (status === 200 && message === "VALID DATA") {

                alert("EVERYTHING IS WORK");
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


    formEmail.addEventListener("submit", async (a) => {
        a.preventDefault();
        const formDataEmail = new FormData(a.target);

        try {
            const fetchChangeEmail = await fetch("http://127.0.0.1:8000/email", {
                method: "POST",
                body: formDataEmail,
            })

            const status = fetchChangeEmail.status
            const message = await fetchChangeEmail.text()
            alert(`STATUS: ${status} MESSAGE: ${message}`);

            if (status === 200 && message === "VALID DATA") {
                alert("EVERYTHING OK");
                window.location.reload();

            } else if (status === 400 && message === "EMAIL IS NOT THE SAME") {
                currentEmailWrong.style.display = 'block';

            } else if (status === 400 && message === "NEW EMAIL WRONG") {
                differentEmailsMessage.style.display = 'block';

            } else if (status === 400 && message === "INCORRECT PASSWORD") {
                incorrectPasswordMessage.style.display = 'block';

            } else if (status === 400 && message === "EMAIL IS NOT DIFFERENT") {
                theSameEmail.style.display = 'block';

            } else {
                someErrorEmailMessage.style.display = 'block';

            };

        } catch ( error ) {
            console.error("ERROR: ", error)
        };

    });

    
    deleteAccountForm.addEventListener("submit", async (form) => {
        form.preventDefault();
        const formData = new FormData(form.target);

        try {
            
            const deleteAccountFetch = await fetch("http://127.0.0.1:8000/delete", {
                method: "POST",
                body: formData,
            })

            const status = deleteAccountFetch.status
            const message = await deleteAccountFetch.text()
            alert(`Status: ${status} Message: ${message}`);

            if (status === 200 && message === "VALID DATA") {

            } else if (status === 400 && message === "INCORRECT PASSWORD") {
                incorrectPasswordDeleteAccount.style.display = 'block';

            } else if (status === 400 && message === "INCORRECT EMAIL") {
                incorrectEmailDeleteAccount.style.display = 'block';

            } else {
                someErrorDeleteAccount.style.display = 'block';

            };

        } catch (error) {
            console.error("ERROR: ", error);
            alert("SOME ERROR");
        };
    });
});
