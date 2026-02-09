export function isUserLoggedIn() {
    const token = localStorage.getItem('token');
    return token !== null && token !== '';
}

function redirectHomePage() {
    if(isUserLoggedIn()) {
        window.location.href = '/followed-posts';
    } else {
        window.location.href = '/login';
    }
}

redirectHomePage();