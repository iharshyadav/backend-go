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
)

func main() {
	// setup config
	cfg := config.Configuration()

	// router setup
	router := http.NewServeMux()

	router.HandleFunc("GET /",func(w http.ResponseWriter , r *http.Request){
		w.Write([]byte("Welcome to backend server"))
	})


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