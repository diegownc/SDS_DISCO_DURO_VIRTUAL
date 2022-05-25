package api

import (
	"fmt"
	"net/http"

	db "github.com/diegownc/SDS_DISCO_DURO_VIRTUAL/db"
	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	username string `json:"username" binding:"required "`
	password string `json:"password" binding:"required"`
}

func (server *Server) login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.LoginParams{
		Username: req.username,
		Password: req.password,
	}

	res := db.LoginUsuario(arg.Username, arg.Password)

	/*
		if !res{
			ctx.JSON(http.StatusInternalServerError, errorResponse((err)))
			return
		}

		ctx.JSON(http.StatusOK, account)
	*/
	fmt.Println("El resultado es: ", res)
}
