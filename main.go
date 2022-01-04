package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	DateTest "DateTestModule/DateTestFolder"
)

func main() {
	var (
		httpAddr = flag.String("http", ":3000", "http listen address")
	)
	flag.Parse()
	ctx := context.Background()
	// our napodate service
	srv := DateTest.NewService()
	errChan := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	// mapping endpoints
	endpoints := DateTest.Endpoints{
		GetEndpoint:      DateTest.MakeGetEndpoint(srv),
		StatusEndpoint:   DateTest.MakeStatusEndpoint(srv),
		ValidateEndpoint: DateTest.MakeValidateEndpoint(srv),
	}

	// HTTP transport
	go func() {
		log.Println("DateTest is listening on port:", *httpAddr)
		handler := DateTest.NewHTTPServer(ctx, endpoints)
		errChan <- http.ListenAndServe(*httpAddr, handler)
	}()

	log.Fatalln(<-errChan)
}
