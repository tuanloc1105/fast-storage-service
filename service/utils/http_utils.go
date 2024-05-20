package utils

import (
	"context"
	"crypto/tls"
	"errors"
	"fast-storage-go-service/constant"
	"fast-storage-go-service/log"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
)

type ConsumeApiOption struct {
	Ctx         context.Context
	Url         string
	Method      string
	Header      map[string]string
	Payload     string
	IsVerifySsl bool
}

func ConsumeApi(option ConsumeApiOption) (string, error) {

	if !slices.Contains(constant.ValidMethod, option.Method) {
		return "", errors.New("invalid method")
	}
	var client *http.Client
	if option.IsVerifySsl {
		client = &http.Client{}
	} else {
		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client = &http.Client{Transport: customTransport}
	}
	req, err := http.NewRequest(option.Method, option.Url, strings.NewReader(option.Payload))
	if err != nil {
		log.WithLevel(
			constant.Error,
			option.Ctx,
			"ConsumeApi - http.NewRequest - error: "+err.Error(),
		)
		return "", err
	}

	log.WithLevel(
		constant.Info,
		option.Ctx,
		curlBuilder(option.Url, option.Payload, option.Header),
	)

	for k, v := range option.Header {
		req.Header.Add(k, v)
	}
	res, err := client.Do(req)
	if err != nil {
		log.WithLevel(
			constant.Error,
			option.Ctx,
			"ConsumeApi - client.Do - error: "+err.Error(),
		)
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("%v", err.Error())
		}
	}(res.Body)

	resHeader := map[string][]string(res.Header)

	headerString := ""

	for k, v := range resHeader {
		if IsSensitiveField(k) {
			headerString += fmt.Sprintf("\n\t\t- %s: %s", k, "***")
		} else {
			headerString += fmt.Sprintf("\n\t\t- %s: %s", k, strings.Join(v, ", "))
		}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.WithLevel(
			constant.Error,
			option.Ctx,
			"ConsumeApi - io.ReadAll - error: "+err.Error(),
		)
		return "", err
	}
	result := string(body)
	log.WithLevel(
		constant.Info,
		option.Ctx,
		"\t- status: %s\n\t- header: %s\n\t- payload: %s",
		res.Status,
		headerString,
		result,
	)

	return result, nil

}

func curlBuilder(url string, payload string, header map[string]string) string {
	curlCommand := "curl "
	curlCommand += "'" + url + "' "
	for k, v := range header {
		curlCommand += "-H '" + k + ": " + v + "' "
	}
	if payload != "" {
		curlCommand += "-X POST -d '" + payload + "'"
	} else {
		curlCommand += "-X GET"
	}
	return curlCommand
}
