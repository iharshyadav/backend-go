package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iharshyadav/backend/internal/config"
	"github.com/iharshyadav/backend/internal/http/handlers/handlerfunction"
	"github.com/iharshyadav/backend/internal/storage/sqlite"
)
 
func main() {
	// setup config
	cfg := config.Configuration()

	// setting up the database
	storage, err := sqlite.New(cfg)

	if err != nil {
		log.Fatal(err)
	}

	slog.Info("database connected", slog.String("env", cfg.Env), slog.Any("storage", storage))

	// router setup
	router := http.NewServeMux()

	router.HandleFunc("POST /api/user",handlerfunction.Create(storage))
	router.HandleFunc("GET /api/user/{id}",handlerfunction.GetById(storage))
	router.HandleFunc("GET /api/user",handlerfunction.GetList(storage))

	// setting up the server
	server := http.Server{
		Addr: cfg.Addr,
		Handler: router,
	}

	fmt.Printf("server is listening %s" , cfg.HTTPServer.Addr)

	done := make(chan os.Signal , 1)

	signal.Notify(done,os.Interrupt,syscall.SIGINT,syscall.SIGTERM)

	go func () {
		err := server.ListenAndServe()
        if err != nil {
			log.Fatal("Failed to start server")
		}
	} ()

	<-done

	slog.Info("shutting down the server")

	ctx , cancel := context.WithTimeout(context.Background(),5*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown the server",slog.String("error",err.Error()))
	}

	slog.Info("server shutdown successfully!!!")
	
}