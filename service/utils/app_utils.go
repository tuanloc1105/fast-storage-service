package utils

import (
	"bytes"
	"context"
	"fast-storage-go-service/constant"
	"fast-storage-go-service/log"
	"fmt"
	"math/big"
	"math/rand"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CheckAndSetTraceId(c *gin.Context) {
	if traceId, _ := c.Get(string(constant.TraceIdLogKey)); traceId == nil || traceId == "" {
		c.Set(string(constant.TraceIdLogKey), strings.Replace(uuid.New().String(), "-", constant.EmptyString, -1))
	}
}

func GetTraceId(c *gin.Context) string {
	if traceId, _ := c.Get(string(constant.TraceIdLogKey)); traceId == nil || traceId == "" {
		return ""
	} else {
		return traceId.(string)
	}
}

func RoundHalfUpBigFloat(input *big.Float) {
	delta := constant.DeltaPositive

	if input.Sign() < 0 {
		delta = constant.DeltaNegative
	}
	input.Add(input, new(big.Float).SetFloat64(delta))
}

func GetPointerOfAnyValue[T any](a T) *T {
	return &a
}

func Shellout(ctx context.Context, command string, isLog ...bool) (string, string, error) {
	currentCommandRunningId := strings.Replace(uuid.New().String(), "-", constant.EmptyString, -1)
	if len(isLog) < 1 || (len(isLog) == 1 && isLog[0]) {
		log.WithLevel(
			constant.Info,
			ctx,
			"[%s] == Start to executing command: %s",
			currentCommandRunningId,
			HideSensitiveInformationOfCurlCommand(command),
		)
	}
	var cmd *exec.Cmd
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("bash", "-c", command)
	case "windows":
		cmd = exec.Command("cmd", "/c", command)
	default:
		log.WithLevel(
			constant.Error,
			ctx,
			"[%s] == %s not implemented",
			currentCommandRunningId,
			runtime.GOOS,
		)
		return "", "", fmt.Errorf("%s not implemented", runtime.GOOS)
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exitCode := cmd.ProcessState.ExitCode()
	stdoutString := strings.TrimPrefix(strings.TrimSuffix(stdout.String(), "\n"), "\n")
	stderrString := strings.TrimPrefix(strings.TrimSuffix(stderr.String(), "\n"), "\n")
	if len(isLog) < 1 || (len(isLog) == 2 && isLog[1]) {
		var finalStdoutString string
		if IsStringAJson(stdoutString) {
			finalStdoutString = HideSensitiveJsonField(stdoutString)
		} else {
			finalStdoutString = stdoutString
		}
		log.WithLevel(
			constant.Info,
			ctx,
			`[%s] == command result:
    - status code: %d
    - stdout: 
%s
    - stderr: 
%s`,
			currentCommandRunningId,
			exitCode,
			finalStdoutString,
			stderrString,
		)
	}
	return stdoutString, stderrString, err
}

func ShelloutAtSpecificDirectory(ctx context.Context, command, directory string, isLog ...bool) (string, string, error) {
	currentCommandRunningId := strings.Replace(uuid.New().String(), "-", constant.EmptyString, -1)
	if len(isLog) < 1 || (len(isLog) == 1 && isLog[0]) {
		log.WithLevel(
			constant.Info,
			ctx,
			"[%s] == Start to executing command: %s",
			currentCommandRunningId,
			HideSensitiveInformationOfCurlCommand(command),
		)
	}
	if directory == "" {
		return "", "", fmt.Errorf("`directory` can not be empty")
	}
	var cmd *exec.Cmd
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("bash", "-c", command)
	case "windows":
		cmd = exec.Command("cmd", "/c", command)
	default:
		log.WithLevel(
			constant.Error,
			ctx,
			"[%s] == %s not implemented",
			currentCommandRunningId,
			runtime.GOOS,
		)
		return "", "", fmt.Errorf("%s not implemented", runtime.GOOS)
	}
	cmd.Dir = directory
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exitCode := cmd.ProcessState.ExitCode()
	stdoutString := strings.TrimPrefix(strings.TrimSuffix(stdout.String(), "\n"), "\n")
	stderrString := strings.TrimPrefix(strings.TrimSuffix(stderr.String(), "\n"), "\n")
	if len(isLog) < 1 || (len(isLog) == 2 && isLog[1]) {
		var finalStdoutString string
		if IsStringAJson(stdoutString) {
			finalStdoutString = HideSensitiveJsonField(stdoutString)
		} else {
			finalStdoutString = stdoutString
		}
		log.WithLevel(
			constant.Info,
			ctx,
			`[%s] == command result:
    - status code: %d
    - stdout: 
%s
    - stderr: 
%s`,
			currentCommandRunningId,
			exitCode,
			finalStdoutString,
			stderrString,
		)
	}
	return stdoutString, stderrString, err
}

func GenerateRandomString(length int) string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = constant.Charset[seededRand.Intn(len(constant.Charset))]
	}

	return string(randomString)
}
