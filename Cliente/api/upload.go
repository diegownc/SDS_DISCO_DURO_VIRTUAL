package api

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FileItemRequest struct {
	Key      string
	FileName string
	Content  []byte
}

func (server *Server) upload(ctx *gin.Context) {
	file, fileheader, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tokenUsuario := ctx.Request.PostFormValue("tokenUsuario")
	username := ctx.Request.PostFormValue("username")
	//Encriptamos el usuario
	clavePublica := leerClavePublica()

	usernameCifradoRSA, err := RsaEncrypt([]byte(username), []byte(clavePublica))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Antes de enviarlo convierto el contenido a base64
	usernameCifrado := base64.StdEncoding.EncodeToString(usernameCifradoRSA)

	//Preparo el request para enviarselo al servidor
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var fileitem FileItemRequest
	fileitem.Key = "file"
	fileitem.FileName = fileheader.Filename
	fileitem.Content = buf.Bytes()

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("tokenUsuario", tokenUsuario)
	bodyWriter.WriteField("username", usernameCifrado)
	fileWriter, err := bodyWriter.CreateFormFile(fileitem.Key, fileitem.FileName)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fileWriter.Write(fileitem.Content)
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	//Aqui estoy haciendo la petición al servidor... y estoy intentando averiguar como enviar el "file"
	url := "https://localhost:8081/upload"

	req2, err := http.NewRequest("POST", url, bodyBuf)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	req2.Header.Set("Content-Type", contentType)
	authorizationString := "Bearer " + tokenUsuario
	req2.Header.Set("Authorization", authorizationString)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req2)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	ctx.JSON(http.StatusOK, string(body))
}

func (server *Server) getNameFiles(ctx *gin.Context) {
	tokenUsuario := ctx.Request.PostFormValue("tokenUsuario")
	username := ctx.Request.PostFormValue("username")
	//Encriptamos el usuario
	clavePublica := leerClavePublica()

	usernameCifradoRSA, err := RsaEncrypt([]byte(username), []byte(clavePublica))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Antes de enviarlo convierto el contenido a base64
	usernameCifrado := base64.StdEncoding.EncodeToString(usernameCifradoRSA)

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("tokenUsuario", tokenUsuario)
	bodyWriter.WriteField("username", usernameCifrado)
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	url := "https://localhost:8081/nameFiles"
	req2, err := http.NewRequest("POST", url, bodyBuf)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	req2.Header.Set("Content-Type", contentType)
	authorizationString := "Bearer " + tokenUsuario
	req2.Header.Set("Authorization", authorizationString)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req2)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	defer resp.Body.Close()

	//Obtenemos la respuesta del servidor
	body, _ := ioutil.ReadAll(resp.Body)

	ctx.JSON(http.StatusOK, string(body))
}

func (server *Server) download(ctx *gin.Context) {
	tokenUsuario := ctx.Request.PostFormValue("tokenUsuario")
	username := ctx.Request.PostFormValue("username")
	idfile := ctx.Request.PostFormValue("idfile")

	clavePublica := leerClavePublica()

	usernameCifradoRSA, err := RsaEncrypt([]byte(username), []byte(clavePublica))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	idfileCifradoRSA, err := RsaEncrypt([]byte(idfile), []byte(clavePublica))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Antes de enviarlo convierto el contenido a base64
	usernameCifrado := base64.StdEncoding.EncodeToString(usernameCifradoRSA)
	idfileCifrado := base64.StdEncoding.EncodeToString(idfileCifradoRSA)

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("tokenUsuario", tokenUsuario)
	bodyWriter.WriteField("username", usernameCifrado)
	bodyWriter.WriteField("idfile", idfileCifrado)
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	url := "https://localhost:8081/download"
	req2, err := http.NewRequest("POST", url, bodyBuf)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	req2.Header.Set("Content-Type", contentType)
	authorizationString := "Bearer " + tokenUsuario
	req2.Header.Set("Authorization", authorizationString)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req2)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	defer resp.Body.Close()

	filename := "archivo" + strconv.Itoa(time.Now().Nanosecond()) + strconv.Itoa(time.Now().Second())
	file, err := os.Create("temp-files/" + filename) // crea el fichero de destino (servidor)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	defer file.Close() // cierra el fichero al salir de ámbito

	//Guardamos en file el contenido de lo que devuelve el servidor
	io.Copy(file, resp.Body) // copia desde el Body del request al fichero con streaming

	//ctx.Header("Content-Type", "application/octet-stream")
	//Force browser download
	//ctx.Header("Content-Disposition", "attachment; filename="+file.Name())
	//Browser download or preview
	//ctx.Header("Content-Disposition", "inline;filename="+file.Name())
	//ctx.Header("Content-Transfer-Encoding", "binary")
	//ctx.Header("Cache-Control", "no-cache")
	//ctx.File("temp-files/patata")

	ctx.JSON(http.StatusOK, "Se ha descargado un archivo llamado "+filename)
}
