package api

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

type errorType struct {
	Result bool   `json:"result"`
	Msg    string `json:"msg"`
}

func NewServer() *Server {
 
	server := &Server{}
	router := gin.Default()
	router.SetTrustedProxies([]string{"localhost"})

	router.POST("/login", server.login)
	router.POST("/registrar", server.registrar)
	router.POST("/upload", server.uploadFile)

	authRoutes := router.Group("/").Use(authMiddleware())
	authRoutes.GET("/users", server.getUsers)
	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) errorType {
	rsp := errorType{
		Result: false,
		Msg:    err.Error(),
	}

	//return gin.H{"error": err.Error()}
	return rsp
}
