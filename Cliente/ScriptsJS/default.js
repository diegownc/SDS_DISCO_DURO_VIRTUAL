
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
                getData(document.getElementById("tokenUsuario").value,  document.getElementById("usernameLogin").value)
            }else{
                alert(json.msg)
            }
        }
        console.log(json);});
}

function getData(tokenUsuario, username){ 
    var data = new FormData()
 
    data.append('tokenUsuario', tokenUsuario)
    data.append('username', username)
    
    const url_web = 'http://localhost:8080/nameFiles';
    fetch(url_web, {
        method: 'POST',
        mode: 'cors',
        body: data
    })
    .then(response => {return response.json()})
    .then(obj => {printUserFiles(obj)})
}

 
function printUserFiles( response ) {
    let json = JSON.parse(response);
    fileList = json.msg

    fileList = fileList.replaceAll(")", "")
    fileList = fileList.replaceAll( "(", "")

    console.log(fileList)
    arrayData = fileList.split(",")
    console.log(arrayData)

    var arrayFiles = []
    var arrayIds = [ ]   
    arrayData.forEach( function(entry){
        if( isNaN(entry))
            arrayFiles.push(entry)
        else    
            arrayIds.push(entry)    
    } )

    console.log(arrayFiles)
    console.log(arrayIds)

    var newDiv = document.createElement("div")
    newDiv.setAttribute( 'id'  , 'files ');
    newDiv.setAttribute( 'class' ,  'files')
    for(let i = 0; i  < arrayFiles.length ; i++){
        
        var individualDiv = document.createElement("div")
        individualDiv.setAttribute( 'id'  ,  arrayFiles[i])  

        var newText = document.createTextNode(arrayFiles[i])
        individualDiv.appendChild(newText)

        var newButton = document.createElement("button")
        newButton.setAttribute( 'id'  ,  arrayIds[i]);
        newButton.setAttribute( 'class' ,  'button')
        newButton.textContent = 'Descargar';
        individualDiv.appendChild(newButton)

        newDiv.appendChild(individualDiv)
    }
    var currentDiv = document.getElementById("uploadFile" + '\n')
    document.body.insertBefore(newDiv, currentDiv)
    
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

