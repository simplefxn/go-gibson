package metrics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/simplefxn/go-gibson/pkg/config"
)

type Server struct {
	srv *http.Server
}

func init() {

}

func New(conf config.Metrics) *Server {
	mux := http.DefaultServeMux

	mux.Handle(conf.Path, promhttp.Handler())

	srv := &Server{
		srv: &http.Server{
			Addr:    fmt.Sprintf(":%v", conf.Port),
			Handler: mux,
		},
	}

	return srv
}

func (s *Server) ListenAndServe() {
	s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
