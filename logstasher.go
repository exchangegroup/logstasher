// Package logstasher is a Martini middleware that prints logstash-compatiable
// JSON to a given io.Writer for each HTTP request.
package logstasher

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/codegangsta/negroni"
)

type Logger struct {
	*bufio.Writer
}

type logstashEvent struct {
	Timestamp string              `json:"@timestamp"`
	Version   int                 `json:"@version"`
	Method    string              `json:"method"`
	Path      string              `json:"path"`
	Status    int                 `json:"status"`
	Size      int                 `json:"size"`
	Duration  float64             `json:"duration"`
	Params    map[string][]string `json:"params,omitempty"`
}

func New(w io.Writer) *Logger {
	return &Logger{bufio.NewWriter(w)}
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	start := time.Now()

	res := rw.(negroni.ResponseWriter)
	next(rw, req)

	event := logstashEvent{
		time.Now().Format(time.RFC3339),
		1,
		req.Method,
		req.URL.Path,
		res.Status(),
		res.Size(),
		time.Since(start).Seconds() * 1000.0,
		map[string][]string(req.Form),
	}

	output, err := json.Marshal(event)
	if err != nil {
		// Should this be fatal?
		log.Printf("Unable to JSON-ify our event (%#v): %v", event, err)
		return
	}

	l.Write(output)
	l.WriteByte('\n')
	err = l.Flush()
	if err != nil {
		log.Printf("logstasher could not write: %v", err)
	}
}
