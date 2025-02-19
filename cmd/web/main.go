package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// NewJSONHandler - tạo 1 logHandle
	// slog.New - tạo ra 1 logging struct có logHandle sử lý
	// HandlerOptions - mưc độ log
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{logger: logger}

	logger.Info("starting server", "addr", addr)

	err := http.ListenAndServe(*addr, app.routers())
	logger.Error(err.Error())
	os.Exit(1)
}
