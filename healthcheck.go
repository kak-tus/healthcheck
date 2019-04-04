package healthcheck

import (
	"fmt"
	"net/http"

	"git.aqq.me/go/app/appconf"
	"git.aqq.me/go/app/applog"
	"git.aqq.me/go/app/event"
	"github.com/iph0/conf"
	"go.uber.org/zap"
)

type healthcheckConfig struct {
	Listen string
}

type server struct {
	log      *zap.SugaredLogger
	listener *http.Server
}

// State type
type State int

const (
	_ = iota
	// StatePassing - check in passing state
	StatePassing State = iota
	// StateWarning - check in warning state
	StateWarning State = iota
	// StateCritical - check in critical state
	StateCritical State = iota
)

var stateMap = map[State]int{
	1: 200,
	2: 429,
	3: 500,
}

var srv *server

func init() {
	event.Init.AddHandler(
		func() error {
			cnf := appconf.GetConfig()["healthcheck"]

			var config healthcheckConfig
			err := conf.Decode(cnf, &config)
			if err != nil {
				return err
			}

			srv = &server{
				log: applog.GetLogger().Sugar(),
				listener: &http.Server{
					Addr: config.Listen,
				},
			}

			go func() {
				err = srv.listener.ListenAndServe()
				if err != nil && err != http.ErrServerClosed {
					srv.log.Error(err)
				}
			}()

			srv.log.Info("Started healthcheck listener")

			return nil
		},
	)

	event.Stop.AddHandler(
		func() error {
			srv.log.Info("Stop healthcheck listener")

			err := srv.listener.Shutdown(nil)
			if err != nil {
				return err
			}

			srv.log.Info("Stopped healthcheck listener")

			return nil
		},
	)
}

// Add add new HTTP healthcheck
func Add(path string, f func() (State, string)) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		state, text := f()

		if state != StatePassing {
			w.WriteHeader(stateMap[state])
		}

		_, err := fmt.Fprintf(w, text)
		if err != nil {
			srv.log.Error(err)
		}
	})
}

// AddReq add new HTTP healthcheck with http.Request parameter
// to allow get some data from request
func AddReq(path string, f func(*http.Request) (State, string)) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		srv.log.Debug("Request ", path)

		state, text := f(r)

		srv.log.Debug("Response state: ", state)

		if state != StatePassing {
			w.WriteHeader(stateMap[state])
		}

		_, err := fmt.Fprintf(w, text)
		if err != nil {
			srv.log.Error(err)
		}
	})
}
