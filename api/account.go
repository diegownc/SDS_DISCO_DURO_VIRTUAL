package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"io"
	"log"
	"os"

	db "github.com/diegownc/SDS_DISCO_DURO_VIRTUAL/db"
	"github.com/diegownc/SDS_DISCO_DURO_VIRTUAL/token"
	"github.com/gin-gonic/gin"
)

const AccessTokenDuration = time.Duration(time.Minute)

type loginRequest struct {
	User     string `form:"username" json:"username" xml:"username" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

type loginResponse struct {
	Result      bool   `json:"result"`
	AccessToken string `json:"access_token"`
}

func (server *Server) registrar(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fmt.Println("Usuario Cifrado: ", req.User)
	fmt.Println()
	fmt.Println("Contraseña Cifrada: ", req.Password)
	fmt.Println()

	//Obtenemos nuestra clave privada
	clavePrivada := leerClavePrivada()

	//Convertimos el cotenido recibido de base64 a bytes[]
	usernameCifradoBytes, err := base64.Encoding.Strict(*base64.StdEncoding).DecodeString(req.User)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	passwordCifradoBytes, err := base64.Encoding.Strict(*base64.StdEncoding).DecodeString(req.Password)
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

	passwordDescifrado, err := RsaDecrypt(passwordCifradoBytes, []byte(clavePrivada))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fmt.Println("Usuario descifrado: ", string(usernameDescifrado))
	fmt.Println()
	fmt.Println("Password descifrado: ", string(passwordDescifrado))

	res := db.RegistroUsuario(string(usernameDescifrado), string(passwordDescifrado))

	if res {
		ctx.JSON(http.StatusOK, "Te has registrado correctamente")
	} else {
		ctx.JSON(http.StatusOK, "Ha ocurrido un error")
	}

}

func (server *Server) login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fmt.Println("Usuario Cifrado: ", req.User)
	fmt.Println()
	fmt.Println("Contraseña Cifrada: ", req.Password)
	fmt.Println()

	//Obtenemos nuestra clave privada
	clavePrivada := leerClavePrivada()

	//Convertimos el cotenido recibido de base64 a bytes[]
	usernameCifradoBytes, err := base64.Encoding.Strict(*base64.StdEncoding).DecodeString(req.User)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	passwordCifradoBytes, err := base64.Encoding.Strict(*base64.StdEncoding).DecodeString(req.Password)
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

	passwordDescifrado, err := RsaDecrypt(passwordCifradoBytes, []byte(clavePrivada))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fmt.Println("Usuario descifrado: ", string(usernameDescifrado))
	fmt.Println()
	fmt.Println("Password descifrado: ", string(passwordDescifrado))

	res := db.LoginUsuario(string(usernameDescifrado), string(passwordDescifrado))

	tokenMaker, err := token.NewJWTMaker("12345678123456781234567812345678")

	accessToken, err := tokenMaker.CreateToken(
		req.User,
		AccessTokenDuration,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	rsp := loginResponse{
		Result:      res,
		AccessToken: accessToken,
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) getUsers(ctx *gin.Context) {

	fmt.Println("METODO QUE REQUIERE AUTENTIFICACION")

}

func (server *Server) uploadFile(ctx *gin.Context) {
	name := ctx.PostForm("name")
	fmt.Println(name)

	file, header, err := ctx.Request.FormFile("upload")
	if err != nil {
		ctx.String(http.StatusBadRequest, "Bad request")
		return
	}
	filename := header.Filename

	fmt.Println(file, err, filename)

	out, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	ctx.JSON(http.StatusCreated, "Upload succesful")
}
