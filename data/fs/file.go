package fs

import (
	"io/ioutil"
	"strings"
)

func ReadFile(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	src := string(data)
	if err != nil {
		println(err)
		return "", err
	}
	return src, nil
}

func GetParentDir(filePath string) string {
	dirPaths := strings.Split(filePath, "/")
	return strings.Join(dirPaths[:len(dirPaths)-1], "/")
}

func BReadFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		println(err)
		return nil, err
	}
	return data, nil
}
func ReadLines(filePath string) []string {
	src, err := ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return strings.Split(src, "\n")
}

func ReadCsv(filePath string, separator string) ([]string, int) {
	src, err := ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	arr := strings.Split(src, "\n")
	return strings.Split(strings.Join(arr, separator), separator), len(arr[0])
}
