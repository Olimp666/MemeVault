package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/Olimp666/MemeVault/internal/app"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	) // ctrl+c, quit, stop, reload
	defer cancel()

	app := app.New()

	if err := app.Start(ctx); err != nil {
		log.Panicln("can't start application:", err)
	}

	if err := app.Wait(ctx, cancel); err != nil {
		log.Panicln("All systems closed with errors. LastError: ", err)
	}

	log.Println("All systems closed without errors")
}
