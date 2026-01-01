package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nepalsaurav/taskflow/models"
	"github.com/nepalsaurav/taskflow/pkg"
)

func (apiRouter *ApiRouter) MailRouter() {

	router := apiRouter.router

	router.POST("/postfix/config", func(ctx *gin.Context) {
		postfixConfig := pkg.PostfixConfig{}

		if err := ctx.ShouldBindJSON(&postfixConfig); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}
		postfixconfig, err := pkg.SetPostfixConfig(postfixConfig)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := models.SetSetting(apiRouter.db, "postfix_config", postfixconfig); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "postfix configuration updated successfully",
		})

	})

	router.GET("/postfix/queue", func(ctx *gin.Context) {
		entry, err := pkg.GetPostfixQueue()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}
		ctx.JSON(http.StatusOK, entry)

	})

	router.GET("/postfix/log", func(ctx *gin.Context) {
		report, err := pkg.PostfixLogDetail()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"report": report,
		})
	})

	router.GET("/settings", func(ctx *gin.Context) {
		db, err := models.DefaultDBConnect("database/models.db")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		settings, err := models.GetAllSettingsAsMap(db)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, settings)
	})

}
