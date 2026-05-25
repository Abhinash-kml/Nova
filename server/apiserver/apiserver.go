package apiserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/abhinash-kml/nova/server/config"
)

type HttpServer struct {
	internalServer   http.Server
	beforeStartHooks []ServerHook
	afterStartHooks  []ServerHook
	beforeStopHooks  []ServerHook
	afterStopHooks   []ServerHook
	errChan          chan error
	ctx              context.Context
	cancel           context.CancelFunc
}

type ServerHook func(context.Context) error

func New(ctx context.Context, config config.HttpServerConfig, handler http.Handler) *HttpServer {
	ctx, cancel := context.WithCancel(ctx)
	return &HttpServer{
		internalServer: http.Server{
			Addr:           config.Address,
			ReadTimeout:    time.Duration(config.ReadTimeout),
			WriteTimeout:   time.Duration(config.WriteTimeout),
			IdleTimeout:    time.Duration(config.IdleTimeout),
			MaxHeaderBytes: config.MaxHeaderBytes,
			Handler:        handler,
		},
		errChan: make(chan error, 1),
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (s *HttpServer) Start() error {
	for _, function := range s.beforeStartHooks {
		err := function(context.Background())
		if err != nil {
			return err
		}
	}

	// Start the actual server in a separate goroutine and forward errors to the errChan
	go func() {
		s.errChan <- s.internalServer.ListenAndServe()
	}()

	for {
		select {
		case <-time.After(time.Second * 2):
			// Execute after start hooks in separate goroutine to prevent blocking
			go func() {
				for _, function := range s.afterStartHooks {
					err := function(context.Background())
					if err != nil {
						// Handle in some way
					}
				}
			}()

			// Monitor for errors in errChan after successful server start
			go func() {
				err := <-s.errChan
				if err != nil && err != http.ErrServerClosed {
					fmt.Println("Server crashed after start")
				}
			}()
		case err := <-s.errChan: // Server failed within 2 secs
			return err
		}
	}
}

func (s *HttpServer) Stop() {
	s.errChan <- s.internalServer.Shutdown(context.Background())
}

func (s *HttpServer) ListenForTermination() {

}

func (s *HttpServer) AddBeforeStartHook(function ServerHook) {
	s.beforeStartHooks = append(s.beforeStartHooks, function)
}

func (s *HttpServer) AddAfterStartHook(function ServerHook) {
	s.afterStartHooks = append(s.afterStartHooks, function)
}

func (s *HttpServer) AddBeforeStopHook(function ServerHook) {
	s.beforeStopHooks = append(s.beforeStopHooks, function)
}

func (s *HttpServer) AddAfterStopHook(function ServerHook) {
	s.afterStopHooks = append(s.afterStopHooks, function)
}
