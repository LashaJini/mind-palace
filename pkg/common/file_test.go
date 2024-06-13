package common

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FileExists(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() (string, error)
		expected  bool
		shouldErr bool
	}{
		{
			name: "File exists",
			setup: func() (string, error) {
				file, err := os.CreateTemp("", "testfile")
				if err != nil {
					return "", err
				}
				file.Close()
				return file.Name(), nil
			},
			expected:  true,
			shouldErr: false,
		},
		{
			name: "File does not exist",
			setup: func() (string, error) {
				return "non_existent_file", nil
			},
			expected:  false,
			shouldErr: false,
		},
		{
			name: "Directory exists",
			setup: func() (string, error) {
				dir, err := os.MkdirTemp("", "testdir")
				if err != nil {
					return "", err
				}
				return dir, nil
			},
			expected:  false,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			path, err := tt.setup()
			if err != nil && !tt.shouldErr {
				t.Fatalf("Failed to set up test: %v", err)
			}

			exists, err := FileExists(path)
			if tt.shouldErr {
				assert.Error(err, "Expected an error but got none")
			} else {
				assert.NoError(err, "Did not expect an error but got one")
				assert.Equal(tt.expected, exists, "Expected file existence does not match")
			}

			// Clean up
			if path != "" && !tt.shouldErr {
				os.Remove(path)
			}
		})
	}
}

func Test_CopyFile(t *testing.T) {
	tests := []struct {
		name       string
		srcContent []byte
		expected   string
		setupFunc  func() (string, string, error)
		shouldFail bool
	}{
		{
			name:       "Successful copy",
			srcContent: []byte("Hello, World!"),
			expected:   "Hello, World!",
			setupFunc: func() (string, string, error) {
				srcFile, err := os.CreateTemp("", "src")
				if err != nil {
					return "", "", err
				}
				dstFile, err := os.CreateTemp("", "dst")
				if err != nil {
					return "", "", err
				}
				os.Remove(dstFile.Name()) // Ensure dst file does not exist
				return srcFile.Name(), dstFile.Name(), nil
			},
			shouldFail: false,
		},
		{
			name: "Source file does not exist",
			setupFunc: func() (string, string, error) {
				dstFile, err := os.CreateTemp("", "dst")
				if err != nil {
					return "", "", err
				}
				return "non_existent_file", dstFile.Name(), nil
			},
			shouldFail: true,
		},
		{
			name:       "Destination file already exists",
			srcContent: []byte("Hello, World!"),
			setupFunc: func() (string, string, error) {
				srcFile, err := os.CreateTemp("", "src")
				if err != nil {
					return "", "", err
				}
				dstFile, err := os.CreateTemp("", "dst")
				if err != nil {
					return "", "", err
				}
				return srcFile.Name(), dstFile.Name(), nil
			},
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			src, dst, err := tt.setupFunc()
			assert.NoError(err, "Failed to set up test")

			// If we have source content to write, do it now
			if tt.srcContent != nil {
				err := os.WriteFile(src, tt.srcContent, 0644)
				assert.NoError(err, "Failed to write source file")
			}

			err = CopyFile(src, dst)
			if (err != nil) != tt.shouldFail {
				assert.Error(err, "Expected an error but got none")
			}

			if !tt.shouldFail {
				assert.NoError(err, "Did not expect an error but got one")

				// Verify the contents of the destination file
				dstContent, err := os.ReadFile(dst)
				assert.NoError(err, "Failed to read destination file")
				assert.Equal(tt.expected, string(dstContent), "Expected content does not match")
			}

			os.Remove(src)
			os.Remove(dst)
		})
	}
}

func Test_IsTextFile(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected bool
	}{
		{
			name:     "Valid text file",
			content:  []byte("Hello, 世界"),
			expected: true,
		},
		{
			name:     "Binary file",
			content:  []byte{0xff, 0xfe, 0xfd},
			expected: false,
		},
		{
			name:     "Empty file",
			content:  []byte{},
			expected: true, // Empty file is considered valid text
		},
		{
			name:     "Single valid rune",
			content:  []byte{0xe4, 0xb8, 0xad}, // 中
			expected: true,
		},
		{
			name:     "Single invalid rune",
			content:  []byte{0xe4, 0xb8, 0x00},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, err := os.CreateTemp("", "testfile")
			if err != nil {
				t.Fatalf("Failed to create temporary file: %v", err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write(tt.content); err != nil {
				t.Fatalf("Failed to write to temporary file: %v", err)
			}

			// Close the file so it can be reopened by the function
			if err := tmpfile.Close(); err != nil {
				t.Fatalf("Failed to close temporary file: %v", err)
			}

			result, err := IsTextFile(tmpfile.Name())
			if err != nil {
				t.Fatalf("IsTextFile(%q) returned error: %v", tmpfile.Name(), err)
			}
			if result != tt.expected {
				t.Errorf("IsTextFile(%q) = %v; want %v", tmpfile.Name(), result, tt.expected)
			}
		})
	}
}

func Test_IsText(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected bool
	}{
		{
			name:     "Valid UTF-8",
			input:    []byte("Hello, 世界"),
			expected: true,
		},
		{
			name:     "Invalid UTF-8",
			input:    []byte{0xff, 0xfe, 0xfd},
			expected: false,
		},
		{
			name:     "Empty input",
			input:    []byte{},
			expected: true, // Empty input is considered valid text
		},
		{
			name:     "Single valid rune",
			input:    []byte{0xe4, 0xb8, 0xad}, // 中
			expected: true,
		},
		{
			name:     "Single invalid rune",
			input:    []byte{0xe4, 0xb8, 0x00},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsText(tt.input)
			if result != tt.expected {
				t.Errorf("IsText(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
