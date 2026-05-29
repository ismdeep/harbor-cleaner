package harbor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	ctx := context.Background()
	clientNew, err := NewClient(ctx, endpoint, username, password)
	assert.NoError(t, err)
	_ = clientNew
}
