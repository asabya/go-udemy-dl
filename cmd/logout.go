package cmd

import (
	"github.com/Sab94/go-udemy-dl/core"
	"github.com/Sab94/go-udemy-dl/repo"
	"github.com/spf13/cobra"
)

func initLogout(dl *core.Downloader) {
	logoutCmd := &cobra.Command{
		Use:   "logout",
		Short: "Remove locally stored session",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := repo.RemoveRepo(dl.Root)
			if err != nil {
				cmd.Printf("Logout failed : %s\n", err.Error())
				return err
			}
			cmd.Println("Logout Successful")
			return nil
		},
	}
	rootCmd.AddCommand(logoutCmd)
}
