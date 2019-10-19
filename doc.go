/*
Library for Consul HTTP healthcheck integration in go codee

Usage example

  hlth := healthcheck.NewHandler()

  hlth.Add("/ping", func() (healthcheck.State, string) {
    return healthcheck.StatePassing, "ok"
  })
  hlth.Add("/piiiing", func() (healthcheck.State, string) {
    return healthcheck.StateCritical, "not ok"
  })

  go http.ListenAndServe("0.0.0.0:9000", hlth)
*/
package healthcheck
