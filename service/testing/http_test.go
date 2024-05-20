package main

import (
	"context"
	"fast-storage-go-service/constant"
	"fast-storage-go-service/log"
	"fast-storage-go-service/utils"
	"testing"
)

func TestHttpConsume(t *testing.T) {
	var ctx = context.Background()
	consumeApiResult, consumeApiError := utils.ConsumeApi(
		utils.ConsumeApiOption{
			Ctx:         ctx,
			Url:         "https://polliwog-one-rarely.ngrok-free.app/fast_storage/api/v1/auth/login",
			Method:      "POST",
			Header:      make(map[string]string),
			Payload:     `{"request":{"username":"admin","password":"admin"}}`,
			IsVerifySsl: false,
		},
	)

	if consumeApiError != nil {
		log.WithLevel(
			constant.Error,
			ctx,
			consumeApiError.Error(),
		)
	} else {
		log.WithLevel(
			constant.Info,
			ctx,
			consumeApiResult,
		)
	}
}
