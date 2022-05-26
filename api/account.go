package api

import (
	"net/http"
 	"fmt"
	"time"
	 "github.com/diegownc/SDS_DISCO_DURO_VIRTUAL/token"
	db "github.com/diegownc/SDS_DISCO_DURO_VIRTUAL/db"
	"github.com/gin-gonic/gin"
)

const AccessTokenDuration = time.Duration(time.Minute)  
   
type loginRequest struct {
	User     string `form:"username" json:"username" xml:"username" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

type loginResponse struct {
	Result bool `json:"result"`
	AccessToken string `json:"access_token"`
}

func (server *Server) registrar(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res := db.RegistroUsuario(req.User, req.Password)

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

	res := db.LoginUsuario(req.User, req.Password)

	tokenMaker, err := token.NewJWTMaker("12345678123456781234567812345678")

	accessToken, err := tokenMaker.CreateToken(
		req.User,
		AccessTokenDuration,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
 	} 

	rsp := loginResponse{ 
		Result : res,
		AccessToken: accessToken,
	} 
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) getUsers(ctx *gin.Context) {
 
	fmt.Println("METODO QUE REQUIERE AUTENTIFICACION")

}
