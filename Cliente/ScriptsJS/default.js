
function sendData(username, password, accion){

    let headers = new Headers();

    headers.append('Content-Type', 'application/json');
    const url_web = 'http://localhost:8080/'+ accion;
    fetch(url_web ,{
        method: 'POST',
        mode: 'cors',
        headers: headers,
        body: JSON.stringify({user: username, password: password})
    })
    .then(response => {return response.json()})
    .then(obj => {
        let json = JSON.parse(obj);
        
        if(accion == "registrar"){
            alert(json.msg);
        }else{
            if(json.result){//LoginOK
                document.getElementById("tokenUsuario").value = json.access_token
                document.getElementById("usernameLogin").value = document.getElementById("username").value
                document.getElementById("inicio").style.display = "none"
                document.getElementById("uploadFile").style.display = "block"
            }else{
                alert(json.msg)
            }
        }
        console.log(json);});
}

document.getElementById("registrar")
.addEventListener("click", (evt) =>{
    evt.preventDefault();

    let username = document.getElementById("username").value;
    let password = document.getElementById("password").value;
    sendData(username, password, "registrar");
})

document.getElementById("login")
.addEventListener("click", (evt) =>{
    evt.preventDefault();

    let username = document.getElementById("username").value;
    let password = document.getElementById("password").value;
    sendData(username, password, "login");
})
