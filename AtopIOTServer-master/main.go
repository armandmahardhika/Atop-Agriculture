package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	_ "github.com/austinjan/AtopIOTServer/config"
	"github.com/austinjan/AtopIOTServer/httpserver"
	"github.com/spf13/viper"
)

func main() {

	fmt.Println("Version: ", viper.GetString("version"))
	//init log system
	logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	//Server start
	log.Println("Server start")
	httpCtx, httpDone := context.WithCancel(context.Background())
	go httpserver.Run(httpCtx)

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	httpDone()

	log.Println("Server shutdown")
	os.Exit(0)
}
