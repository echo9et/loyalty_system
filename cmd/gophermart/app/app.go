package app

import (
	storage "gophermart.ru/internal/infrastructure/storage"
	interfaces "gophermart.ru/internal/interfaces/http"
)

type App struct {
	server  *interfaces.Server
	storage *storage.Database
}

func New() (*App, error) {
	app := &App{}

	var err error
	if app.storage, err = storage.NewDatabase("host=localhost user=username password=123321 dbname=echo9et sslmode=disable"); err != nil {
		return nil, err
	}

	if app.server, err = interfaces.New(); err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) Start(sAddress string) {
	a.server.Engine.Run(sAddress)
}
