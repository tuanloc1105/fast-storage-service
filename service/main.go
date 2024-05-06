package main

import (
	"context"
	"errors"
	"fast-storage-go-service/config"
	"fast-storage-go-service/constant"
	"fast-storage-go-service/controller"
	"fast-storage-go-service/log"
	"fast-storage-go-service/payload"
	"fast-storage-go-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	var ctx = context.Background()

	log.WithLevel(
		constant.Info,
		ctx,
		">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Application starting <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<",
	)
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

	db, err := config.InitDatabaseConnection()
	if err != nil {
		panic(err)
	}

	controller.AuthenticateController(router, db)
	controller.StorageController(router, db)

	if isKeycloakInfoSet := config.CheckKeycloakInfo(); !isKeycloakInfoSet {
		panic(errors.New("keycloak is required to run this application"))
	}

	if isMountFolderFromEnvSet := config.MountFolderLocationConfig(); !isMountFolderFromEnvSet {
		panic(errors.New("a mount folder is required to run this application, consider to add a mount folder directory path to the environment"))
	}

	applicationPort := "8080"

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
