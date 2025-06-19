package user_test

import (
	"testing"

	"user-service/backoffice/user"
	"user-service/pkg/tu"

	"github.com/stretchr/testify/assert"
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
