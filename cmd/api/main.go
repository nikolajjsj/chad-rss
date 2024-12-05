package main

import (
	"chad-rss/internal/jobs"
	"chad-rss/internal/server"
	"fmt"
)

func main() {
	server := server.NewServer()

	go jobs.RunJobs()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
