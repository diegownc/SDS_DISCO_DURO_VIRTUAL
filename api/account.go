package api

import (
	"net/http"

	db "github.com/diegownc/SDS_DISCO_DURO_VIRTUAL/db"
	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	User     string `form:"username" json:"username" xml:"username" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
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

	if res {
		ctx.JSON(http.StatusOK, "Credenciales correctas")
	} else {
		ctx.JSON(http.StatusOK, "Credenciales incorrectas")
	}
}
