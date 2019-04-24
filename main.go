package main

import (
	"config"
	_ "controller"
	"fmt"
	"log"
	_ "logger"
	"os"
	"os/signal"
	_ "room"
	_ "rooms"
	"service"
	"syscall"
	"time"
)

func main() {

	var (
		port              = config.GetString("port")
		sm                = service.NewSessionManager()
		stopped chan bool = make(chan bool)
	)

	if port == "" {
		log.Println("config, port lost")
		fmt.Println("config error, port lost")
		os.Exit(1)
	}

	serv := service.Create(port)
	go func(stop chan<- bool) {
		serv.Start(sm)
		stop <- true
	}(stopped)

	sig_usr1 := make(chan os.Signal, 1)
	signal.Notify(sig_usr1, syscall.SIGUSR1)
	for {
		select {
		case <-stopped:
			log.Println("Service Stopped", time.Now())
		case <-sig_usr1:
			log.Println("Receive USR1")
		}
	}

}
