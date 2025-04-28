package main

import (
	"backend-go/internal/api"
)

func main() {
	srv := api.NewApiServer()
	if err := srv.ListenAndServe("localhost:3000"); err != nil {
		panic(err)
	}
}
