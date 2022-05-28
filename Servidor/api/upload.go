package api

import (
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

	idfolder := db.ObtenerIdFolder(username)
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

	idfolder := db.ObtenerIdFolder(username)
	res := db.ObtenerArchivosUsuario(strconv.Itoa(idfolder))

	rsp := registryResponse{
		Result: true,
		Msg:    res,
	}
	ctx.JSON(http.StatusOK, rsp)
}
