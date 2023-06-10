package main

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

//go:embed static/* templates/*
var assets embed.FS

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
	defer logger.Sync()

	var staticFiles http.FileSystem
	if development {
		staticFiles = http.Dir(".")
	} else {
		staticFiles = http.FS(assets)
	}

	r := chi.NewRouter()
	r.Use(requestLogger(logger))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		file, err := assets.ReadFile("templates/layout.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(file)
	})
	r.Handle("/static/*", http.FileServer(staticFiles))

	logger.Info("starting server", zap.String("listen-address", listenAddress), zap.Bool("development", development))
	if err := http.ListenAndServe(listenAddress, r); err != nil {
		logger.Error("server error", zap.Error(err))
	}
}

func requestLogger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrapper := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()
			next.ServeHTTP(wrapper, r)

			logger.Info("request completed",
				zap.String("method", r.Method),
				zap.Int("status", wrapper.Status()),
				zap.Int("bytes", wrapper.BytesWritten()),
				zap.String("remote", r.RemoteAddr),
				zap.String("host", r.Host),
				zap.String("path", r.URL.Path),
				zap.String("user-agent", r.UserAgent()),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}
