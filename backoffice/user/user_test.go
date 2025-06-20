package user_test

import (
	"context"
	"testing"
	"time"

	"user-service/backoffice/internal/kctx"
	"user-service/backoffice/user"
	"user-service/pkg/tu"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	tc := tu.Setup()
	defer tc.Teardown()

	params := &user.CreateUserRequest{
		Name:  "Thunder",
		Email: "thunder@example.com",
	}

	createdUser, err := user.CreateUser(tc.Ctx(), params)
	assert.NoError(t, err)
	assert.Equal(t, "Thunder", createdUser.Name)
	assert.Equal(t, "thunder@example.com", createdUser.Email)
}

func TestGetUserArticles(t *testing.T) {
	tc := tu.Setup()
	defer tc.Teardown()

	now := time.Now()
	userID := "test-user"
	authorID := "test-author"
	ctx := tc.Ctx()
	ctx = kctx.NewUserIDContext(ctx, userID)

	// Case 1: Invalid limit
	t.Run("bad request - invalid limit", func(t *testing.T) {
		_, err := user.GetUserArticles(ctx, authorID, 0)
		assert.ErrorContains(t, err, "limit must be between 1 and 100")

		_, err = user.GetUserArticles(ctx, authorID, 101)
		assert.ErrorContains(t, err, "limit must be between 1 and 100")
	})

	// Case 2: No user ID in context
	t.Run("unauthorized - no user in context", func(t *testing.T) {
		_, err := user.GetUserArticles(context.Background(), authorID, 10)
		assert.ErrorContains(t, err, "unauthorized")
	})

	// Insert articles for testing
	for i := range 5 {
		tc.CreateArticles(t, &user.Article{
			ID:        i,
			Title:     "test",
			Content:   "test",
			AuthorID:  authorID,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	// Case 3: Successful data retrieval
	t.Run("success - all articles", func(t *testing.T) {
		articles, err := user.GetUserArticles(ctx, authorID, 10)
		require.NoError(t, err)
		assert.Len(t, articles, 5)
		for _, article := range articles {
			assert.Equal(t, authorID, article.AuthorID)
		}
	})

	// Case 4: No articles found
	t.Run("success - no articles found", func(t *testing.T) {
		articles, err := user.GetUserArticles(ctx, "non-article-user", 10)
		require.NoError(t, err)
		assert.Len(t, articles, 0)
	})

	// Case 5: Limit works correctly
	t.Run("success - limit works", func(t *testing.T) {
		articles, err := user.GetUserArticles(ctx, authorID, 3)
		require.NoError(t, err)
		assert.Len(t, articles, 3)
	})
}
