package cmd

import (
	"errors"
	"log"

	"github.com/Sab94/go-udemy-dl/core"
	"github.com/spf13/cobra"
)

func initLogin(dl *core.Downloader) {
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Login to Udemy",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			// Parse username flag
			username, _ := cmd.Flags().GetString("username")

			// Parse password flag
			password, _ := cmd.Flags().GetString("password")

			log.Printf("Got username and password : %s %s\n", username, password)
			if username == "" && password == "" {
				err := errors.New("Username and password cannot be blank")
				return err
			}

			dl.GetLogin()
			dl.DoLogin(username, password)
			// dl.List()
			return nil
		},
	}
	rootCmd.AddCommand(loginCmd)
	loginCmd.PersistentFlags().StringP("username", "u", "", "username")
	loginCmd.PersistentFlags().StringP("password", "p", "", "password")
}
