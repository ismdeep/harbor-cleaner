package harbor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_ListRepoTags(t *testing.T) {
	tags, err := client.ListRepoTags("public", "debian")
	assert.NoError(t, err)
	for i, tag := range tags {
		t.Logf("tags[%v] = %+v", i, tag)
	}
}

func TestClient_DeleteTag(t *testing.T) {
	//err := client.DeleteTag("public", TagInfo{
	//	Repository: "debian",
	//	Name:       "13-fix",
	//	Digest:     "sha256:09d27af7a1f8cc4ab597de642c019511229361558685cf07809b732cb988f87b",
	//	Image:      "docker.example.com/public/debian",
	//})
	//assert.NoError(t, err)
}
