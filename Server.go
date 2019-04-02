package main

import (
	"F22/config"
	"F22/db"
	"F22/handlers"
	"F22/route"
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"F22/err"
)

// Making sure all the defined interfaces
// abiding the contract
var (
	_ error = &err.UIError{}
)

func main() {

	//passing the config file with flag conf
	var configFile = flag.String("conf", "", "configuration file")
	flag.Parse()
	if flag.NFlag() != 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	//Configuring logger for the app
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	logger := log.New(os.Stdout, "F22: ", log.LstdFlags|log.Lshortfile)

	//Parsing configuration file
	logger.Println("Parsing configuration...")
	cfg := config.Parse(*configFile)


	// Connecting to DB
	logger.Println("Initializing mongodb master session...")
	db.ConnectDB()
	defer db.Obj.Close()

	//Provider holds application-wide variables
	logger.Println("Initializing provider...")
	Provider := handlers.NewProvider(logger, cfg, db.Obj)

	logger.Println("Initializing routes...")
	router := route.NewRouter(Provider)

	server := &http.Server{
		Addr:           cfg.HTTPAddress + ":" + strconv.Itoa(cfg.Port),
		Handler:        router,
		ReadTimeout:    time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Graceful shut down of server
	graceful := make(chan os.Signal)
	signal.Notify(graceful, syscall.SIGINT)
	signal.Notify(graceful, syscall.SIGTERM)
	go func() {
		<-graceful
		logger.Println("Shutting down server...")
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Fatalf("Could not do graceful shutdown: %v\n", err)
		}
	}()


	logger.Println("Listening server on ", server.Addr)
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("Listen: %s\n", err)
	}

	logger.Println("Server gracefully stopped")

}

