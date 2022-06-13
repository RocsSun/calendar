package utils

import (
	"errors"
	"os"
)

func ReadFile(name string) (*os.File, error) {
	if IsFile(name) {
		return os.Open(name)
	}
	return nil, errors.New("文件不存在。")
}

func IsFile(name string) bool {
	if name == "" {
		return false
	}
	info, err := os.Stat(name)
	if err != nil {
		return false
	}

	if !info.IsDir() {
		return true
	}
	return false
}
