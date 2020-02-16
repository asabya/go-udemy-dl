package cmd

import (
	"github.com/Sab94/go-udemy-dl/core"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-udemy-dl",
	Short: "Go Udemy Downloader",
	Long: `
   Go Udemy Downloader is another udemy downloader written in go for
	 personal offline use.
	`,
}

func Execute() error {
	dl := core.New()
	initLogin(dl)
	initLogout(dl)
	initList(dl)
	rootCmd.SilenceUsage = true
	return rootCmd.Execute()
}
