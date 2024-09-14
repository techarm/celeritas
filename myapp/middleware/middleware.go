package middleware

import (
	"myapp/data"

	"github.com/techarm/celeritas"
)

type Middleware struct {
	App    *celeritas.Celeritas
	Models data.Models
}
