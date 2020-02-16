package cmd

import (
	"errors"
	"os"

	"github.com/Sab94/go-udemy-dl/core"
	"github.com/Sab94/go-udemy-dl/repo"
	"github.com/spf13/cobra"
)

func initList(dl *core.Downloader) {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List Subscribed Courses",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isLoggedin := repo.IsInitialized(dl.Root + string(os.PathSeparator) + "session")
			if isLoggedin {
				return nil
			}
			return errors.New("Please login to see list")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := dl.List()
			if err != nil {
				cmd.Printf("List failed : %s\n", err.Error())
				return err
			}
			return nil
		},
	}
	rootCmd.AddCommand(listCmd)
}
