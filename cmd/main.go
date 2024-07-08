package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/soltanireza65/gp-ts/internal/handlers"
	"github.com/soltanireza65/gp-ts/internal/store"
)

func main() {
	todos := []store.Todo{
		{Title: "Task 1", Done: false},
		{Title: "Task 2", Done: true},
		{Title: "Task 3", Done: false},
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/healthcheck", handlers.NewHealthcheckHandler().Execute)

	r.Post("/todos", handlers.NewCreateTodoHandler(handlers.CreateTodoHandlerParams{Todos: &todos}).Execute)

	r.Get("/todos", handlers.NewGetTodosHandler(handlers.GetTodosHandlerParams{Todos: &todos}).Execute)

	server := &http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Println("Starting server on http://localhost:8000")

		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("Error starting server: %v\n", err)

		}
	}()

	fmt.Println("Press Ctrl+C to stop the server")

	<-sigCh

	fmt.Println("Stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Error stopping server: %v\n", err)
	}

	fmt.Println("Server stopped")
}
