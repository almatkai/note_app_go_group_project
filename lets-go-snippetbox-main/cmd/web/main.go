package main

import (
	"context"
	"crypto/tls"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"snippetbox.eegurt.net/internal/models"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type application struct {
	infoLog        *log.Logger
	errorLog       *log.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("addr", "localhost:4000", "HTTP network address") //
	dsn := flag.String("dsn", "postgres://web:789@localhost:5432/snippetbox", "connection login to database")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)                  //New Custom Logger INFO
	errorLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile) //New Custom Logger ERROR

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

	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		infoLog:        infoLog,
		errorLog:       errorLog,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	db, err := pgxpool.Connect(context.Background(), dsn)
	//ctx := context.Background()
	//Эта функция возвращает пустой контекст.
	//Она должна использоваться только на высоком уровне
	//(в main или обработчике запросов высшего уровня).
	//Он может быть использован для получения других контекстов,
	//которые мы обсудим позже.
	if err != nil {
		return nil, err
	}
	return db, err
}
