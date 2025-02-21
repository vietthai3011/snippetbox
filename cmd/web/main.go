package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/vietthai3011/snippetbox/internal/models"
)

type application struct {
	logger        *slog.Logger
	snippet       *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")

	dsn := flag.String("dsn", "web:123456@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	// NewJSONHandler - tạo 1 logHandle
	// slog.New - tạo ra 1 logging struct có logHandle sử lý
	// HandlerOptions - mưc độ log
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dsn)

	if err != nil {
		logger.Error(err.Error())
	}

	templateCache, err := newTemplateCache()

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()
	app := &application{
		logger:        logger,
		snippet:       &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	logger.Info("starting server", "addr", *addr)

	err = http.ListenAndServe(*addr, app.routers())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}
