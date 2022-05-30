package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	db "github.com/diegownc/SDS_DISCO_DURO_VIRTUAL/db"
	"github.com/diegownc/SDS_DISCO_DURO_VIRTUAL/token"
	"github.com/gin-gonic/gin"
)

const AccessTokenDuration = time.Duration(time.Hour)

type loginRequest struct {
	User     string `form:"username" json:"username" xml:"username" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

//Definimos los tipos de respuesta (todos tienen un Result para controlarlo en el js)
type loginResponseSuccess struct {
	Result      bool   `json:"result"`
	AccessToken string `json:"access_token"`
}

type loginResponseFailure struct {
	Result bool   `json:"result"`
	Msg    string `json:"msg"`
}

type registryResponse struct {
	Result bool   `json:"result"`
	Msg    string `json:"msg"`
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
		//Creamos la carpeta del cliente
		_, err = crearDirectorioSiNoExiste("ArchivosUsuarios/" + strconv.Itoa(db.ObtenerIdFolder(string(usernameDescifrado))))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		//Creamos la carpeta versiones...
		res, err = crearDirectorioSiNoExiste("ArchivosUsuarios/" + strconv.Itoa(db.ObtenerIdFolder(string(usernameDescifrado))) + "/versiones")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		rsp := registryResponse{
			Result: res,
			Msg:    "Registrado correctamente.",
		}
		ctx.JSON(http.StatusOK, rsp)

	} else {
		rsp := registryResponse{
			Result: res,
			Msg:    "Este usuario ya esta registrado.",
		}
		ctx.JSON(http.StatusOK, rsp)
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

	if res {
		rsp := loginResponseSuccess{
			Result:      res,
			AccessToken: accessToken,
		}
		ctx.JSON(http.StatusOK, rsp)

	} else {
		rsp := loginResponseFailure{
			Result: res,
			Msg:    "Credenciales incorrectas",
		}
		ctx.JSON(http.StatusOK, rsp)
	}

}

func (server *Server) getUsers(ctx *gin.Context) {

	fmt.Println("METODO QUE REQUIERE AUTENTIFICACION")

}

func crearDirectorioSiNoExiste(directorio string) (bool, error) {
	var todoOk bool = false
	var err error = nil
	if _, err = os.Stat(directorio); os.IsNotExist(err) {
		err = os.Mkdir(directorio, 0755)
		if err != nil {
			panic(err)
		}
		todoOk = true
	}

	return todoOk, err
}
