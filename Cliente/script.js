
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
    .then(json => {console.log(json)});
}

function sendData2(username, password, accion){
    
    const Url2 = 'http://localhost:8080/' + accion;
    const data2={
        user: username,
        password: password
    }
    axios({
        method: 'post',
        Url: Url2,
        data: {
            data2
        },
        body: JSON.stringify({user: username, password: password})
    })
    .then(data => console.log(data));
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

