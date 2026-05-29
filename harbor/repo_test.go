package harbor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_ListRepositories(t *testing.T) {
	repositories, err := client.ListRepositories("library")
	assert.NoError(t, err)
	for i, repository := range repositories {
		t.Logf("repositories[%v] = %v\n", i, repository)
	}
}
