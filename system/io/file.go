package io

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Exists(path string) bool {
	//os.Stat获取文件信息
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// 按行读取文件内容
func ReadFileByLine(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return ReadFileContentByLine(file)
}

// 按行读取文件内容
func ReadFileContentByLine(file *os.File) ([]string, error) {
	defer file.Close()
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	content := make([]string, 0)
	for fileScanner.Scan() {
		content = append(content, fileScanner.Text())
	}
	return content, nil
}

// EnsurePath is used to make sure a path exists
func EnsurePath(path string, dir bool) error {
	if !dir {
		path = filepath.Dir(path)
	}
	return os.MkdirAll(path, 0755)
}

// 读取文件
func ReadFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	return byteValue, err
}

// 将内容存入到文件中，自动换行
func WriteToFile(filePath string, data string) error {
	f, err := os.OpenFile(filePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(data + "\n"); err != nil {
		err = fmt.Errorf("写入数据到文件%s时出错,错误信息:%s", filePath, err.Error())
		return err
	}
	return nil
}

// 将内容存入到文件中，自动换行
func WriteLinesToFile(filePath string, lines []string) error {
	f, err := os.OpenFile(filePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, eachLine := range lines {
		if _, err := f.WriteString(eachLine + "\n"); err != nil {
			err = fmt.Errorf("写入数据到文件%s时出错,错误信息:%s", filePath, err.Error())
			return err
		}
	}
	return nil
}
