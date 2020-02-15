package cmd

import (
	"github.com/Sab94/go-udemy-dl/core"
	"github.com/spf13/cobra"
)

func initList(dl *core.Downloader) {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List Subscribed Cources",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			dl.List()
			return nil
		},
	}
	rootCmd.AddCommand(listCmd)
}
