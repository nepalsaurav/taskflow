package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nepalsaurav/taskflow/models"
	"gorm.io/gorm"
)

type ApiRouter struct {
	router *gin.Engine
	db     *gorm.DB
}

func (apiRouter *ApiRouter) initRoute() {
	apiRouter.router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	apiRouter.MailRouter()

}

func Serve(address string) {

	db, err := models.DefaultDBConnect("database/models.db")

	if err != nil {
		panic(err)
	}

	apiRouter := &ApiRouter{
		router: gin.Default(),
		db:     db,
	}

	apiRouter.initRoute()

	apiRouter.router.Run(address)
}
