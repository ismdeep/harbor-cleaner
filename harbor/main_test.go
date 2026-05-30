package harbor

import (
	"context"
	"os"
	"testing"
)

var (
	endpoint = "https://docker.example.com"
	username = "admin"
	password = "Harbor12345"
)

var client *Client

func TestMain(m *testing.M) {
	if os.Getenv("HARBOR_ENDPOINT") != "" {
		endpoint = os.Getenv("HARBOR_ENDPOINT")
	}
	if os.Getenv("HARBOR_USERNAME") != "" {
		username = os.Getenv("HARBOR_USERNAME")
	}
	if os.Getenv("HARBOR_PASSWORD") != "" {
		password = os.Getenv("HARBOR_PASSWORD")
	}

	var err error
	client, err = NewClient(context.Background(), endpoint, username, password, 4)
	if err != nil {
		panic(err)
	}

	m.Run()
}
