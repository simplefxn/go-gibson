package metrics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/simplefxn/go-gibson/pkg/config"
	"github.com/simplefxn/go-gibson/pkg/logger"
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
	logger.Log.Infof("DEBUG %v", s)
	if s != nil {
		err := s.srv.Shutdown(ctx)
		return err
	}
	return fmt.Errorf("server has not stareted yet")
}
