package config

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	*echo.Echo
	config  Configuration
	closers []func() // See `AddCloser` method.``
}

func NewServer(c Configuration) *Server {
	e := echo.New()
	srv := &Server{
		Echo:   e,
		config: c,
	}
	return srv
}

func (srv *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", srv.config.Host, srv.config.Port)
	log.Print("Server Running at :", addr)

	// BuildRouter
	srv.buildRouter()

	// Start server in a goroutine
	go func() {
		if err := srv.Echo.Start(addr); err != nil && err != http.ErrServerClosed {
			srv.Logger.Fatal("shutting down the server")
		}
	}()

	// Gracefull shutdown
	timeWait := 15 * time.Second
	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)
	<-signChan
	log.Println("Shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), timeWait)
	defer cancel()
	log.Println("Stop http server")
	if err := srv.Echo.Shutdown(ctx); err == context.DeadlineExceeded {
		log.Print("Halted active connections")
	} else if err != nil {
		log.Fatal(err)
	}
	close(signChan)
	log.Printf("Completed")
	return nil
}

// All root endpoints are registered here.
func (srv *Server) buildRouter() {
	srv.Any("/health", func(c echo.Context) error {
		return c.JSON(200, "OK")
	})
}

// registers application-level middlewares.
func (srv *Server) registerMiddlewares() {
	if srv.config.EnableRequestLog {
		// srv.Logger()
	}
	srv.Use(middleware.CORS())
}
