package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/config"
	"github.com/lashajini/mind-palace/pkg/models"
	rpcclient "github.com/lashajini/mind-palace/pkg/rpc/client"
	pb "github.com/lashajini/mind-palace/pkg/rpc/client/gen/proto"
	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/lashajini/mind-palace/pkg/storage/vdatabase"
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

	isText, err := isTextFile(file)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if !isText {
		fmt.Printf("Error: File %s is not a text file\n", file)
		os.Exit(1)
	}

	resourceID := uuid.New()
	fileExtension := filepath.Ext(file)
	dst := filepath.Join(config.MindPalaceOriginalResourcePath(currentUser, true), resourceID.String()+fileExtension)
	if err := copyFile(file, dst); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cfg := config.NewConfig()
	// TODO: move to py
	vdatabase.InitVDB(cfg)
	rpcClient := rpcclient.NewClient(cfg)
	db := database.InitDB(cfg)
	memory := models.NewMemory()

	ctx := context.Background()
	tx := database.NewMultiInstruction(ctx, db.DB())

	tx.Begin()
	memoryID, _ := models.InsertMemoryTx(tx, memory)
	resource := models.NewResource(resourceID, memoryID, config.MindPalaceOriginalResourcePath(currentUser, false))
	models.InsertResourceTx(tx, resource)
	tx.Commit()

	userCfg, _ := config.ReadUserConfig(currentUser)
	fmt.Println(userCfg)

	vectors, _ := rpcClient.Add(ctx, &pb.Memory{
		Steps: userCfg.Steps(),
		File:  dst,
		Type:  fileExtension,
	})
	fmt.Println(vectors)

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

func isTextFile(filename string) (bool, error) {
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

		if !isText(buf[:n]) {
			return false, nil
		}

		if err == io.EOF {
			break
		}
	}

	return true, nil
}

func isText(data []byte) bool {
	for len(data) > 0 {
		r, size := utf8.DecodeRune(data)
		if r == utf8.RuneError && size == 1 {
			return false
		}

		data = data[size:]
	}
	return true
}
