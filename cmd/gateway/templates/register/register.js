document.getElementById('registerForm').addEventListener('submit', function(event) {
    event.preventDefault();

    let formData = {
        "username":  document.getElementById('username').value,
        "email": document.getElementById("email").value,
        "password": document.getElementById("password").value
        // "login": "login1",
        // "password": "password1"
    }

    let confirmPassword = document.getElementById('confirmPassword').value;

    if (formData.password !== confirmPassword) {
        alert("Passwords do not match");
        return;
    }

    fetch('/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(formData)
    })
    .then(response => {
        if (response.ok) {
            window.location.href = "/chats"
        } else if (response.status == 409) {
            alert("user with such email already exists");
            return;
        } else {
            throw new Error('Network response was not ok');
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });

});

document.querySelector('.login').addEventListener('click', function() {
    window.location.href = "/static/login/login.html"
});

document.querySelector('.google-login').addEventListener('click', function() {
    window.location.href = "/auth/google"
});

document.querySelector('.yandex-login').addEventListener('click', function() {
    window.location.href = "/auth/yandex"
});