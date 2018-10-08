package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nameoffnv/time"
)

var (
	flagPort = flag.Int("p", 37, "Port to listening")
)

func main() {
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	done := make(chan struct{}, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	timeServer := time.New(time.RFCTime)

	go func() {
		<-sigs
		if err := timeServer.Close(); err != nil {
			fmt.Println(err)
		}
		close(done)
	}()

	if err := timeServer.Listen(fmt.Sprintf(":%d", *flagPort)); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	<-done
}
