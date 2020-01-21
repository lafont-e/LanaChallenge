package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/elafont/CbreChallenge/hangman"
	"github.com/elafont/CbreChallenge/server"
	"github.com/gorilla/mux"
)

const DEFAULTHOST = "localhost"
const DEFAULTWEBPORT = "8080"

var (
	ipWeb = flag.String("ip", DEFAULTHOST+":"+DEFAULTWEBPORT, "Web IP:PORT used to listen ie: *:8080, :8080, localhost")
	help  = flag.Bool("help", false, "Print Usage options")
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	flag.Usage = Usage
	flag.Parse()

	if *help {
		Usage()
	}

	go signals()

	// Create Server
	s := &server.Server{
		Router: mux.NewRouter(),
		Games:  make([]*hangman.Hangman, 0, 8),
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}

	s.RegisterRoutes()

	log.Fatal(s.Start(*ipWeb))
}

func signals() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGQUIT) // Kill signal not needed as its handled by the OS
	sig := <-sigCh
	log.Printf("Signal received %v\n", sig)
	fmt.Fprintf(os.Stderr, "Signal received %v\n", sig)
	os.Exit(1)
}

func Usage() {
	fmt.Println("Usage: ", os.Args[0], "[-ip, -help]")
	fmt.Println("   ie: ", os.Args[0], "-ip *:8081")
	os.Exit(2)
}
