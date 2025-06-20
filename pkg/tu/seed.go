package tu

import (
	"testing"
	"user-service/backoffice/user"

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

func (ctx *Context) CreateArticles(t *testing.T, article *user.Article) {
	t.Helper()

	_, err := pgctx.Exec(ctx.Ctx(), `
			INSERT INTO articles (id, title, content, author_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, now(), now())
		`, article.ID, article.Title, article.Content, article.AuthorID)
	require.NoError(t, err)
}
