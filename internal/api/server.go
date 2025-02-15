package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (h *Handler) Serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", h.config.Port),
		Handler:      h.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownErrorChannel := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		h.logger.Info().Msgf("caught %s", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownErrorChannel <- err
		}

		h.logger.Info().Msgf("completing background tasks at addr: %s ", srv.Addr)

		h.wg.Wait()

		close(shutdownErrorChannel)
	}()

	h.logger.Info().Msgf("starting server %s based on %s", srv.Addr, h.config.Env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	for err := range shutdownErrorChannel {
		if err != nil {
			return err
		}
	}

	h.logger.Info().Msgf("stopped server %s %s", srv.Addr, h.config.Env)

	return nil
}
