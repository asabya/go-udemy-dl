package cmd

import (
	"github.com/Sab94/go-udemy-dl/core"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-udemy-dl",
	Short: "Go Udemy Downloader",
	Long: `
   Go Udemy Downloader is a udemy downloader for porsonal offline use.
	`,
}

func Execute() error {
	dl := core.New()
	initLogin(dl)
	initList(dl)
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	return rootCmd.Execute()
}
