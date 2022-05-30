
function getData(tokenUsuario, username, esversion, idfile){ 
    var data = new FormData()
 
    data.append('tokenUsuario', tokenUsuario);
    data.append('username', username);
    data.append('esversion', esversion);
    data.append('idfile', idfile)
    
    const url_web = 'http://localhost:8080/nameFiles';
    fetch(url_web, {
        method: 'POST',
        mode: 'cors',
        body: data
    })
    .then(response => {return response.json()})
    .then(obj => {
        if (esversion == 0) {
            printUserFiles(obj)
        }else{
            printUserFilesVersiones(obj)
        }
    })
}

 

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
    .then(obj => { 
        let json = JSON.parse(obj);
        alert(json.msg);
        getData(document.getElementById("tokenUsuario").value,  document.getElementById("usernameLogin").value, 0, -1)})
}

function sendDataDownload(tokenUsuario, username, idfile, esversion){
    var data = new FormData()
    data.append('tokenUsuario', tokenUsuario)
    data.append('username', username)
    data.append('idfile', idfile)
    data.append('esversion', esversion)

    const url_web = 'http://localhost:8080/download';
    fetch(url_web, {
        method: 'POST',
        mode: 'cors',
        body: data
    })
    .then(response => {return response.json()})
    .then(obj => {alert(obj)})
}

function sendDataFile(tokenUsuario, username, idfile){
    var data = new FormData()
    data.append('tokenUsuario', tokenUsuario)
    data.append('username', username)
    data.append('idfile', idfile)

    const url_web = 'http://localhost:8080/getFileProperties';
    fetch(url_web, {
        method: 'POST',
        mode: 'cors',
        body: data
    })
    .then(response => {return response.json()})
    .then(obj => {
        document.getElementById("uploadFile").style.display = "none";
        document.getElementById("files").style.display = "none";
        document.getElementById("view_File").style.display = "block";
        document.getElementById("view_size").value = obj.size + " bytes";
        document.getElementById("view_img").src = obj.path;
        document.getElementById("view_txt").value = obj.content;
        let comment = obj.comment.replaceAll('"', '')
        document.getElementById("view_comment").value = comment;
    })
}


function sendDataDelete(tokenUsuario, username, idfile, esversion){
    var data = new FormData()
    data.append('tokenUsuario', tokenUsuario);
    data.append('username', username);
    data.append('idfile', idfile);
    data.append('esversion', esversion);
    
    const url_web = 'http://localhost:8080/delete';
    fetch(url_web, {
        method: 'POST',
        mode: 'cors',
        body: data
    })
    .then(response => {return response.json()})
    .then(obj => {
        if(esversion == 1){
            getData(document.getElementById("tokenUsuario").value,  document.getElementById("usernameLogin").value, 1, idfile);
        }else{
            getData(document.getElementById("tokenUsuario").value,  document.getElementById("usernameLogin").value, 0, -1);
        }
    })
}

function sendDataUpdateComment(tokenUsuario, idfile, comment){
    var data = new FormData()
    data.append('tokenUsuario', tokenUsuario)
    data.append('comment', comment)
    data.append('idfile', idfile)

    const url_web = 'http://localhost:8080/updateComment';
    fetch(url_web, {
        method: 'POST',
        mode: 'cors',
        body: data
    })
    .then(response => {return response.json()})
    .then(obj => {document.getElementById("view_comment").value = comment; console.log(obj)})
}

function printUserFiles( response ) {
    try{
        //ya existia la tabla...
        document.getElementById("filesversiones").value; //provocamos una excepcion
        var tableexists2 = document.getElementById("filesversiones");
        tableexists2.remove();
    }catch{
    }
    try{
        //ya existia la tabla...
        document.getElementById("files").value; //provocamos una excepcion
        var tableexists = document.getElementById("files");
        tableexists.remove();
    }catch{
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

        var divAll = document.createElement("div")
        divAll.setAttribute( 'id'  , 'files');
        
        var table = document.createElement("table")
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


        var td4 = document.createElement("td")
        newText = document.createElement("label")
        newText.setAttribute( 'class' , 'titleTable')
        newText.textContent = "Ver Archivo"
        td4.appendChild(newText);

        tr.appendChild(td1);
        tr.appendChild(td2);
        tr.appendChild(td3);
        tr.appendChild(td4);
        table.appendChild(tr);
        
        for(let i = 0; i  < arrayFiles.length ; i++){
            tr = document.createElement("tr");
            td1 = document.createElement("td");

            newText = document.createTextNode(arrayFiles[i])
            td1.appendChild(newText);

            td2 = document.createElement("td")

            var newButton = document.createElement("img")
            newButton.setAttribute( 'src'  ,  "img/upload.png");
            newButton.setAttribute( 'width'  ,  "30px");
            newButton.setAttribute( 'height'  ,  "30px");        
            newButton.addEventListener("click", function () {
                let tokenUsuario = document.getElementById("tokenUsuario").value;
                let username = document.getElementById("usernameLogin").value;
            
                sendDataDownload(tokenUsuario, username, arrayIds[i], 0);
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
            
                sendDataDelete(tokenUsuario, username, arrayIds[i], 0);
            });
            td3.appendChild(newButton);

            td4 = document.createElement("td")

            var newButton = document.createElement("img")
            newButton.setAttribute( 'src'  ,  "img/file.png");
            newButton.setAttribute( 'width'  ,  "30px");
            newButton.setAttribute( 'height'  ,  "30px");        
            newButton.addEventListener("click", function () {
                let tokenUsuario = document.getElementById("tokenUsuario").value;
                let username = document.getElementById("usernameLogin").value;
                document.getElementById("view_id").value = arrayIds[i];
                
                sendDataFile(tokenUsuario, username, arrayIds[i]);
                document.getElementById("view_filename").value = arrayFiles[i];
                getData(tokenUsuario, username, 1, arrayIds[i]);
            });
            td4.appendChild(newButton);

            
            tr.appendChild(td1);
            tr.appendChild(td2);
            tr.appendChild(td3);
            tr.appendChild(td4);
            

            table.appendChild(tr);
            divAll.appendChild(table);
        }
        var currentDiv = document.getElementById("uploadFile" + '\n')
        document.body.insertBefore(divAll, currentDiv)
    } 
}

function printUserFilesVersiones( response ) {
    try{
        //ya existia la tabla...
        document.getElementById("filesversiones").value; //provocamos una excepcion
        var tableexists = document.getElementById("filesversiones");
        tableexists.remove();
    }catch{
        //Creamos la tabla por primera vez...
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

        var divAll = document.createElement("div")
        divAll.setAttribute( 'id'  , 'filesversiones');
        
        var table = document.createElement("table")
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
            newButton.setAttribute( 'src'  ,  "img/upload.png");
            newButton.setAttribute( 'width'  ,  "30px");
            newButton.setAttribute( 'height'  ,  "30px");        
            newButton.addEventListener("click", function () {
                let tokenUsuario = document.getElementById("tokenUsuario").value;
                let username = document.getElementById("usernameLogin").value;
                
                sendDataDownload(tokenUsuario, username, arrayIds[i], 1);
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
            
                sendDataDelete(tokenUsuario, username, arrayIds[i], 1);
            });
            td3.appendChild(newButton);
            
            tr.appendChild(td1);
            tr.appendChild(td2);
            tr.appendChild(td3);
            
            table.appendChild(tr);
            divAll.appendChild(table);
        }
        var currentDiv = document.getElementById("view_table" + '\n')
        document.body.insertBefore(divAll, currentDiv)
    } 
}

document.getElementById("view_buttontxt")
.addEventListener("click", (evt) =>{
    evt.preventDefault();

    document.getElementById("view_txt").style.display = "block";
    document.getElementById("view_img").style.display = "none";
})

document.getElementById("view_buttonimg")
.addEventListener("click", (evt) =>{
    evt.preventDefault();

    document.getElementById("view_txt").style.display = "none";
    document.getElementById("view_img").style.display = "block";
})

document.getElementById("view_buttonback")
.addEventListener("click", (evt) =>{
    evt.preventDefault();

    document.getElementById("view_File").style.display = "none";
    document.getElementById("uploadFile").style.display = "block";
    document.getElementById("files").style.display = "block";
    try{
        document.getElementById("filesversiones").value; //provocamos una excepcion
        var tableexists = document.getElementById("files");
        tableexists.remove();
    }catch{
        console.log("se ha producido una excepción");
    }
    let tokenUsuario = document.getElementById("tokenUsuario").value;
    let username = document.getElementById("usernameLogin").value;
    getData(tokenUsuario, username, 0, -1);
})

document.getElementById("view_save")
.addEventListener("click", (evt) =>{
    evt.preventDefault();

    let tokenUsuario = document.getElementById("tokenUsuario").value;
    let comment = document.getElementById("view_comment").value;
    let idfile = document.getElementById("view_id").value
    sendDataUpdateComment(tokenUsuario, idfile, comment)
})