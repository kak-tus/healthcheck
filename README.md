# Library for Consul HTTP healthcheck integration in go code

[Documentation](https://godoc.org/github.com/kak-tus/healthcheck).

Package realise a little snippets to simplify Consul healthcheck integration.

I planed to add Kubernetes support, but there is a popular
[healthcheck](https://github.com/heptiolabs/healthcheck) library for Kubernetes.

So I delegate this library to only Consul support.

## Example

```
  hlth := healthcheck.NewHandler()

  hlth.Add("/ping", func() (healthcheck.State, string) {
    return healthcheck.StatePassing, "ok"
  })
  hlth.Add("/piiiing", func() (healthcheck.State, string) {
    return healthcheck.StateCritical, "not ok"
  })

  go http.ListenAndServe("0.0.0.0:9000", hlth)
```
