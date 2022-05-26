package api

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	tokenMaker tokenMaker 
}

func NewServer() (*Server, error) {
	tokenMaker, err := token.NewJWTMaker("12345678123456781234567812345678")

	if err != nil {
		return nil, fmt.Error("Cannot create token maker: %w", err)
	}
	server := &Server{
		tokenMaker: tokenMaker,
	}
	router := gin.Default()

	router.POST("/login", server.login)
	router.POST("/registrar", server.registrar)

	server.router = router
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
