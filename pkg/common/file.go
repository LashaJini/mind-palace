package common

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"
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

func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return info.IsDir(), nil
}

func RemoveAllFiles(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := filepath.Join(dir, file.Name())

		if err := os.Remove(filePath); err != nil {
			return err
		}
	}

	return nil
}

type SQLTemplate struct {
	Namespace string
}

func NewSQLTemplates(schemas []string) []SQLTemplate {
	sqlTemplates := make([]SQLTemplate, 0, len(schemas))

	for _, schema := range schemas {
		sqlTemplates = append(sqlTemplates, SQLTemplate{Namespace: schema})
	}

	return sqlTemplates
}

func (s *SQLTemplate) Inject(sqlBuffer *bytes.Buffer, path string) error {
	f, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	tmpl, err := template.New("sql").Parse(string(f))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, s)
	if err != nil {
		return err
	}

	_, err = sqlBuffer.WriteString(buf.String())

	return err
}
