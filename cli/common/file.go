package common

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode/utf8"
)

func FileExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err == nil {
		return !stat.IsDir(), nil // File or dir
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	dstExists, err := FileExists(dst)
	if err != nil {
		return fmt.Errorf("failed to check if destination file exists: %w", err)
	}

	if dstExists {
		return fmt.Errorf("destination file already exists: %s", dst)
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	err = dstFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	return nil
}

func IsTextFile(filename string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buf := make([]byte, 4096)
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return false, err
		}

		if !IsText(buf[:n]) {
			return false, nil
		}

		if err == io.EOF {
			break
		}
	}

	return true, nil
}

func IsText(data []byte) bool {
	for len(data) > 0 {
		r, size := utf8.DecodeRune(data)
		if r == utf8.RuneError && size == 1 {
			return false
		}

		data = data[size:]
	}
	return true
}
