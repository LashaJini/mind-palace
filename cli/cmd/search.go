package cli

import (
	"fmt"

	"github.com/lashajini/mind-palace/pkg/api"
	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/mperrors"
	vdbrpc "github.com/lashajini/mind-palace/pkg/rpc/vdb"
	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/spf13/cobra"
)

var SEARCH_TEXT string

var searchCmd = &cobra.Command{
	Use:   "search -t [TEXT]",
	Short: "search <short description>",
	Long:  "search <long description>",
	Run:   Search,
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringVarP(&SEARCH_TEXT, "text", "t", "", "text to search")

	_ = searchCmd.MarkFlagRequired("text")
}

func Search(cmd *cobra.Command, args []string) {
	text, _ := cmd.Flags().GetString("text")

	cfg := common.NewConfig()
	currentUser, err := common.CurrentUser()
	mperrors.On(err).Wrap("failed to get current user").Exit()

	db := database.InitDB(cfg)
	db.SetSchema(db.ConstructSchema(currentUser))
	vdbGrpcClient := vdbrpc.NewGrpcClient(cfg, currentUser)

	searchResponse, err := api.Search(cmd.Context(), text, db, vdbGrpcClient)
	mperrors.On(err).Exit()

	fmt.Println(searchResponse)
}
