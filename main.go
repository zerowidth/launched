package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "launchd",
	Short: "launchd web app server",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

var development = false

func init() {
	rootCmd.PersistentFlags().BoolVarP(&development, "development", "d", false, "run development mode to live-reload templates/static files")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func serve() {
	fmt.Println("Hello, world!")
}
