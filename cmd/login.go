package cmd

import (
	"errors"

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

			if username == "" && password == "" {
				err := errors.New("Username and password cannot be blank")
				return err
			}
			err := dl.GetLogin()
			if err != nil {
				cmd.Printf("Login failed : %s\n", err.Error())
				return err
			}
			err = dl.DoLogin(username, password)
			if err != nil {
				cmd.Printf("Login failed : %s\n", err.Error())
				return err
			}
			cmd.Println("Login Successful")
			return nil
		},
	}
	rootCmd.AddCommand(loginCmd)
	loginCmd.PersistentFlags().StringP("username", "u", "", "username")
	loginCmd.PersistentFlags().StringP("password", "p", "", "password")
}
