function sendDataUpload(tokenUsuario, username){
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
    .then(obj => { })
}

function sendDataDownload(tokenUsuario, username, idfile){
    var data = new FormData()
    data.append('tokenUsuario', tokenUsuario)
    data.append('username', username)
    data.append('idfile', idfile)

    const url_web = 'http://localhost:8080/download';
    fetch(url_web, {
        method: 'POST',
        mode: 'cors',
        body: data
    })
    .then(response => {return response.json()})
    .then(obj => {alert(obj)})
}

document.getElementById("subir")
.addEventListener("click", (evt) =>{
    evt.preventDefault();

    let tokenUsuario = document.getElementById("tokenUsuario").value;
    let username = document.getElementById("usernameLogin").value;

    sendDataUpload(tokenUsuario, username);
})

 

