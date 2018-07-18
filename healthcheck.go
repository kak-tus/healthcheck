package healthcheck

import (
	"fmt"
	"net/http"

	"git.aqq.me/go/app/appconf"
	"git.aqq.me/go/app/applog"
	"git.aqq.me/go/app/event"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

type healthcheckConfig struct {
	Listen string
}

type server struct {
	logger   *zap.SugaredLogger
	listener *http.Server
}

// State type
type State int

// StatePassing - check in passing state
const StatePassing State = 200

// StateWarning - check in warning state
const StateWarning State = 429

// StateCritical - check in critical state
const StateCritical State = 500

// Check implements healthcheck function
type Check interface {
	Check() (State, string)
}

var srv *server

func init() {
	appconf.Require("file:healthcheck.yml")

	event.Init.AddHandler(
		func() error {
			cnf := appconf.GetConfig()["healthcheck"]

			var config healthcheckConfig
			err := mapstructure.Decode(cnf, &config)
			if err != nil {
				return err
			}

			logger := applog.GetLogger()

			srv = &server{
				logger: logger,
				listener: &http.Server{
					Addr: config.Listen,
				},
			}

			go func() {
				err = srv.listener.ListenAndServe()
				if err != nil && err != http.ErrServerClosed {
					logger.Error(err)
				}
			}()

			return nil
		},
	)

	event.Stop.AddHandler(
		func() error {
			err := srv.listener.Shutdown(nil)
			if err != nil {
				return err
			}
			return nil
		},
	)
}

// AddCheck add new HTTP check
func AddCheck(path string, check Check) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		srv.logger.Debug("Request ", path)

		state, text := check.Check()

		srv.logger.Debug("Response state: ", state)

		if state != StatePassing {
			w.WriteHeader(int(state))
		}

		_, err := fmt.Fprintf(w, text)
		if err != nil {
			srv.logger.Error(err)
		}
	})
}
