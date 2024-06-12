document.getElementById('setNameForm').addEventListener('submit', function(event) {
    event.preventDefault();


    let formData = {
        "username": document.getElementById("username").value,
    }

    fetch("http://localhost:8081/set_username", {

        method: 'POST',
        body: JSON.stringify(formData),
        mode: 'cors',
    }).catch((e) => { alert(e) });

    return false;
});