package apiserver

import (
	"context"
	"net/http"
	"time"

	"github.com/abhinash-kml/nova/server/config"
	"go.uber.org/zap"
)

type HttpServer struct {
	internalServer   http.Server
	beforeStartHooks []Hook
	afterStartHooks  []Hook
	beforeStopHooks  []Hook
	afterStopHooks   []Hook
	errChan          chan error
	ctx              context.Context
	cancel           context.CancelFunc
	logger           *zap.Logger
}

type Hook func(context.Context) error

func New(ctx context.Context, config config.HttpServerConfig, handler http.Handler, logger *zap.Logger) *HttpServer {
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
		logger:  logger,
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
						s.logger.Error("Failed to execute after start hooks", zap.Error(err))
					}
				}
			}()

			// Monitor for errors in errChan after successful server start
			go func() {
				err := <-s.errChan
				if err != nil && err != http.ErrServerClosed {
					s.logger.Error("Server crashed after start", zap.Error(err))
				}
			}()

			// Start listening for termination
			go s.listenForTermination()
		case err := <-s.errChan: // Server failed within 2 secs
			return err
		}
	}
}

func (s *HttpServer) Stop(ctx context.Context) {
	s.errChan <- s.internalServer.Shutdown(ctx)
}

func (s *HttpServer) listenForTermination() {
	<-s.ctx.Done() // This will block untill context is done
	s.logger.Info("Server terminated as a result of context completion")
	s.Stop(context.Background())
}

func (s *HttpServer) AddBeforeStartHook(function Hook) {
	s.beforeStartHooks = append(s.beforeStartHooks, function)
}

func (s *HttpServer) AddAfterStartHook(function Hook) {
	s.afterStartHooks = append(s.afterStartHooks, function)
}

func (s *HttpServer) AddBeforeStopHook(function Hook) {
	s.beforeStopHooks = append(s.beforeStopHooks, function)
}

func (s *HttpServer) AddAfterStopHook(function Hook) {
	s.afterStopHooks = append(s.afterStopHooks, function)
}
