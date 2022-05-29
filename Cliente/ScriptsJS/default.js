
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
       ;});
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
    try{
        //ya existia la tabla...
        document.getElementById("files").value; //provocamos una excepcion
        var tableexists = document.getElementById("files")
        tableexists.remove()
    }catch{
        //Creamos la tabla por primera vez...
        console.log("No existia");
    }
    
    let json = JSON.parse(response);
    fileList = json.msg

    if(fileList != ""){    
        fileList = fileList.replaceAll(")", "")
        fileList = fileList.replaceAll( "(", "")

        
        arrayData = fileList.split(",")
        

        var arrayFiles = []
        var arrayIds = [ ]   
        arrayData.forEach( function(entry){
            if( isNaN(entry))
                arrayFiles.push(entry)
            else    
                arrayIds.push(entry)    
        } )

    

        var table = document.createElement("table")
        table.setAttribute( 'id'  , 'files');
        table.setAttribute( 'class' ,  'files');
        table.setAttribute( 'align' ,  'center');
        
        var tr = document.createElement("tr");
        var td1 = document.createElement("td");
        var newText = document.createElement("label");
        newText.setAttribute( 'class' , 'titleTable');
        newText.textContent = "Nombre del archivo";
        td1.appendChild(newText);

        var td2 = document.createElement("td")
        newText = document.createElement("label")
        newText.setAttribute( 'class' , 'titleTable')
        newText.textContent = "Descargar"
        td2.appendChild(newText);


        var td3 = document.createElement("td")
        newText = document.createElement("label")
        newText.setAttribute( 'class' , 'titleTable')
        newText.textContent = "Eliminar"
        td3.appendChild(newText);


        tr.appendChild(td1);
        tr.appendChild(td2);
        tr.appendChild(td3);
        table.appendChild(tr);
        
        for(let i = 0; i  < arrayFiles.length ; i++){
            tr = document.createElement("tr");
            td1 = document.createElement("td");

            newText = document.createTextNode(arrayFiles[i])
            td1.appendChild(newText);

            td2 = document.createElement("td")

            var newButton = document.createElement("img")
            newButton.setAttribute( 'src'  ,  "img/file.png");
            newButton.setAttribute( 'width'  ,  "30px");
            newButton.setAttribute( 'height'  ,  "30px");        
            newButton.addEventListener("click", function () {
                let tokenUsuario = document.getElementById("tokenUsuario").value;
                let username = document.getElementById("usernameLogin").value;
            
                sendDataDownload(tokenUsuario, username, arrayIds[i]);
            });
            td2.appendChild(newButton);

            
            td3 = document.createElement("td")

            var newButton = document.createElement("img")
            newButton.setAttribute( 'src'  ,  "img/delete.png");
            newButton.setAttribute( 'width'  ,  "30px");
            newButton.setAttribute( 'height'  ,  "30px");        
            newButton.addEventListener("click", function () {
                let tokenUsuario = document.getElementById("tokenUsuario").value;
                let username = document.getElementById("usernameLogin").value;
            
                sendDataDelete(tokenUsuario, username, arrayIds[i]);
                getData(document.getElementById("tokenUsuario").value,  document.getElementById("usernameLogin").value);
            });
            td3.appendChild(newButton);

            tr.appendChild(td1);
            tr.appendChild(td2);
            tr.appendChild(td3);
            
            table.appendChild(tr);
        }
        var currentDiv = document.getElementById("uploadFile" + '\n')
        document.body.insertBefore(table, currentDiv)
    } 
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

