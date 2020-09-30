package apiserver

import (
	"database/sql"
	"gismart-rest-api/internal/app/store/sqlstore"
	"github.com/gorilla/sessions"
	"net/http"
	"strings"
)

func Start(config *Config, user string, password string) error {
	databaseURL := config.DatabaseURL
	if user != "" {
		params := []string{user, password, config.DatabaseURL}
		databaseURL = strings.Join(params, " ")
	}
	db, err := newDB(databaseURL)
	if err != nil {
		return err
	}
	defer db.Close()
	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	srv := newServer(store, sessionStore)

	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}
	if  err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}