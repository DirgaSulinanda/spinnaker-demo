package web

import (
	"context"
	"fmt"
	htmltemplate "html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/ariefrahmansyah/spinnaker-demo/router"
	"github.com/ariefrahmansyah/spinnaker-demo/template"
	"github.com/cockroachdb/cmux"
	"golang.org/x/net/netutil"
)

func getAsset(name string) ([]byte, error) {
	file, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %s", name, err)
	}

	return file, nil
}

// Options for the web handler.
type Options struct {
	ListenAddress  string
	MaxConnections int

	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	Version *Version
}

// Handler serves various HTTP endpoints.
type Handler struct {
	logger   *log.Logger
	options  *Options
	template *template.Template
	router   *router.Router
}

// New initializes a new web handler.
func New(logger *log.Logger, options *Options) *Handler {
	if logger == nil {
		logger = log.New(ioutil.Discard, "", 0)
	}

	baseTemplate := "web/ui/templates/_base.html"
	templatePath := "web/ui/templates"
	templateFuncMap := htmltemplate.FuncMap{}
	template := template.New(baseTemplate, templatePath, getAsset)
	template.Funcs(templateFuncMap)

	handler := &Handler{
		logger:   logger,
		options:  options,
		template: template,
	}

	router := router.New()
	router.Get("/ping", handler.Ping)
	router.Get("/version", handler.Version)

	router.Get("/", handler.Canary)
	router.Get("/canary", handler.Canary)

	handler.router = router

	return handler
}

// Run serves the web handler.
func (h *Handler) Run(ctx context.Context) error {
	log.Println("Running web server...")

	// Create the main listener
	listener, err := net.Listen("tcp", h.options.ListenAddress)
	if err != nil {
		return err
	}
	listener = netutil.LimitListener(listener, h.options.MaxConnections)

	// Listner multiplexer
	listenerMux := cmux.New(listener)
	httpListener := listenerMux.Match(cmux.HTTP1Fast())

	// TODO: grpc listener

	// HTTP handler
	httpHandler := &http.Server{
		Handler:     h.router,
		ReadTimeout: h.options.ReadTimeout,
		ErrorLog:    h.logger,
	}

	// Start listening
	errCh := make(chan error)
	go func() {
		errCh <- httpHandler.Serve(httpListener)
	}()
	go func() {
		errCh <- listenerMux.Serve()
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		httpHandler.Shutdown(ctx)
		return nil
	}
}

// Router of web handler.
func (h *Handler) Router() *router.Router {
	return h.router
}

// Ping writes pong.
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

// Canary renders canary page.
func (h *Handler) Canary(w http.ResponseWriter, req *http.Request) {
	backgroundColor := "deepskyblue"
	// backgroundColor = "springgreen"

	args := map[string]interface{}{
		"backgroundColor": backgroundColor,
	}
	h.template.ExecuteTemplate(w, "canary.html", args)
}
