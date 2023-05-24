package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := App{}
	app.Initialize()
	app.PopulateDatabase()

	// this is to listen to the interrupt signal and close the database connection gracefully
	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		<-gracefulStop
		app.CloseDatabase()
		os.Exit(0)
	}()

	app.Run("localhost:10000")
}
