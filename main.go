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
)

func main() {
    var (
        httpAddr = flag.String("http", ":3000", "http listen address")
    )
    flag.Parse()
    ctx := context.Background()
    // our napodate service
    srv := NewService()
    errChan := make(chan error)

    go func() {
        c := make(chan os.Signal, 1)
        signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
        errChan <- fmt.Errorf("%s", <-c)
    }()

    // mapping endpoints
    endpoints := Endpoints{
        GetEndpoint:      MakeGetEndpoint(srv),
        StatusEndpoint:   MakeStatusEndpoint(srv),
        ValidateEndpoint: MakeValidateEndpoint(srv),
    }

    // HTTP transport
    go func() {
        log.Println("napodate is listening on port:", *httpAddr)
        handler := NewHTTPServer(ctx, endpoints)
        errChan <- http.ListenAndServe(*httpAddr, handler)
    }()

    log.Fatalln(<-errChan)
}
