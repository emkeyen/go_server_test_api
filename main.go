package main

import (
	"log"
	"net/http"

	"github.com/emkeyen/go_server_test_api/httpserver"
)

func main() {
	// init with test data
	httpserver.Mu.Lock()
	httpserver.Users[1] = httpserver.User{ID: 1, Name: "Test User1"}
	httpserver.NextID = 2
	httpserver.Mu.Unlock()

	// register handlers
	http.HandleFunc("/", httpserver.GetRoot)
	http.HandleFunc("/hello", httpserver.GetHello)
	http.HandleFunc("/user", httpserver.UserHandler)

	log.Println("Starting server on :3333")
	log.Fatal(http.ListenAndServe(":3333", nil))
}
