package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	mrpb "google.golang.org/genproto/googleapis/api/monitoredres"

	"cloud.google.com/go/logging"
)

var logger *logging.Logger

var (
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	service   = os.Getenv("GAE_SERVICE")
	version   = os.Getenv("GAE_VERSION")
)

func main() {
	client, err := logging.NewClient(context.Background(), fmt.Sprintf("projects/%s", os.Getenv("GOOGLE_CLOUD_PROJECT")))
	if err != nil {
		panic(err)
	}
	defer client.Close()
	logger = client.Logger("app_logs", logging.CommonResource(&mrpb.MonitoredResource{
		Type: "gae_app",
		Labels: map[string]string{
			"module_id":  service,
			"project_id": projectID,
			"version_id": version,
		},
	}))

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

	logger.Log(logging.Entry{
		Severity: logging.Info,
		Payload:  "this is info logging!!!",
		Trace:    traceID,
	})
	logger.Log(logging.Entry{
		Severity: logging.Warning,
		Payload:  "this is info warining!!!",
		Trace:    traceID,
	})
	logger.Log(logging.Entry{
		Severity: logging.Error,
		Payload:  "this is info error!!!",
		Trace:    traceID,
	})
	logger.Log(logging.Entry{
		Severity: logging.Alert,
		Payload:  "this is info alert!!!",
		Trace:    traceID,
	})
}
