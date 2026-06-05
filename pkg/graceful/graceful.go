package graceful

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server interface {
	Run() error
	Shutdown(context.Context) error
	Timeout() time.Duration
}

func StartServerWithGracefulShutdown(server Server) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)

	notifyCh := make(chan os.Signal, 1)
	signal.Notify(notifyCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(notifyCh)

	go func(s Server) {
		if err := s.Run(); err != nil {
			err = fmt.Errorf("error starting server: %w", err)
			errCh <- err
		}
	}(server)

	select {
	case err := <-errCh:
		if err == nil {
			return nil
		}

		log.Println("Server start error:", err)
		return err
	case <-notifyCh:
		log.Println("Shutdown signal received")
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, server.Timeout())
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Println("Error during server shutdown:", err)
			return err
		}

		err := <-errCh
		// Ignore the expected ErrServerClosed error during shutdown
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println("Server returned error after shutdown:", err)
			return err
		}

		log.Println("Server successfully shutdown")
		return nil
	}
}
