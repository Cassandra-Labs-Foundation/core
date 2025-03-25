package main

import (
    "log"
    
    tb "github.com/tigerbeetle/tigerbeetle-go"
    "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

func main() {
    // Create the client with cluster ID 0
    client, err := tb.NewClient(
        types.Uint128{0, 0},
        []string{"34.53.7.189:3000"},
    )
    if err != nil {
        log.Fatalf("Error creating client: %v", err)
    }
    defer client.Close()

    log.Println("TigerBeetle client created successfully")
}