

document.getElementById("subir")
.addEventListener("click", (evt) =>{
    evt.preventDefault();

    alert(document.getElementById('tokenUsuario').value);
})