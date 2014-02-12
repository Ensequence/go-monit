## go-monit

package monit provides a metric reporting mechanism for webapps

Basic usage:

```go
m = monit.NewMonitor(monit.Config{
	Host: "https://myhost.com/reporting/",
	Base: map[string]interface{}{
		"auth": "maybeINeedThis?"
	},
})
m.Start()
```
