package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zrotrasukha/snippetbox/internal/modeles"
)

type config struct {
	addr string
}

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *modeles.SnippetModel
	users          *modeles.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	cfg := config{
		addr: ":4000",
	}

	dsn := flag.String("dsn", "web:pass@tcp(:33060)/snippetbox?parseTime=true", "MySQL data source name")

	flag.StringVar(&cfg.addr, "addr", cfg.addr, "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ltime|log.Ldate)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ltime|log.Lshortfile|log.Ldate)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionsManager := scs.New()
	sessionsManager.Store = mysqlstore.New(db)
	sessionsManager.Lifetime = 12 * time.Hour

	sessionsManager.Cookie.Secure = true

	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &modeles.SnippetModel{DB: db},
		users:          &modeles.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionsManager,
	}
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}
	srv := &http.Server{
		Addr:         cfg.addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", cfg.addr)
	err = srv.ListenAndServeTLS("tls/cert.pem", "tls/key.pem")
	errorLog.Fatal(err)
}
