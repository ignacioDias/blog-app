const $registerForm = document.getElementById("register-form");

$registerForm.addEventListener('submit', async function (event) {
    event.preventDefault();

    const formData = new FormData($registerForm);
    const username = formData.get('username');
    const email = formData.get('email');
    const password = formData.get('password');

    try {
        const response = await fetch('/api/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                username: username,
                email: email,
                password: password
            })
        });

        const data = await response.json();

        if (response.ok) {
            alert('Registration successful! Please login.');
            window.location.href = '/login';
        } else {
            alert(data.error || 'Registration failed');
        }
    } catch (error) {
        console.error('Register error:', error);
        alert('An error occurred during registration');
    }
});