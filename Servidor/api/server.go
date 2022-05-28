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

	authRoutes := router.Group("/").Use(authMiddleware())
	authRoutes.GET("/users", server.getUsers)
	authRoutes.POST("/nameFiles", server.getNameFiles)
	authRoutes.POST("/upload", server.uploadFile)
	authRoutes.POST("/download", server.download)
	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	//return server.router.Run(address)
	return server.router.RunTLS(address, "certificados/localhost.crt", "certificados/localhost.key")
}

func errorResponse(err error) errorType {
	rsp := errorType{
		Result: false,
		Msg:    err.Error(),
	}

	//return gin.H{"error": err.Error()}
	return rsp
}
