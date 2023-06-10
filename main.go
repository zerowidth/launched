package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "launchd",
	Short: "launchd web app server",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

var development bool
var listenAddress string

func init() {
	rootCmd.PersistentFlags().BoolVarP(&development, "development", "d", false, "run development mode to live-reload templates and static files")
	rootCmd.PersistentFlags().StringVarP(&listenAddress, "listen-address", "l", "localhost:3000", "address to listen on")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func serve() {
	logger, _ := zap.NewDevelopment()
	mux := http.NewServeMux()

	logger.Info("starting server", zap.String("listen-address", listenAddress), zap.Bool("development", development))
	if err := http.ListenAndServe(listenAddress, mux); err != nil {
		logger.Error("failed to start server", zap.Error(err))
	}
}
