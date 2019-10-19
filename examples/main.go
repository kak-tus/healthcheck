package main

import (
	"net/http"

	"github.com/kak-tus/healthcheck"
)

func main() {
	hlth := healthcheck.NewHandler()

	hlth.Add("/ping", func() (healthcheck.State, string) {
		return healthcheck.StatePassing, "ok"
	})
	hlth.Add("/piiiing", func() (healthcheck.State, string) {
		return healthcheck.StateCritical, "not ok"
	})

	err := http.ListenAndServe(":9000", hlth)
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
