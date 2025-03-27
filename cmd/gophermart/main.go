package main

import (
	"fmt"
	"log/slog"

	// "github.com/echo9et/loyalty_system/cmd/gophermart/app"
	"gophermart.ru/cmd/gophermart/app"
	config "gophermart.ru/internal"
)

func main() {
	app, err := app.New()
	if err != nil {
		slog.Error(fmt.Sprintln(err))
		return
	}
	err = app.Start(config.Get().AddrServer)

	if err != nil {
		slog.Error(err.Error())
	}
}
