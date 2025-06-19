package tu

import (
	"testing"

	"github.com/acoshift/pgsql/pgctx"
	"github.com/moonrhythm/randid"
	"github.com/stretchr/testify/require"
)

func (ctx *Context) CreateUser(t *testing.T, name string) string {
	t.Helper()

	userID := randid.MustGenerate()
	_, err := pgctx.Exec(ctx.Ctx(), `
		insert into users (id, name)
		values ($1, $2)
	`, userID, name)
	require.NoError(t, err)
	return userID.String()
}
