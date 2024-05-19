package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/config"
	"github.com/spf13/cobra"
)

var (
	FILE    string
	PREVIEW bool
)

var addCmd = &cobra.Command{
	Use:   "add -f [FILE]",
	Short: "add <short description>",
	Long:  "add <long description>",
	Run:   Add,
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&FILE, "file", "f", "", "file location")
	addCmd.MarkFlagRequired("file")
	addCmd.Flags().BoolVarP(&PREVIEW, "preview", "p", false, "preview result")
}

func Add(cmd *cobra.Command, args []string) {
	currentUser, err := config.CurrentUser()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if currentUser == "" {
		fmt.Println("Error: There are no users available.")
		fmt.Printf("\nCreate one by using: mind-palace user --new <name>\n\n")
		os.Exit(1)
	}

	file, _ := cmd.Flags().GetString("file")

	exists, err := fileExists(file)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if !exists {
		fmt.Printf("Error: File %s does not exist\n", file)
		os.Exit(1)
	}

	id := uuid.New()
	fileExtension := filepath.Ext(file)
	dst := filepath.Join(config.MindPalaceOriginalResourcePath(currentUser), id.String()+fileExtension)
	if err := copyFile(file, dst); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	add(args...)
}

func add(args ...string) {}

func fileExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err == nil {
		return !stat.IsDir(), nil // File or dir
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

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
