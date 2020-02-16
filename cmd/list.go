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
