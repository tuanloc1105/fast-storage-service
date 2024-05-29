package log

import (
	"context"
	"fast-storage-go-service/constant"
	"fast-storage-go-service/utils/splunk/v2"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

func WithLevel(level constant.LogLevelType, ctx context.Context, content string, parameter ...any) {

	// ensure that ctx is never nil
	if ctx == nil {
		ctx = context.Background()
		ctx = context.WithValue(ctx, constant.UsernameLogKey, "nil ctx input")
		ctx = context.WithValue(ctx, constant.TraceIdLogKey, "nil ctx input")
	}

	timeZoneLocation, timeLoadLocationErr := time.LoadLocation("Asia/Ho_Chi_Minh")
	if timeLoadLocationErr != nil {
		return
	}
	currentTimestamp := time.Now().In(timeZoneLocation)
	usernameFromContext := ctx.Value(constant.UsernameLogKey)
	traceIdFromContext := ctx.Value(constant.TraceIdLogKey)
	username := constant.EmptyString
	traceId := constant.EmptyString
	if usernameFromContext != nil {
		username = usernameFromContext.(string)
	}
	if traceIdFromContext != nil {
		traceId = traceIdFromContext.(string)
	}
	// fmt.Println(strings.Compare(string(level), string(constant.LogLevelType("INFO"))))
	podName, _ := os.LookupEnv("POD_NAME")
	var message string
	if len(parameter) < 1 {
		message = fmt.Sprintf(
			constant.LogPattern,
			podName,
			traceId,
			username,
			content,
		)
	} else {
		message = fmt.Sprintf(
			constant.LogPattern,
			podName,
			traceId,
			username,
			fmt.Sprintf(content, parameter...),
		)
	}
	switch level {
	case constant.Info:
		log.Info(message)
	case constant.Warn:
		log.Warn(message)
	case constant.Error:
		log.Error(message)
	case constant.Debug:
		log.Debug(message)
	default:
		log.Info(message)
	}

	sendLogToSplunkChan := make(chan struct{})

	go func(done chan struct{}) {
		defer close(done)
		host, token, source, sourcetype, index, splunkInfoIsFullSetInEnv := GetSplunkInformationFromEnvironment()

		if splunkInfoIsFullSetInEnv {
			splunkClient := splunk.NewClient(
				nil,
				host,
				token,
				source,
				sourcetype,
				index,
			)
			err := splunkClient.Log(
				message,
			)
			if err != nil {
				log.Error(err)
			}
		}
	}(sendLogToSplunkChan)

	appendLogToFileError := AppendLogToFile(
		fmt.Sprintf(
			"%s: %s - %s\n",
			currentTimestamp.Format(constant.YyyyMmDdHhMmSsFormat),
			string(level),
			message,
		),
	)
	if appendLogToFileError != nil {
		log.Error(fmt.Sprintf(
			constant.LogPattern,
			podName,
			traceId,
			username,
			fmt.Sprintf("An error has been occurred when appending log to file: %s", appendLogToFileError.Error()),
		))
	}

	// select {
	// case <-sendLogToSplunkChan:
	// 	fmt.Println("Log sent to Splunk successfully.")
	// }
}

// GetSplunkInformationFromEnvironment
// SPLUNK_HOST: "https://{your-splunk-URL}:8088/services/collector",
// SPLUNK_TOKEN: "{your-token}",
// SPLUNK_SOURCE: "{your-source}",
// SPLUNK_SOURCETYPE: "{your-sourcetype}",
// SPLUNK_INDEX: "{your-index}",
func GetSplunkInformationFromEnvironment() (host string, token string, source string, sourcetype string, index string, splunkInfoIsFullSetInEnv bool) {
	var splunkHost, isSplunkHostSet = os.LookupEnv("SPLUNK_HOST")
	var splunkToken, isSplunkTokenSet = os.LookupEnv("SPLUNK_TOKEN")
	var splunkSource, isSplunkSourceSet = os.LookupEnv("SPLUNK_SOURCE")
	var splunkSourcetype, isSplunkSourcetypeSet = os.LookupEnv("SPLUNK_SOURCETYPE")
	var splunkIndex, isSplunkIndexSet = os.LookupEnv("SPLUNK_INDEX")
	if !isSplunkHostSet && !isSplunkTokenSet && !isSplunkSourceSet && !isSplunkSourcetypeSet && !isSplunkIndexSet {
		return "", "", "", "", "", false
	}
	return splunkHost, splunkToken, splunkSource, splunkSourcetype, splunkIndex, true
}

func AppendLogToFile(logContent string) error {
	folder := GetSystemRootFolder()
	timeZoneLocation, timeLoadLocationErr := time.LoadLocation("Asia/Ho_Chi_Minh")
	if timeLoadLocationErr != nil {
		return timeLoadLocationErr
	}
	currentTimestamp := time.Now().In(timeZoneLocation)

	logFileName := fmt.Sprintf(constant.LogFileLocation, currentTimestamp.Year(), int(currentTimestamp.Month()), currentTimestamp.Day())

	// check if logContent folder is existed or not
	if _, directoryStatusError := os.Stat(folder); os.IsNotExist(directoryStatusError) {
		log.Info(fmt.Sprintf("start to create folder %s", folder))
		makeDirectoryAllError := os.MkdirAll(folder, 0755)
		if makeDirectoryAllError != nil {
			return makeDirectoryAllError
		}
	}

	// +-----+---+--------------------------+
	// | rwx | 7 | Read, write and execute  |
	// | rw- | 6 | Read, write              |
	// | r-x | 5 | Read, and execute        |
	// | r-- | 4 | Read,                    |
	// | -wx | 3 | Write and execute        |
	// | -w- | 2 | Write                    |
	// | --x | 1 | Execute                  |
	// | --- | 0 | no permissions           |
	// +------------------------------------+

	// +------------+------+-------+
	// | Permission | Octal| Field |
	// +------------+------+-------+
	// | rwx------  | 0700 | User  |
	// | ---rwx---  | 0070 | Group |
	// | ------rwx  | 0007 | Other |
	// +------------+------+-------+
	// O_RDONLY: It opens the file read-only.
	// O_WRONLY: It opens the file write-only.
	// O_RDWR: It opens the file read-write.
	// O_APPEND: It appends data to the file when writing.
	// O_CREATE: It creates a new file if none exists.
	file, openFileError := os.OpenFile(folder+logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if openFileError != nil {
		return openFileError
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("An error has been occurred when defer closing file: %v", err)
		}
	}(file)

	_, writeStringToFileError := file.WriteString(logContent)

	if writeStringToFileError != nil {
		return writeStringToFileError
	}
	return nil
}

func GetSystemRootFolder() string {
	s := os.Getenv("MOUNT_FOLDER")
	return EnsureTrailingSlash(s)
}

func EnsureTrailingSlash(s string) string {
	if !strings.HasSuffix(s, "/") {
		s += "/"
	}
	return s
}

// CustomWriter is a custom log writer that uses a custom logging function
type CustomWriter struct {
}

// Implement the Write method of io.Writer for CustomWriter
func (cw CustomWriter) Write(p []byte) (n int, err error) {
	WithLevel(constant.Debug, context.Background(), string(p))
	return len(p), nil
}
