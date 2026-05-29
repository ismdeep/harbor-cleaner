package harbor

import (
	"context"
	"testing"
)

var (
	endpoint = "https://docker.example.com"
	username = "admin"
	password = "Harbor12345"
)

var client *Client

func TestMain(m *testing.M) {
	var err error
	client, err = NewClient(context.Background(), endpoint, username, password)
	if err != nil {
		panic(err)
	}

	m.Run()
}
