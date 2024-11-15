package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/d4ny4z0rd/godogo/router"
)

func main() {
	r := router.Router()
	fmt.Println("Starting the server on port 3000")

	log.Fatal(http.ListenAndServe(":3000",r))
}
