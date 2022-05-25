package api

import (
	"bytes"
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

	fmt.Println(req.User, req.Password)
	url := "http://localhost:8081/registrar"
	var jsonStr = []byte(`{"username": "` + req.User + `", "password" : "` + req.Password + `"}`)
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

	fmt.Println(req.User, req.Password)
	url := "http://localhost:8081/login"
	var jsonStr = []byte(`{"username": "` + req.User + `", "password" : "` + req.Password + `"}`)
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
