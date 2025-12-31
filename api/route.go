package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiRouter struct {
	router *gin.Engine
}

func (apiRouter *ApiRouter) initRoute() {
	apiRouter.router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	MailRouter(apiRouter.router)

}

func Serve(address string) {

	apiRouter := &ApiRouter{
		router: gin.Default(),
	}

	apiRouter.initRoute()

	apiRouter.router.Run(address)
}
