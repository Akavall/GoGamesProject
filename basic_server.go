package main

import (
        "net/http"
        "log"
        "fmt"
)

func dice_roll(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("Content-type", "text/plain")

    // Parse URL and POST data into the request.Form
    err := request.ParseForm()
    if err != nil {
        log.Fatal(response, fmt.Sprintf("error parsing url %v", err), 500)
    }

    my_dice := InitDefaultDice(6)
    side := my_dice.Roll()
    log.Printf("Rolled: %d for request...", side.numerical_value)

    // Actual response sent to web client
    fmt.Fprintf(response, "%d", side.numerical_value)
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", dice_roll)
    log.Printf("Started dumb Dice web server! Try it on http://localhost:8000")
    err := http.ListenAndServe(":8000", mux)
    
    if err != nil {
        log.Fatal(err)
    }
}
