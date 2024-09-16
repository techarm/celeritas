package mailer

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var pool *dockertest.Pool
var resource *dockertest.Resource

var mailer = Mail{
	Domain:      "localhost",
	Templates:   "./testdata/mail",
	Host:        "localhost",
	Port:        1027,
	Encryption:  "none",
	FromAddress: "me@here.com",
	FromName:    "Joe",
	Jobs:        make(chan Message, 1),
	Result:      make(chan Result, 1),
}

func TestMain(m *testing.M) {
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal("could not connect to docker", err)
	}
	pool = p

	opts := dockertest.RunOptions{
		Repository:   "mailhog/mailhog",
		Tag:          "latest",
		Env:          []string{},
		ExposedPorts: []string{"1025", "8025"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"1025": {
				{HostIP: "0.0.0.0", HostPort: "1027"},
			},
			"8025": {
				{HostIP: "0.0.0.0", HostPort: "8027"},
			},
		},
	}

	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		log.Fatalf("could not start resource: %s", err)
		err = pool.Purge(resource)
		if err != nil {
			log.Fatalf("count not purge resource: %s", err)
		}
	}

	time.Sleep(2 * time.Second)

	go mailer.ListenFromMail()
	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}
