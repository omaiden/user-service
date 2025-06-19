package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchemaFS(t *testing.T) {
	t.Parallel()

	xs, err := fs.ReadDir(".")
	assert.NoError(t, err)
	assert.True(t, len(xs) > 0)
}
