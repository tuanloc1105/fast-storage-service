package main

import (
	"context"
	"fast-storage-go-service/constant"
	"fast-storage-go-service/log"
	"fast-storage-go-service/payload"
	"fast-storage-go-service/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	var ctx = context.Background()
	router := gin.Default()

	router.NoRoute(
		func(context *gin.Context) {
			context.JSON(
				http.StatusNotFound, &payload.Response{
					Trace:        utils.GetTraceId(context),
					ErrorCode:    constant.PageNotFound.ErrorCode,
					ErrorMessage: constant.PageNotFound.ErrorMessage,
				},
			)
		},
	)

	router.NoMethod(
		func(context *gin.Context) {
			context.JSON(
				http.StatusNotFound, &payload.Response{
					Trace:        utils.GetTraceId(context),
					ErrorCode:    constant.MethodNotAllowed.ErrorCode,
					ErrorMessage: constant.MethodNotAllowed.ErrorMessage,
				},
			)
		},
	)

	applicationPort := os.Getenv("APPLICATION_PORT")
	if applicationPort == "" {
		applicationPort = "8080"
	}

	router.GET(constant.BaseApiPath+"/", func(ctx *gin.Context) {
		ctx.Data(
			http.StatusOK,
			constant.ContentTypeHTML,
			[]byte(`
				<h1>Fast Storage Service</h1>
			`),
		)
	})

	log.WithLevel(
		constant.Info,
		ctx,
		"Current directory is: "+utils.GetCurrentDirectory(),
	)

	log.WithLevel(
		constant.Info,
		ctx,
		"Application starting with port: "+applicationPort,
	)

	gitStartUpError := router.Run(":" + applicationPort)
	if gitStartUpError != nil {
		panic(gitStartUpError)
	}

}
