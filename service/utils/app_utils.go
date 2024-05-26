package utils

import (
	"bytes"
	"context"
	"fast-storage-go-service/constant"
	"fast-storage-go-service/log"
	"fmt"
	"math/big"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CheckAndSetTraceId(c *gin.Context) {
	if traceId, _ := c.Get(string(constant.TraceIdLogKey)); traceId == nil || traceId == "" {
		c.Set(string(constant.TraceIdLogKey), uuid.New().String())
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
	if len(isLog) < 1 || (len(isLog) == 1 && isLog[0]) {
		log.WithLevel(
			constant.Info,
			ctx,
			"Start to executing command: %s",
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
			"%s not implemented",
			runtime.GOOS,
		)
		return "", "", fmt.Errorf("%s not implemented", runtime.GOOS)
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exitCode := cmd.ProcessState.ExitCode()
	stdoutString := strings.TrimSuffix(stdout.String(), "\n")
	stderrString := stderr.String()
	if len(isLog) < 1 || (len(isLog) == 2 && isLog[1]) {
		log.WithLevel(
			constant.Info,
			ctx,
			"--- command exit status ---\n%d",
			exitCode,
		)
		if IsStringAJson(stdoutString) {
			log.WithLevel(
				constant.Info,
				ctx,
				"--- stdout ---\n%s",
				HideSensitiveJsonField(stdoutString),
			)
		} else {
			log.WithLevel(
				constant.Info,
				ctx,
				"--- stdout ---\n%s",
				stdoutString,
			)
		}
		log.WithLevel(
			constant.Info,
			ctx,
			"--- stderr ---\n%s",
			stderrString,
		)
	}
	return stdoutString, stderrString, err
}

func ShelloutAtSpecificDirectory(ctx context.Context, command, directory string, isLog ...bool) (string, string, error) {
	if len(isLog) < 1 || (len(isLog) == 1 && isLog[0]) {
		log.WithLevel(
			constant.Info,
			ctx,
			"Start to executing command: %s",
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
			"%s not implemented",
			runtime.GOOS,
		)
		return "", "", fmt.Errorf("%s not implemented", runtime.GOOS)
	}
	cmd.Dir = directory
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exitCode := cmd.ProcessState.ExitCode()
	stdoutString := strings.TrimSuffix(stdout.String(), "\n")
	stderrString := stderr.String()
	if len(isLog) < 1 || (len(isLog) == 2 && isLog[1]) {
		log.WithLevel(
			constant.Info,
			ctx,
			"--- command exit status ---\n%d",
			exitCode,
		)
		if IsStringAJson(stdoutString) {
			log.WithLevel(
				constant.Info,
				ctx,
				"--- stdout ---\n%s",
				HideSensitiveJsonField(stdoutString),
			)
		} else {
			log.WithLevel(
				constant.Info,
				ctx,
				"--- stdout ---\n%s",
				stdoutString,
			)
		}
		log.WithLevel(
			constant.Info,
			ctx,
			"--- stderr ---\n%s",
			stderrString,
		)
	}
	return stdoutString, stderrString, err
}
