module myapp

go 1.22.5

replace github.com/techarm/celeritas => ../celeritas

require github.com/techarm/celeritas v0.0.0-00010101000000-000000000000

require (
	github.com/go-chi/chi/v5 v5.1.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
)
