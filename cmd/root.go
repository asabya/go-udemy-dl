package cmd

import (
	"log"

	"github.com/Sab94/go-udemy-dl/core"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-udemy-dl",
	Short: "Go Udemy Downloader",
	Long: `
   Go Udemy Downloader is a udemy downloader for porsonal offline use.
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		log.Println("Prerun")
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("run")
	},
}

func Execute() error {
	dl := core.New()
	initLogin(dl)
	return rootCmd.Execute()
}
