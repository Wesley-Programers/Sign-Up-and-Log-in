document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('formThemes');
    const radioButtons = document.querySelectorAll('input[name="theme"]');

    const savedTheme = localStorage.getItem('user-theme');

    if (savedTheme) {
        document.body.classList.add(savedTheme);
        
        const radioToCheck = document.querySelector(`input[value="${savedTheme}"]`);
        if (radioToCheck) radioToCheck.checked = true;
    }
    
    form.addEventListener('change', () => {
        const selectedTheme = document.querySelector('input[name="theme"]:checked').value;

        document.body.classList.remove('light', 'dark');
        document.body.classList.add(selectedTheme);

        localStorage.setItem('user-theme', selectedTheme);
    });
});