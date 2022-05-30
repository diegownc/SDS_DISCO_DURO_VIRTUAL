package api

import (
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"strconv"

	db "github.com/diegownc/SDS_DISCO_DURO_VIRTUAL/db"
	"github.com/gin-gonic/gin"
)

type uploadResponse struct {
	Result bool   `json:"result"`
	Msg    string `json:"msg"`
}

func (server *Server) uploadFile(ctx *gin.Context) {

	username := ctx.Request.PostFormValue("username")

	//Obtenemos nuestra clave privada
	clavePrivada := leerClavePrivada()

	//Convertimos el cotenido recibido de base64 a bytes[]
	usernameCifradoBytes, err := base64.Encoding.Strict(*base64.StdEncoding).DecodeString(username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Desencriptamos con nuestra clave privada...
	usernameDescifrado, err := RsaDecrypt(usernameCifradoBytes, []byte(clavePrivada))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	idfolder := db.ObtenerIdFolder(string(usernameDescifrado))

	//Almacenamos en base de datos el nombre del archivo..
	res, versionado, version := db.RegistrarArchivo(file.Filename, "", idfolder)
	if !res {
		rsp := registryResponse{
			Result: res,
			Msg:    "No se ha podido registrar en la base de datos el archivo",
		}
		ctx.JSON(http.StatusOK, rsp)
		return
	}

	var msg string

	if versionado {
		//Es la segunda vez que se almacena el archivo...
		// leer datos de origen
		origen, err := os.Open("ArchivosUsuarios/" + strconv.Itoa(idfolder) + "/" + file.Filename)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		defer origen.Close()
		// crea un nuevo archivo

		destino, err := os.OpenFile("ArchivosUsuarios/"+strconv.Itoa(idfolder)+"/versiones/"+file.Filename+strconv.Itoa(version), os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		// cierra el archivo "destino.txt" al terminar programa
		defer destino.Close()
		// copiar datos
		_, err = io.Copy(destino, origen)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		msg = "Almacenado correctamente, además se ha guardado en el historial de versiones el anterior archivo"

	} else {
		msg = "Almacenado correctamente"
	}

	//Es la primera vez que se almacena el archivo
	err = ctx.SaveUploadedFile(file, "ArchivosUsuarios/"+strconv.Itoa(idfolder)+"/"+file.Filename)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	rsp := registryResponse{
		Result: res,
		Msg:    msg,
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) getNameFiles(ctx *gin.Context) {
	username := ctx.Request.PostFormValue("username")
	esversion := ctx.Request.PostFormValue("esversion")
	idfile := ctx.Request.PostFormValue("idfile")

	//Obtenemos nuestra clave privada
	clavePrivada := leerClavePrivada()

	//Convertimos el cotenido recibido de base64 a bytes[]
	usernameCifradoBytes, err := base64.Encoding.Strict(*base64.StdEncoding).DecodeString(username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Desencriptamos con nuestra clave privada...
	usernameDescifrado, err := RsaDecrypt(usernameCifradoBytes, []byte(clavePrivada))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Convertimos el cotenido recibido de base64 a bytes[]
	esversionCifradoBytes, err := base64.Encoding.Strict(*base64.StdEncoding).DecodeString(esversion)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Desencriptamos con nuestra clave privada...
	esversionDescifrado, err := RsaDecrypt(esversionCifradoBytes, []byte(clavePrivada))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Convertimos el cotenido recibido de base64 a bytes[]
	idfileCifradoBytes, err := base64.Encoding.Strict(*base64.StdEncoding).DecodeString(idfile)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Desencriptamos con nuestra clave privada...
	idfileDescifrado, err := RsaDecrypt(idfileCifradoBytes, []byte(clavePrivada))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var res string

	if string(esversionDescifrado) == "1" {
		res = db.ObtenerArchivosUsuarioVersiones(string(idfileDescifrado))

	} else {
		idfolder := db.ObtenerIdFolder(string(usernameDescifrado))
		res = db.ObtenerArchivosUsuario(strconv.Itoa(idfolder))

	}

	rsp := uploadResponse{
		Result: true,
		Msg:    res,
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) download(ctx *gin.Context) {

	username := ctx.Request.PostFormValue("username")
	idfile := ctx.Request.PostFormValue("idfile")
	esversion := ctx.Request.PostFormValue("esversion")

	//Obtenemos nuestra clave privada
	clavePrivada := leerClavePrivada()

	//Convertimos el cotenido recibido de base64 a bytes[]
	usernameCifradoBytes, err := base64.Encoding.Strict(*base64.StdEncoding).DecodeString(username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Desencriptamos con nuestra clave privada...
	usernameDescifrado, err := RsaDecrypt(usernameCifradoBytes, []byte(clavePrivada))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Convertimos el cotenido recibido de base64 a bytes[]
	idfileCifradoBytes, err := base64.Encoding.Strict(*base64.StdEncoding).DecodeString(idfile)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Desencriptamos con nuestra clave privada...
	idfileDescifrado, err := RsaDecrypt(idfileCifradoBytes, []byte(clavePrivada))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Convertimos el cotenido recibido de base64 a bytes[]
	esversionCifradoBytes, err := base64.Encoding.Strict(*base64.StdEncoding).DecodeString(esversion)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Desencriptamos con nuestra clave privada...
	esversionDescifrado, err := RsaDecrypt(esversionCifradoBytes, []byte(clavePrivada))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	idfolder := db.ObtenerIdFolder(string(usernameDescifrado))

	var path string
	if string(esversionDescifrado) == "1" {
		filename := db.ObtenerFileNameVersion(string(idfileDescifrado))
		path = "ArchivosUsuarios/" + strconv.Itoa(idfolder) + "/versiones/" + filename
	} else {
		filename := db.ObtenerFileName(string(idfileDescifrado))
		path = "ArchivosUsuarios/" + strconv.Itoa(idfolder) + "/" + filename
	}

	ctx.Status(http.StatusOK)
	ctx.File(path)
}

func (server *Server) delete(ctx *gin.Context) {
	username := ctx.Request.PostFormValue("username")
	idfile := ctx.Request.PostFormValue("idfile")
	esversion := ctx.Request.PostFormValue("esversion")

	//Obtenemos nuestra clave privada
	clavePrivada := leerClavePrivada()

	//Convertimos el cotenido recibido de base64 a bytes[]
	usernameCifradoBytes, err := base64.Encoding.Strict(*base64.StdEncoding).DecodeString(username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Desencriptamos con nuestra clave privada...
	usernameDescifrado, err := RsaDecrypt(usernameCifradoBytes, []byte(clavePrivada))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Convertimos el cotenido recibido de base64 a bytes[]
	idfileCifradoBytes, err := base64.Encoding.Strict(*base64.StdEncoding).DecodeString(idfile)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Desencriptamos con nuestra clave privada...
	idfileDescifrado, err := RsaDecrypt(idfileCifradoBytes, []byte(clavePrivada))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Convertimos el cotenido recibido de base64 a bytes[]
	esversionCifradoBytes, err := base64.Encoding.Strict(*base64.StdEncoding).DecodeString(esversion)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Desencriptamos con nuestra clave privada...
	esversionDescifrado, err := RsaDecrypt(esversionCifradoBytes, []byte(clavePrivada))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var path string
	var done bool
	idfolder := db.ObtenerIdFolder(string(usernameDescifrado))

	if string(esversionDescifrado) == "1" {
		filename := db.ObtenerFileNameVersion(string(idfileDescifrado))
		path = "ArchivosUsuarios/" + strconv.Itoa(idfolder) + "/versiones/" + filename
		done = db.EliminarFileNameVersion(string(idfileDescifrado))
	} else {
		filename := db.ObtenerFileName(string(idfileDescifrado))
		path = "ArchivosUsuarios/" + strconv.Itoa(idfolder) + "/" + filename
		done = db.EliminarFileName(string(idfileDescifrado))
	}

	if !done {
		rsp := uploadResponse{
			Result: false,
			Msg:    "No se ha podido eliminar",
		}
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	rsp := uploadResponse{
		Result: true,
		Msg:    "Eliminado correctamente",
	}

	err = os.Remove(path)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) getComment(ctx *gin.Context) {
	idfile := ctx.Request.PostFormValue("idfile")

	//Obtenemos nuestra clave privada
	clavePrivada := leerClavePrivada()

	//Convertimos el cotenido recibido de base64 a bytes[]
	idfileCifradoBytes, err := base64.Encoding.Strict(*base64.StdEncoding).DecodeString(idfile)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Desencriptamos con nuestra clave privada...
	idfileDescifrado, err := RsaDecrypt(idfileCifradoBytes, []byte(clavePrivada))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	comment := db.ObtenerComment(string(idfileDescifrado))

	ctx.JSON(http.StatusOK, comment)
}

func (server *Server) updateComment(ctx *gin.Context) {
	idfile := ctx.Request.PostFormValue("idfile")
	comment := ctx.Request.PostFormValue("comment")

	//Obtenemos nuestra clave privada
	clavePrivada := leerClavePrivada()

	//Convertimos el cotenido recibido de base64 a bytes[]
	idfileCifradoBytes, err := base64.Encoding.Strict(*base64.StdEncoding).DecodeString(idfile)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Desencriptamos con nuestra clave privada...
	idfileDescifrado, err := RsaDecrypt(idfileCifradoBytes, []byte(clavePrivada))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Convertimos el cotenido recibido de base64 a bytes[]
	commentCifradoBytes, err := base64.Encoding.Strict(*base64.StdEncoding).DecodeString(comment)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Desencriptamos con nuestra clave privada...
	commentDescifrado, err := RsaDecrypt(commentCifradoBytes, []byte(clavePrivada))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	done := db.ModificarComment(string(idfileDescifrado), string(commentDescifrado))
	if !done {
		rsp := uploadResponse{
			Result: false,
			Msg:    "No se ha podido modificar",
		}
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	rsp := uploadResponse{
		Result: true,
		Msg:    "Modificado correctamente",
	}
	ctx.JSON(http.StatusOK, rsp)
}
