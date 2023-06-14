package main

import (
	"embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/form"
	"github.com/go-playground/validator/v10"
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

// single instance caches structs
var decoder *form.Decoder

func init() {
	rootCmd.PersistentFlags().BoolVarP(&development, "development", "d", false, "run development mode to live-reload templates and static files")
	rootCmd.PersistentFlags().StringVarP(&listenAddress, "listen-address", "l", "localhost:3000", "address to listen on")

	decoder = form.NewDecoder()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type PlistField struct {
	Value string
	Error string
}

type PlistForm map[string]PlistField

func NewPlistForm(plist LaunchdPlist, errors validator.ValidationErrorsTranslations) PlistForm {
	form := PlistForm{}
	v := reflect.ValueOf(plist)
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		value := v.Field(i).Interface().(string)
		form[field.Name] = PlistField{
			Value: fmt.Sprintf("%v", value),
			Error: errors[v.Type().Name()+"."+field.Name],
		}
	}
	return form
}

func serve() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	var staticFiles fs.FS
	if development {
		staticFiles = os.DirFS(".")
	} else {
		staticFiles = assets
	}

	r := chi.NewRouter()
	r.Use(requestLogger(logger))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		form := PlistForm{}
		layout := template.Must(template.ParseFS(staticFiles, "templates/layout.html", "templates/form.html"))
		layout.Execute(w, form)
	})

	r.Post("/plist", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		plist := NewPlistFromForm(r.PostForm)
		errors := plist.Validate()
		if errors != nil {
			form := NewPlistForm(plist, errors)
			layout := template.Must(template.ParseFS(staticFiles, "templates/layout.html", "templates/form.html"))
			layout.Execute(w, form)
			return
		}
		http.Redirect(w, r, "/plist/"+plist.Encode(), http.StatusSeeOther)
	})

	r.Get("/plist/{encoded}", func(w http.ResponseWriter, r *http.Request) {
		proto := r.Header.Get("X-Forwarded-Proto")
		if proto == "" {
			proto = "http"
		}
		host := r.Host
		url := fmt.Sprintf("%s://%s", proto, host)
		decoded, _ := base64.RawURLEncoding.DecodeString(chi.URLParam(r, "encoded"))
		plist := NewPlistFromJSON(string(decoded))
		context := struct {
			Plist   LaunchdPlist
			RootURL string
		}{
			Plist:   plist,
			RootURL: url,
		}
		layout := template.Must(template.ParseFS(staticFiles, "templates/layout.html", "templates/plist.html"))
		layout.Execute(w, context)
	})

	r.Get("/plist/{encoded}/install", func(w http.ResponseWriter, r *http.Request) {
		proto := r.Header.Get("X-Forwarded-Proto")
		if proto == "" {
			proto = "http"
		}
		host := r.Host
		url := fmt.Sprintf("%s://%s", proto, host)
		decoded, _ := base64.RawURLEncoding.DecodeString(chi.URLParam(r, "encoded"))
		plist := NewPlistFromJSON(string(decoded))
		context := struct {
			Plist   LaunchdPlist
			RootURL string
		}{
			Plist:   plist,
			RootURL: url,
		}
		layout := template.Must(template.ParseFS(staticFiles, "templates/install.sh"))
		r.Header.Set("Content-Type", "text/plain; charset=utf-8")
		layout.Execute(w, context)
	})

	r.Get("/plist/{encoded}.xml", func(w http.ResponseWriter, r *http.Request) {
		decoded, _ := base64.RawURLEncoding.DecodeString(chi.URLParam(r, "encoded"))
		plist := NewPlistFromJSON(string(decoded))
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(plist.PlistXML()))
	})

	r.Get("/plist/{encoded}/download", func(w http.ResponseWriter, r *http.Request) {
		decoded, _ := base64.RawURLEncoding.DecodeString(chi.URLParam(r, "encoded"))
		plist := NewPlistFromJSON(string(decoded))
		w.Header().Set("Content-Type", "application/xml")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.plist", plist.Label()))
		w.Write([]byte(plist.PlistXML()))
	})

	r.Get("/plist/{encoded}/edit", func(w http.ResponseWriter, r *http.Request) {
		decoded, _ := base64.RawURLEncoding.DecodeString(chi.URLParam(r, "encoded"))
		plist := NewPlistFromJSON(string(decoded))
		form := NewPlistForm(plist, nil)
		layout := template.Must(template.ParseFS(staticFiles, "templates/layout.html", "templates/form.html"))
		layout.Execute(w, form)
	})

	r.Get("/help", func(w http.ResponseWriter, r *http.Request) {
		layout := template.Must(template.ParseFS(staticFiles, "templates/layout.html", "templates/help.html"))
		layout.Execute(w, nil)
	})
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		ico, _ := fs.ReadFile(staticFiles, "static/favicon.ico")
		w.Write(ico)
	})

	r.Handle("/static/*", http.FileServer(http.FS(staticFiles)))

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

			remoteAddr := r.RemoteAddr
			if r.Header.Get("X-Forwarded-For") != "" {
				remoteAddr = strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]
			}

			logger.Info("request completed",
				zap.String("method", r.Method),
				zap.Int("status", wrapper.Status()),
				zap.Int("bytes", wrapper.BytesWritten()),
				zap.String("remote", remoteAddr),
				zap.String("host", r.Host),
				zap.String("url", r.URL.String()),
				zap.String("path", r.URL.Path),
				zap.String("user-agent", r.UserAgent()),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}
