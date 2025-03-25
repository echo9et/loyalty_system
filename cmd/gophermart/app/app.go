package app

import (
	config "gophermart.ru/internal"
	storage "gophermart.ru/internal/infrastructure/storage"
	interfaces "gophermart.ru/internal/interfaces/http"
)

type App struct {
	server  *interfaces.Server
	storage *storage.Database
}

func New() (*App, error) {
	var err error

	if _, err = config.ParseFlags(); err != nil {
		return nil, err
	}

	app := &App{}

	if app.storage, err = storage.NewDatabase(config.Get().AddrDatabase); err != nil {
		return nil, err
	}

	if app.server, err = interfaces.New(); err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) Start(addr string) error {
	err := a.server.Engine.Run(addr)
	if err != nil {
		return err
	}
	return nil
}
