package utils

import (
	"bufio"
	"context"
	"errors"
	"fast-storage-go-service/constant"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
)

func GetCurrentDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

func ReadFileFromPath(path ...string) []byte {
	if len(path) == 0 {
		return nil
	}
	resultPath := filepath.Join(path...)
	log.Info("Read file from path: %s\n", resultPath)
	file, err := os.Open(resultPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	stat, fileStatErr := file.Stat()
	if fileStatErr != nil {
		fmt.Println(fileStatErr)
		return nil
	}
	defer func(file *os.File) {
		fileCloseErr := file.Close()
		if fileCloseErr != nil {
			return
		}
	}(file)

	buffer := make([]byte, stat.Size())
	/*
		for {
			bytesRead, readFileErr := file.Read(buffer)
			if readFileErr != nil {
				if readFileErr != io.EOF {
					fmt.Println(readFileErr)
				}
				break
			}
			fmt.Println(string(buffer[:bytesRead])) // Print content from buffer
		}
	*/
	_, bufioReadErr := bufio.NewReader(file).Read(buffer)
	if bufioReadErr != nil && bufioReadErr != io.EOF {
		fmt.Println(bufioReadErr)
		return nil
	}
	return buffer
}

func FileEncryption(ctx context.Context, filePathToEncrypt string) error {
	command := fmt.Sprintf(constant.PythonEncryptFileCommand, constant.FileCryptoSecretKeyPath, filePathToEncrypt)
	if _, fileEncryptionStderr, fileEncryptionError := Shellout(ctx, command); fileEncryptionError != nil || fileEncryptionStderr != "" {
		return errors.New(fmt.Sprint(fileEncryptionError, fileEncryptionStderr))
	} else {
		return nil
	}
}

func FileDecryption(ctx context.Context, filePathToEncrypt string) error {
	command := fmt.Sprintf(constant.PythonDecryptFileCommand, constant.FileCryptoSecretKeyPath, filePathToEncrypt)
	if _, fileDecryptionStderr, fileDecryptionError := Shellout(ctx, command); fileDecryptionError != nil || fileDecryptionStderr != "" {
		return errors.New(fmt.Sprint(fileDecryptionError, fileDecryptionStderr))
	} else {
		return nil
	}
}

func ListAndFindFileInDirectory(ctx context.Context, path string, inputFileNameToFind ...string) ([]string, error) {
	result := make([]string, 0)
	info, err := os.Stat(path)
	if os.IsNotExist(err) || !info.IsDir() {
		return result, errors.New("input directory does not exist or it is not a directory")
	}

	isListAllFile := true
	if len(inputFileNameToFind) > 0 {
		isListAllFile = false
	}

	command := fmt.Sprintf(constant.PythonListAllFileCommand, path)
	if pythonListFileStdout, pythonListFileStderr, pythonListFileError := Shellout(ctx, command); pythonListFileError != nil || pythonListFileStderr != "" {
		return result, errors.New(fmt.Sprint(pythonListFileError, pythonListFileStderr))
	} else {
		if strings.Contains(pythonListFileStdout, "empty directory") {
			return result, nil
		}
		if isListAllFile {
			for _, currentLine := range strings.Split(pythonListFileStdout, "\n") {
				result = append(result, strings.TrimSpace(currentLine))
			}
		} else {
			contentToCompare := inputFileNameToFind[0]
			for _, currentLine := range strings.Split(pythonListFileStdout, "\n") {
				trimSpaceCurrentLine := strings.TrimSpace(currentLine)
				currentLineFileName := filepath.Base(trimSpaceCurrentLine)
				if strings.Contains(currentLineFileName, contentToCompare) || currentLineFileName == contentToCompare {
					result = append(result, trimSpaceCurrentLine)
				}
			}
		}
	}
	return result, nil
}

func ConvertImageIntoWebpBase64(imagePath string) (string, error) {
	command := fmt.Sprintf(constant.PythonImageReaderCommand, imagePath)
	if pythonConvertImageStdout, pythonConvertImageStderr, pythonConvertImageError := Shellout(context.Background(), command); pythonConvertImageError != nil || pythonConvertImageStderr != "" {
		return "", errors.New(fmt.Sprint(pythonConvertImageError, pythonConvertImageStderr))
	} else {
		convertImageResultArray := strings.Split(pythonConvertImageStdout, "\n")
		if convertImageResultArray[0] == "false" {
			return "", errors.New("can not convert image or it is not an image")
		} else {
			base64Result := convertImageResultArray[1]
			return base64Result, nil
		}
	}
}
