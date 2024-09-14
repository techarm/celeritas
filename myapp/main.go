package main

import (
	"myapp/data"
	"myapp/handlers"
	"myapp/middleware"

	"github.com/techarm/celeritas"
)

type application struct {
	App        *celeritas.Celeritas
	Middleware *middleware.Middleware
	Handlers   *handlers.Handlers
	Models     data.Models
}

func main() {
	c := initApplication()
	c.App.ListenAndServe()
}
