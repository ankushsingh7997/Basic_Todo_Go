package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	todoRouter "github.com/ankush/todo/Router"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/thedevsaddam/renderer"
)

var rnd *renderer.Render

// so using pointers - so a pointer is like a sign post that tells where something lives in memory insted of copying the whole object

// db is not the database itself but a pointer to the database. It's like saying, "Here's the address of the database, go there to find it."

// Why use * (pointers) here?

// Efficient Memory Use: Instead of copying the whole database object (which could be big), we just pass around its address, which is small and fast.

// Shared Access: When different parts of your program use db, they're all looking at the same database. If one part makes changes, everyone sees those changes.

const port string = ":3001"

type todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"createdAt"`
}

func main() {
	// Close server gracefully
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/api/v1/todo", todoRouter.TodoHandlers())

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Listening on port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Error while connecting to server: %v\n", err)
		}
	}()

	<-stopChan
	log.Println("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("Server gracefully stopped")
}
