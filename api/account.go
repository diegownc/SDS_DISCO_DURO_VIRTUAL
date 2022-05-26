package api

import (
	"net/http"

	db "github.com/diegownc/SDS_DISCO_DURO_VIRTUAL/db"
	"github.com/gin-gonic/gin"
)

const{
	AccessTokenDuration time.Duration("15m") 
}

type loginRequest struct {
	User     string `form:"username" json:"username" xml:"username" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

type loginResponse struct {
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

	accessToken, err := server.tokenMaker.CreateToken(
		req.User,
		AccessTokenDuration,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
 	} 

	rsp := loginResponse{
		AccessToken: accessToken
	} 
	ctx.JSON(http.StatusOK, rsp)
}
