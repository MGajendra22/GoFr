package main

import (
	"github.com/MGajendra22/GoFr/migrations"
	"gofr.dev/pkg/gofr"
)

type Task struct {
	ID     int    `json:"id"`
	Desc   string `json:"desc"`
	Status bool   `json:"status"`
	Userid int    `json:"userid"`
}

func main() {

	app := gofr.New()

	app.Migrate(migrations.All())

	err := app.AddRESTHandlers(&Task{})
	if err != nil {
		return
	}

	app.Run()

}
