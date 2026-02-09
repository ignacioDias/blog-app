// import { isUserLoggedIn } from './script.js';

const $loginForm = document.getElementById("login-form");

$loginForm.addEventListener('submit', async function (event) {
    event.preventDefault();

    const formData = new FormData($loginForm);
    const username = formData.get('username');
    const password = formData.get('password');

    try {
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                username: username,
                password: password
            })
        });

        const data = await response.json();

        if (response.ok) {
            localStorage.setItem('token', data.token);
            localStorage.setItem('user', JSON.stringify(data.user));
            
            window.location.href = '/';
        } else {
            alert(data.error || 'Login failed');
        }
    } catch (error) {
        console.error('Login error:', error);
        alert('An error occurred during login');
    }
});