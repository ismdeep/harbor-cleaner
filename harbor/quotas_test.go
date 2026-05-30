package harbor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_Quotas(t *testing.T) {
	quotas, err := client.Quotas()
	assert.NoError(t, err)
	for i, quota := range quotas {
		t.Logf("%d: %+v", i, quota)
	}
}
