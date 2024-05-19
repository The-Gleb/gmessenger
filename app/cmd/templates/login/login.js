document.getElementById('loginForm').addEventListener('submit', function(event) {
    event.preventDefault();
    
    let formData = {
        "email": document.getElementById("email").value,
        "password": document.getElementById("password").value
        // "login": "login1",
        // "password": "password1"
    }

    fetch("/login", {
        method: 'POST',
        body: JSON.stringify(formData),
        mode: 'cors',
    }).then((response) => {
        if (response.ok) {

            window.location.href = "/chats"
        } else {
            throw 'unauthorized';
        }
    }).catch((e) => { alert(e) });

    return false;
});

document.querySelector('.google-login').addEventListener('click', function() {
    window.location.href = "/auth/google"
});

document.querySelector('.yandex-login').addEventListener('click', function() {
    window.location.href = "/auth/yandex"
});

document.querySelector('.signup').addEventListener('click', function() {
    window.location.href = "/static/register/register.html"
});