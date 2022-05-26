package api

import (
	"github.com/gin-gonic/gin"
	"github.com/diegownc/SDS_DISCO_DURO_VIRTUAL/token"
)

const (
	authorizationHeaderKey = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware( ) gin.HandlerFunc{
	return func(ctx *gin.Context){
		authorizationHeader _= ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0{
			err := errors.New("Authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return 
		}

		fields .= strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("Invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return 
		}

		authorizationType := strings.ToLower(fields[0])
		if(authorizationType != authorizationTypeBearer){
			err := errors.New("Unsupported authorization type")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return 
		}

		access_token := fields[1]
		payload, err := tokenMaker.VerifyToken(access_token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}