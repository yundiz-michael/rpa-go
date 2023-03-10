package common

import (
	"bufio"
	"fmt"
	"os"
)

func Write0(fileName string, bytes []byte) (int, error) {
	var f *os.File
	var err error

	if CheckFileExist(fileName) {
		err = os.RemoveAll(fileName)
	}
	f, err = os.Create(fileName)
	if err != nil {
		fmt.Println("file create fail")
		return 0, err
	}
	write := bufio.NewWriter(f)
	nn, err := write.Write(bytes)
	if err != nil {
		return 0, err
	}
	err = write.Flush()
	if err != nil {
		return 0, err
	}
	return nn, nil
}

func CheckFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
