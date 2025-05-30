// JavaScript for login page (login.html)
document.addEventListener('DOMContentLoaded', function () {
    const form = document.querySelector('#login-form');
    if (!form) {

        return;
    }

    form.addEventListener('submit', async function (e) {
        e.preventDefault();
        const body = new FormData(form);
        const response = await fetch("/auth/login", {
            method: 'POST',
            body: JSON.stringify(Object.fromEntries(body)),
            headers: {
                'Content-Type': 'application/json',
            },
        });
        if (response.ok) {
            console.log(response)
            // save token to localStorage and set it in the header when redirecting
            const data = await response.json();
            document.cookie = `token=${data.token}; path=/; secure; samesite=strict`;
            // Redirect to the home page on successful login
            window.location.href = '/';
        } else {
            // Handle error response
            const errorText = await response.text();
            const errorDiv = document.createElement('div');
            errorDiv.className = 'error';
            errorDiv.textContent = errorText;
            form.prepend(errorDiv);
        }
        const btn = form.querySelector('button[type="submit"]');
        if (btn) btn.disabled = true;
    });
});
