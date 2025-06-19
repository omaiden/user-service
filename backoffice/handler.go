package backoffice

import (
	"github.com/acoshift/arpc/v2"
	"github.com/moonrhythm/httpmux"
	"user-service/backoffice/user"
)

func Mount(mux *httpmux.Mux, am *arpc.Manager) {
	mux.Handle("POST /users", am.Handler(user.CreateUser))
}
