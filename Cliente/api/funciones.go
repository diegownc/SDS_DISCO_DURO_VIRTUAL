package api

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	User     string `form:"user" json:"user" xml:"user" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

func (server *Server) registrar(ctx *gin.Context) {
	var req loginRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	//Encriptamos el usuario y la contraseña
	clavePublica := leerClavePublica()

	usernameCifradoRSA, err := RsaEncrypt([]byte(req.User), []byte(clavePublica))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	passwordCifradoRSA, err := RsaEncrypt([]byte(req.Password), []byte(clavePublica))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Antes de enviarlo convierto el contenido a base64
	usernameCifrado := base64.StdEncoding.EncodeToString(usernameCifradoRSA)
	passwordCifrado := base64.StdEncoding.EncodeToString(passwordCifradoRSA)

	url := "http://localhost:8081/registrar"
	var jsonStr = []byte(`{"username": "` + usernameCifrado + `", "password" : "` + passwordCifrado + `"}`)
	req2, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	req2.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req2)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	ctx.JSON(http.StatusOK, gin.H{"Msg": string(body)})

}

func (server *Server) login(ctx *gin.Context) {
	var req loginRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Encriptamos el usuario y la contraseña
	clavePublica := leerClavePublica()

	usernameCifradoRSA, err := RsaEncrypt([]byte(req.User), []byte(clavePublica))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	passwordCifradoRSA, err := RsaEncrypt([]byte(req.Password), []byte(clavePublica))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	usernameCifrado := base64.StdEncoding.EncodeToString(usernameCifradoRSA)
	passwordCifrado := base64.StdEncoding.EncodeToString(passwordCifradoRSA)

	//Hago una petición al server.go
	url := "http://localhost:8081/login"
	var jsonStr = []byte(`{"username": "` + usernameCifrado + `", "password" : "` + passwordCifrado + `"}`)
	req2, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	req2.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req2)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	ctx.JSON(http.StatusOK, gin.H{"Msg": string(body)})

}
