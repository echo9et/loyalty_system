package main

import (
	"fmt"
	"log/slog"

	// "github.com/echo9et/loyalty_system/cmd/gophermart/app"
	"gophermart.ru/cmd/gophermart/app"
)

func main() {
	app, err := app.New()
	if err != nil {
		slog.Error(fmt.Sprintln(err))
		return
	}
	err = app.Start("127.0.0.1:8000")

	if err != nil {
		slog.Error(err.Error())
	}
}
