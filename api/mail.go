package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nepalsaurav/taskflow/pkg"
)

func MailRouter(router *gin.Engine) {
	router.GET("/index_mail", func(ctx *gin.Context) {
		maildir := &pkg.Maildir{}

		resp, err := maildir.IndexMail()

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
		}
		ctx.JSON(http.StatusOK, gin.H{
			"msgType": "success",
			"resp":    resp,
		})
	})

	router.POST("/postfix/config", func(ctx *gin.Context) {
		postfixConfig := pkg.PostfixConfig{}

		if err := ctx.ShouldBindJSON(&postfixConfig); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}
		if err := pkg.SetPostfixConfig(postfixConfig); err != nil {
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

}
