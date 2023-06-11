package main

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"text/template"
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

//go:embed static templates
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

	var fs fs.FS
	if development {
		fs = os.DirFS(".")
	} else {
		fs = assets
	}

	r := chi.NewRouter()
	r.Use(requestLogger(logger))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		layout := template.Must(template.ParseFS(fs, "templates/layout.html", "templates/form.html"))
		layout.Execute(w, nil)
	})
	r.Post("/plist", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		plist := createPlistFromForm(r.PostForm)
		http.Redirect(w, r, "/plist/"+plist, http.StatusSeeOther)
	})
	r.Get("/plist/{encoded}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		decoded, _ := base64.RawURLEncoding.DecodeString(chi.URLParam(r, "encoded"))
		w.Write(decoded)
	})
	r.Handle("/static/*", http.FileServer(http.FS(fs)))

	logger.Info("starting server", zap.String("listen-address", listenAddress), zap.Bool("development", development))
	if err := http.ListenAndServe(listenAddress, r); err != nil {
		logger.Error("server error", zap.Error(err))
	}
}

// JSON encoded launchd plist, for encoded use in URL path.
// Uses string types to easily allow for empty values.
type LaunchdPlist struct {
	Name              string `json:"name,omitempty"`
	Command           string `json:"command,omitempty"`
	StartInterval     string `json:"start_interval,omitempty"`
	Minute            string `json:"minute,omitempty"`
	Hour              string `json:"hour,omitempty"`
	DayOfMonth        string `json:"day_of_month,omitempty"`
	Month             string `json:"month,omitempty"`
	Weekday           string `json:"weekday,omitempty"`
	RunAtLoad         string `json:"run_at_load,omitempty"`
	RestartOnCrash    string `json:"restart_on_crash,omitempty"`
	StartOnMount      string `json:"start_on_mount,omitempty"`
	QueueDirectories  string `json:"queue_directories,omitempty"`
	Environment       string `json:"environment,omitempty"`
	User              string `json:"user,omitempty"`
	Group             string `json:"group,omitempty"`
	WorkingDirectory  string `json:"working_directory,omitempty"`
	StandardOutPath   string `json:"standard_out_path,omitempty"`
	StandardErrorPath string `json:"standard_error_path,omitempty"`
}

func createPlistFromForm(values url.Values) string {
	// parse the values into a LaunchdPlist
	plist := LaunchdPlist{
		Name:              values.Get("name"),
		Command:           values.Get("command"),
		StartInterval:     values.Get("start_interval"),
		Minute:            values.Get("minute"),
		Hour:              values.Get("hour"),
		DayOfMonth:        values.Get("day_of_month"),
		Month:             values.Get("month"),
		Weekday:           values.Get("weekday"),
		RunAtLoad:         values.Get("run_at_load"),
		RestartOnCrash:    values.Get("restart_on_crash"),
		StartOnMount:      values.Get("start_on_mount"),
		QueueDirectories:  values.Get("queue_directories"),
		Environment:       values.Get("environment"),
		User:              values.Get("user"),
		Group:             values.Get("group"),
		WorkingDirectory:  values.Get("working_directory"),
		StandardOutPath:   values.Get("standard_out_path"),
		StandardErrorPath: values.Get("standard_error_path"),
	}

	encoded, _ := json.Marshal(plist)
	return base64.RawURLEncoding.EncodeToString(encoded)
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
