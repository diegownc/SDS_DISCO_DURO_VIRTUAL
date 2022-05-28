function sendData2(tokenUsuario, username){
    var data = new FormData()
    data.append('file', document.getElementById("file").files[0])
    data.append('tokenUsuario', tokenUsuario)
    data.append('username', username)
    
    const url_web = 'http://localhost:8080/upload';
    fetch(url_web, {
        method: 'POST',
        mode: 'cors',
        body: data
    })
    .then(response => {return response.json()})
    .then(obj => {console.log(obj)})
}

document.getElementById("subir")
.addEventListener("click", (evt) =>{
    evt.preventDefault();

    let tokenUsuario = document.getElementById("tokenUsuario").value;
    let username = document.getElementById("usernameLogin").value;

    sendData2(tokenUsuario, username);
})