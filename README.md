# logstasher

logstasher is a Negroni middleware that prints logstash-compatible JSON to an `io.Writer` for each HTTP request.

Here's an example from one of the Go microservices we have at @bikeexchange :

``` json
{
  "@timestamp":"2014-03-01T19:08:06+11:00","@version":1,"method":"GET",
  "path":"/locations/slugs/VIC/Williams-Landing","status":200,"size":238,
  "duration":14.059902000000001,"params":{"country":["au"]}
}
```

Used in conjunction with the [rotating file writer](http://github.com/mipearson/rfw) it allows for rotatable logs ready to feed directly into logstash with no parsing.

### Example

``` go
package main

import (
  "log"

  "github.com/codegangsta/negroni"
  "github.com/exchangegroup/logstasher"
  "github.com/mipearson/rfw"
)

func main() {
  n := negroni.Classic()

  logstashLogFile, err := rfw.Open("hello.log", 0644)
  if err != nil {
    log.Fatalln(err)
  }
  defer logstashLogFile.Close()
  n.Use(logstasher.Logger(logstashLogFile))

  n.Get("/", func() string {
    return "Hello world!"
  })
  n.Run()
}
```

```
## logstash.conf
input {
  file {
    path => ["hello.log"]
    codec => "json"
  }
}
```
