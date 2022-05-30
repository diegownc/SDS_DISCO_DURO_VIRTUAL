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

type FileResponse struct {
	Result  bool   `json:"result"`
	Path    string `json:"path"`
	Size    int    `json:"size"`
	Comment string `json:"comment"`
	Content string `json:"content"`
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
	esversion := ctx.Request.PostFormValue("esversion")
	idfile := ctx.Request.PostFormValue("idfile")

	//Encriptamos el usuario
	clavePublica := leerClavePublica()

	usernameCifradoRSA, err := RsaEncrypt([]byte(username), []byte(clavePublica))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Antes de enviarlo convierto el contenido a base64
	usernameCifrado := base64.StdEncoding.EncodeToString(usernameCifradoRSA)

	esversionCifradoRSA, err := RsaEncrypt([]byte(esversion), []byte(clavePublica))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Antes de enviarlo convierto el contenido a base64
	esversionCifrado := base64.StdEncoding.EncodeToString(esversionCifradoRSA)

	idfileCifradoRSA, err := RsaEncrypt([]byte(idfile), []byte(clavePublica))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Antes de enviarlo convierto el contenido a base64
	idfileCifrado := base64.StdEncoding.EncodeToString(idfileCifradoRSA)

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("username", usernameCifrado)
	bodyWriter.WriteField("esversion", esversionCifrado)
	bodyWriter.WriteField("idfile", idfileCifrado)
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
	esversion := ctx.Request.PostFormValue("esversion")

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

	esversionCifradoRSA, err := RsaEncrypt([]byte(esversion), []byte(clavePublica))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Antes de enviarlo convierto el contenido a base64
	usernameCifrado := base64.StdEncoding.EncodeToString(usernameCifradoRSA)
	idfileCifrado := base64.StdEncoding.EncodeToString(idfileCifradoRSA)
	esversionCifrado := base64.StdEncoding.EncodeToString(esversionCifradoRSA)

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("username", usernameCifrado)
	bodyWriter.WriteField("idfile", idfileCifrado)
	bodyWriter.WriteField("esversion", esversionCifrado)
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
	file, err := os.Create("Descargas/" + filename) // crea el fichero de destino (servidor)
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

	ctx.JSON(http.StatusOK, "Se ha descargado un archivo llamado "+filename+" y se ha guardado en Descargas")
}

func (server *Server) delete(ctx *gin.Context) {
	tokenUsuario := ctx.Request.PostFormValue("tokenUsuario")
	username := ctx.Request.PostFormValue("username")
	idfile := ctx.Request.PostFormValue("idfile")
	esversion := ctx.Request.PostFormValue("esversion")

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

	esversionCifradoRSA, err := RsaEncrypt([]byte(esversion), []byte(clavePublica))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Antes de enviarlo convierto el contenido a base64
	usernameCifrado := base64.StdEncoding.EncodeToString(usernameCifradoRSA)
	idfileCifrado := base64.StdEncoding.EncodeToString(idfileCifradoRSA)
	esversionCifrado := base64.StdEncoding.EncodeToString(esversionCifradoRSA)

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("username", usernameCifrado)
	bodyWriter.WriteField("idfile", idfileCifrado)
	bodyWriter.WriteField("esversion", esversionCifrado)
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	url := "https://localhost:8081/delete"
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

func (server *Server) getFileProperties(ctx *gin.Context) {
	tokenUsuario := ctx.Request.PostFormValue("tokenUsuario")
	username := ctx.Request.PostFormValue("username")
	idfile := ctx.Request.PostFormValue("idfile")
	var esversion string = "0"

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

	esversionCifradoRSA, err := RsaEncrypt([]byte(esversion), []byte(clavePublica))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	//Antes de enviarlo convierto el contenido a base64
	usernameCifrado := base64.StdEncoding.EncodeToString(usernameCifradoRSA)
	idfileCifrado := base64.StdEncoding.EncodeToString(idfileCifradoRSA)
	esversionCifrado := base64.StdEncoding.EncodeToString(esversionCifradoRSA)

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("username", usernameCifrado)
	bodyWriter.WriteField("idfile", idfileCifrado)
	bodyWriter.WriteField("esversion", esversionCifrado)
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

	//Eliminamos la carpeta y la volvemos a crear,,,
	os.RemoveAll("temp-files/")
	_, err = crearDirectorioSiNoExiste("temp-files")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	filename := "archivo" + strconv.Itoa(time.Now().Nanosecond()) + strconv.Itoa(time.Now().Second())
	file, err := os.Create("temp-files/" + filename) // crea el fichero de destino (servidor)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	defer file.Close() // cierra el fichero al salir de ámbito

	//Guardamos en file el contenido de lo que devuelve el servidor
	io.Copy(file, resp.Body) // copia desde el Body del request al fichero con streaming

	var fileSize int64
	fileStat, err := file.Stat()
	if err == nil {
		fileSize = fileStat.Size()
	} else {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Obtenemos el contenido....
	var text string

	if int(fileSize) < 100000 {
		fileContent, err := ioutil.ReadFile("temp-files/" + filename)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		text = string(fileContent)
	} else {
		text = "El archivo es muy pesado y no se pueve previsualizar"
	}
	//Obtenemos el comentario del archivo...
	bodyBuf2 := &bytes.Buffer{}
	bodyWriter2 := multipart.NewWriter(bodyBuf2)
	bodyWriter2.WriteField("idfile", idfileCifrado)
	contentType2 := bodyWriter2.FormDataContentType()
	bodyWriter2.Close()

	url2 := "https://localhost:8081/getComment"
	req3, err := http.NewRequest("POST", url2, bodyBuf2)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	req3.Header.Set("Content-Type", contentType2)
	req3.Header.Set("Authorization", authorizationString)

	client2 := &http.Client{Transport: tr}
	resp2, err := client2.Do(req3)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	defer resp2.Body.Close()

	body2, _ := ioutil.ReadAll(resp2.Body)

	fmt.Println(string(body2))
	rsp := FileResponse{
		Result:  true,
		Path:    "temp-files/" + filename,
		Size:    int(fileSize),
		Comment: string(body2),
		Content: text,
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) updateComment(ctx *gin.Context) {
	tokenUsuario := ctx.Request.PostFormValue("tokenUsuario")
	comment := ctx.Request.PostFormValue("comment")
	idfile := ctx.Request.PostFormValue("idfile")

	clavePublica := leerClavePublica()

	commentCifradoRSA, err := RsaEncrypt([]byte(comment), []byte(clavePublica))
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
	commentCifrado := base64.StdEncoding.EncodeToString(commentCifradoRSA)
	idfileCifrado := base64.StdEncoding.EncodeToString(idfileCifradoRSA)

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("comment", commentCifrado)
	bodyWriter.WriteField("idfile", idfileCifrado)
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	url := "https://localhost:8081/updateComment"
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
