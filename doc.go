/*
Package healthcheck - universal HTTP healthcheck package

Usage example

First, create file with healthcheck configuration

  // healthcheck.yml
	---
	healthcheck:
		listen: ':9000'

Then you can define healthchecks

	package main

	import (
		"git.aqq.me/go/app/appconf"
		"git.aqq.me/go/app/launcher"
		"github.com/iph0/conf/fileconf"
		"github.com/kak-tus/healthcheck"
	)

	func init() {
		fileLdr, err := fileconf.NewLoader("etc")
		if err != nil {
			panic(err)
		}

		appconf.RegisterLoader("file", fileLdr)
	}

	type check struct {
		path string
	}

	func (c check) Check() (healthcheck.State, string) {
		return healthcheck.StatePassing, "ok " + c.path
	}

	func main() {
		launcher.Run(func() error {
			healthcheck.AddCheck("/ping", check{path: "/ping"})
			healthcheck.AddCheck("/status", check{path: "/status"})
			return nil
		})
	}
*/

package healthcheck
