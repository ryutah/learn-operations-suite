package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")

type entry struct {
	Severity string `json:"severity"`
	Message  string `json:"message"`
	Trace    string `json:"logging.googleapis.com/trace"`
}

func (e entry) String() string {
	if e.Severity == "" {
		e.Severity = "INFO"
	}
	v, _ := json.Marshal(e)
	return string(v)
}

func newEntry(severity, message, trace string) entry {
	return entry{
		Severity: severity,
		Message:  message,
		Trace:    trace,
	}
}

func (e entry) toJSON() string {
	v, _ := json.Marshal(e)
	return string(v)
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	http.HandleFunc("/", handle)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	_ = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	traceID := strings.SplitN(r.Header.Get("X-Cloud-Trace-Context"), "/", 2)[0]
	traceID = fmt.Sprintf("projects/%s/traces/%s", projectID, traceID)

	status := r.FormValue("status")
	statusCode, err := strconv.Atoi(status)
	if err != nil {
		log.Println(entry{
			Severity: "WARNING",
			Message:  fmt.Sprintf("status %q can not parse to integer", status),
			Trace:    traceID,
		})
		statusCode = http.StatusOK
	}

	log.Println(entry{
		Severity: "INFO",
		Message:  "this is info logging by stdout!!",
		Trace:    traceID,
	})
	log.Println(entry{
		Severity: "WARNING",
		Message:  "this is warn logging by stdout!!",
		Trace:    traceID,
	})
	log.Println(entry{
		Severity: "ALERT",
		Message:  "this is alert logging by stdout!!",
		Trace:    traceID,
	})

	w.WriteHeader(statusCode)
}
