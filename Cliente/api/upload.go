package api

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

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
	bodyWriter.WriteField("username", username)
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

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("tokenUsuario", tokenUsuario)
	bodyWriter.WriteField("username", username)
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

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("tokenUsuario", tokenUsuario)
	bodyWriter.WriteField("username", username)
	bodyWriter.WriteField("idfile", idfile)
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

	//Obtenemos la respuesta del servidor
	//body, _ := ioutil.ReadAll(resp.Body)

	file, err := os.Create("temp-files/patata") // crea el fichero de destino (servidor)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	defer file.Close() // cierra el fichero al salir de ámbito

	io.Copy(file, resp.Body) // copia desde el Body del request al fichero con streaming

	fmt.Println("enviare el archivo...")
	ctx.Header("Content-Type", "application/octet-stream")
	//Force browser download
	ctx.Header("Content-Disposition", "attachment; filename="+file.Name())
	//Browser download or preview
	ctx.Header("Content-Disposition", "inline;filename="+file.Name())
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Cache-Control", "no-cache")
	ctx.File("temp-files/patata")
}
