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
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	var ctx = context.Background()

	log.WithLevel(
		constant.Info,
		ctx,
		">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Application starting <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<",
	)
	// gin.SetMode(gin.DebugMode)
	router := gin.New()
	err := router.SetTrustedProxies([]string{"192.168.1.0/24", "127.0.0.1"})
	if err != nil {
		log.WithLevel(constant.Error, ctx, "Could not set trusted proxies: %s\n", err.Error())
	}

	router.Use(gin.Recovery())

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
	controller.TotpController(router, db)

	if isKeycloakInfoSet := config.CheckKeycloakInfo(); !isKeycloakInfoSet {
		panic(errors.New("keycloak is required to run this application"))
	}

	if isMountFolderFromEnvSet := config.MountFolderLocationConfig(); !isMountFolderFromEnvSet {
		panic(errors.New("a mount folder is required to run this application, consider to add a mount folder directory path to the environment"))
	}

	applicationPort := "8080"

	// Serve static files
	router.Static("/static", "./static")

	// Serve HTML template
	router.LoadHTMLGlob("templates/*")

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

	go func() {
		ginStartUpError := router.Run(":" + applicationPort)
		if ginStartUpError != nil {
			log.WithLevel(constant.Error, ctx, "Error when running server: %v", ginStartUpError)
			os.Exit(1)
		}
	}()

	log.WithLevel(
		constant.Info,
		ctx,
		"Application started on port: "+applicationPort,
	)
	var applicationOsSignal os.Signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	applicationOsSignal = <-quit

	log.WithLevel(constant.Info, ctx, fmt.Sprintln("Shutting down the server with signal", applicationOsSignal.String()))

}
