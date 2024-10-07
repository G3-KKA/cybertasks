package httpctl

import (
	"context"
	"cybertask/config"
	"cybertask/internal/logger"
	"errors"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type httpController struct {
	server *http.Server
	l      *logger.Logger
	cfg    config.HTTPServer

	shutdown sync.Once
}

func New(l *logger.Logger, cfg config.HTTPServer, mux *gin.Engine) *httpController {
	server := &http.Server{
		Addr:                         cfg.Address,
		Handler:                      mux,
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  0,
		ReadHeaderTimeout:            0,
		WriteTimeout:                 0,
		IdleTimeout:                  0,
		MaxHeaderBytes:               0,
		TLSNextProto:                 nil,
		ConnState:                    nil,
		ErrorLog:                     log.New(l, "", log.Flags()),
		BaseContext:                  nil,
		ConnContext:                  nil,
	}

	httpC := &httpController{
		server:   server,
		l:        l,
		cfg:      cfg,
		shutdown: sync.Once{},
	}

	return httpC
}

// Serve starts listening
func (ctl *httpController) Serve(ctx context.Context) error {

	lis, err := net.Listen("tcp", ctl.cfg.Address)
	if err != nil {
		return err
	}

	errs := make(chan error, 1)
	defer close(errs)

	ctl.l.Info().Msg("http server starting on: " + lis.Addr().String())

	select {
	case <-ctx.Done():
		err = ctl.Shutdown(ctx)
	case errs <- ctl.server.Serve(lis):
	}
	err2 := <-errs

	// If error is anything besides default ServerClosed.

	if !errors.Is(err2, http.ErrServerClosed) {
		return errors.Join(err, err2)
	}

	if err != nil {
		return err
	}

	return nil

}
func (ctl *httpController) Shutdown(ctx context.Context) error {

	var err error

	ctl.shutdown.Do(func() {
		err = ctl.server.Shutdown(ctx)
	})

	if err != nil {
		return err
	}

	return nil
}
