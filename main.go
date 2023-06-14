package main

import (
	"embed"
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
var redisAddress string

//go:embed static templates
var assets embed.FS

// single instance caches structs
var decoder *form.Decoder

func init() {
	rootCmd.PersistentFlags().BoolVarP(&development, "development", "d", false, "run development mode to live-reload templates and static files")
	rootCmd.PersistentFlags().StringVarP(&listenAddress, "listen-address", "l", "localhost:3000", "address to listen on")
	rootCmd.PersistentFlags().StringVarP(&redisAddress, "redis-address", "r", "localhost:6379", "address of redis server")

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

	store := NewPlistStore(redisAddress, os.Getenv("REDIS_PASSWORD"))

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
		id, err := store.Save(plist)
		if err != nil {
			logger.Error("error saving plist", zap.Error(err))
			http.Error(w, "could not save plist", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/plist/"+id, http.StatusSeeOther)
	})

	r.Get("/plist/{id}", func(w http.ResponseWriter, r *http.Request) {
		proto := r.Header.Get("X-Forwarded-Proto")
		if proto == "" {
			proto = "http"
		}
		host := r.Host
		url := fmt.Sprintf("%s://%s", proto, host)

		plist, ok, err := store.Load(chi.URLParam(r, "id"))
		if err != nil {
			logger.Error("error loading plist", zap.Error(err))
			http.Error(w, "could not load plist", http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

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

	r.Get("/plist/{id}/install", func(w http.ResponseWriter, r *http.Request) {
		proto := r.Header.Get("X-Forwarded-Proto")
		if proto == "" {
			proto = "http"
		}
		host := r.Host
		url := fmt.Sprintf("%s://%s", proto, host)

		plist, ok, err := store.Load(chi.URLParam(r, "id"))
		if err != nil {
			logger.Error("error loading plist", zap.Error(err))
			http.Error(w, "could not load plist", http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

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

	r.Get("/plist/{id}.xml", func(w http.ResponseWriter, r *http.Request) {
		plist, ok, err := store.Load(chi.URLParam(r, "id"))
		if err != nil {
			logger.Error("error loading plist", zap.Error(err))
			http.Error(w, "could not load plist", http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(plist.PlistXML()))
	})

	r.Get("/plist/{id}/download", func(w http.ResponseWriter, r *http.Request) {
		plist, ok, err := store.Load(chi.URLParam(r, "id"))
		if err != nil {
			logger.Error("error loading plist", zap.Error(err))
			http.Error(w, "could not load plist", http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/xml")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.plist", plist.Label()))
		w.Write([]byte(plist.PlistXML()))
	})

	r.Get("/plist/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
		plist, ok, err := store.Load(chi.URLParam(r, "id"))
		if err != nil {
			logger.Error("error loading plist", zap.Error(err))
			http.Error(w, "could not load plist", http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

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
