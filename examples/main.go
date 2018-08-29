package main

import (
	"git.aqq.me/go/app/appconf"
	"git.aqq.me/go/app/launcher"
	"github.com/iph0/conf/fileconf"
	"github.com/kak-tus/healthcheck"
)

func init() {
	fileLdr := fileconf.NewLoader("etc")

	appconf.RegisterLoader("file", fileLdr)

	appconf.Require("file:app.yml")
}

func main() {
	launcher.Run(func() error {
		healthcheck.Add("/ping", func() (healthcheck.State, string) {
			return healthcheck.StatePassing, "ok"
		})
		healthcheck.Add("/status", func() (healthcheck.State, string) {
			return healthcheck.StateCritical, "err"
		})
		return nil
	})
}
