package api

import (
	"github.com/gin-gonic/gin"
	 
)

type Server struct {
	router *gin.Engine
 
}

func NewServer() (*Server ) {


 
	server := &Server{
 
	}
	router := gin.Default()

	router.POST("/login", server.login)
	router.POST("/registrar", server.registrar)

	authRoutes := router.Group("/").Use(authMiddleware( ))
	authRoutes.GET("/users", server.getUsers)
	authRoutes.POST("/upload", server.uploadFile)
	server.router = router
	return server 
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
