package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	db "github.com/diegownc/SDS_DISCO_DURO_VIRTUAL/db"
	"github.com/diegownc/SDS_DISCO_DURO_VIRTUAL/token"
	"github.com/gin-gonic/gin"
)

type uploadResponse struct {
	Result bool   `json:"result"`
	Msg    string `json:"msg"`
}

func (server *Server) uploadFile(ctx *gin.Context) {

	tokenUsuario := ctx.Request.PostFormValue("tokenUsuario")
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

	tokenMaker, err := token.NewJWTMaker("12345678123456781234567812345678")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = tokenMaker.VerifyToken(tokenUsuario)
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
	err = ctx.SaveUploadedFile(file, "ArchivosUsuarios/"+strconv.Itoa(idfolder)+"/"+file.Filename)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Almacenamos en base de datos el nombre del archivo..
	res := db.RegistrarArchivo(file.Filename, "", idfolder)
	if !res {
		rsp := registryResponse{
			Result: res,
			Msg:    "No se ha podido registrar en la base de datos el archivo",
		}
		ctx.JSON(http.StatusOK, rsp)
		return
	}

	rsp := registryResponse{
		Result: res,
		Msg:    "Almacenado correctamente.",
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) getNameFiles(ctx *gin.Context) {

	tokenUsuario := ctx.Request.PostFormValue("tokenUsuario")
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

	tokenMaker, err := token.NewJWTMaker("12345678123456781234567812345678")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = tokenMaker.VerifyToken(tokenUsuario)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	idfolder := db.ObtenerIdFolder(string(usernameDescifrado))
	res := db.ObtenerArchivosUsuario(strconv.Itoa(idfolder))

	rsp := uploadResponse{
		Result: true,
		Msg:    res,
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) download(ctx *gin.Context) {

	tokenUsuario := ctx.Request.PostFormValue("tokenUsuario")
	username := ctx.Request.PostFormValue("username")
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

	tokenMaker, err := token.NewJWTMaker("12345678123456781234567812345678")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = tokenMaker.VerifyToken(tokenUsuario)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	idfolder := db.ObtenerIdFolder(string(usernameDescifrado))
	filename := db.ObtenerFileName(string(idfileDescifrado))

	path := "ArchivosUsuarios/" + strconv.Itoa(idfolder) + "/" + filename
	fmt.Println("El path del archivo es.. " + path)

	ctx.Status(200)
	ctx.File(path)
}
