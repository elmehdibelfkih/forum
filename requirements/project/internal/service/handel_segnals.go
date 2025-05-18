package service

import (
	"context"
	"fmt"
	db "forum/internal/db"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func HandleSignals(server *http.Server) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)

	sig := <-signalChannel
	fmt.Printf("\b\b\033[33mReceived signal: %s\n\033[0m", sig)
	fmt.Println("\033[33mShutting down gracefully...\033[0m")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	// TODO: close any other ersources
	db.CloseDB()

	print("\033[31mshutdowning... \033[0m")
	// for range 4 {
	// 	print("\b\b\b\b    ")
	// 	time.Sleep(100 * time.Millisecond)
	// 	print("\b\b\b\b\033[31m... \033[0m")
	// 	time.Sleep(200 * time.Millisecond)
	// }
	println()
	fmt.Println("\033[32mExited cleanly.\033[31m")
}
