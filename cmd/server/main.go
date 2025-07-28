package main

import (
	"context"
	"go-mini-server/core"
	"go-mini-server/core/db/pool"
	"go-mini-server/internal/web"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	Web web.Config  `yaml:"web"`
	DB  pool.Config `yaml:"db"`
}

const configPath = "./config.yaml"

func main() {
	defer recoverPanic()

	c, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	initServices(c)

	startServer(c)
}

func recoverPanic() {
	if message := recover(); message != nil {
		log.Printf("panic recovered: %v", message)
	}
}

func loadConfig() (c Config, err error) {
	err = core.ParseYaml(configPath, &c)
	return
}

func getHandlers() map[string]web.Handler {
	return map[string]web.Handler{
		"/user/get": handleUserGet,
	}
}

func startServer(c Config) {
	handlers := getHandlers()
	s := web.NewServer(c.Web, handlers,
		web.LoggerMiddleware,
	)

	srv := &http.Server{
		Addr:         c.Web.Listen,
		Handler:      s,
		ReadTimeout:  c.Web.Timeout,
		WriteTimeout: c.Web.Timeout,
	}

	go func() {
		log.Printf("server started at %s", c.Web.Listen)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}

	log.Println("server exited properly")
}
