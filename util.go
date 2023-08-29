package wordCloudGenerator

import (
	"os"
	"path/filepath"
)

var currentDirectory string

func init() {
	executablePath, _ := os.Executable()
	currentDirectory = filepath.Dir(executablePath)
}

func fileNameToByteArray(fileName string) ([]byte, error) {
	//check if file exists
	if _, err := os.Stat(fileName); err != nil {
		return nil, err
	}

	//open file
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	//read file
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	//create byte array
	fileByteArray := make([]byte, fileInfo.Size())

	//read file into byte array
	_, err = file.Read(fileByteArray)
	if err != nil {
		return nil, err
	}

	return fileByteArray, nil
}
