package main

import (
	"encoding/json"
	"go.uber.org/zap"
	"log"
	"net/http"
	"net/url"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)


type Response struct {
	Headers http.Header    `json:"headers"`
	Body    map[string]any `json:"body"`
	Path    string         `json:"path"`
	Params  url.Values     `json:"params"`
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET")
}

func handler(w http.ResponseWriter, r *http.Request) {
	logger := loggerWithCorrelateTrace(r.Context())

	logger.Info(
		"request received",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.Any("params", r.URL.Query()),
	)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	var body map[string]any
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		body = map[string]any{}
	}

	enableCors(&w)

	res := Response{
		Headers: r.Header,
		Body:    body,
		Path:    r.URL.Path,
		Params:  r.URL.Query(),
	}

	json.NewEncoder(w).Encode(res)
	logger.Info(
		"response sent",
		zap.Int("status", http.StatusOK),
		zap.Any("response", res),
	)
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8080"
	}

	instrumentedHandler := otelhttp.NewHandler(
		http.HandlerFunc(handler),
		"http-server",
	)

	logger.Info("server started", zap.String("port", port))
	http.ListenAndServe(":"+port, instrumentedHandler)
}
