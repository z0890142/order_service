package main

import (
	"flag"
	"order_service/config"
	"order_service/internal/app"
	"os"
	"os/signal"
	"syscall"
)

var (
	// flagconf is the config flag.
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "./config.yaml", "config path, eg: -conf config.yaml")
}

func handleSignals(server *app.Application) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	server.GetLogger().Infof("signal %s received", <-sigs)
	server.Shutdown()
}

// @title Order Upload API
// @version 1.0
func main() {
	flag.Parse()
	config.LoadConf(flagconf, config.GetConfig())

	server := app.Default()

	server.AddInitHook(app.InitDatabaseHook)
	server.AddInitHook(app.InitGinApplicationHook)

	server.AddDestroyHook(app.DestroyGinApplicationHook)

	go handleSignals(server)
	server.Run()
}
