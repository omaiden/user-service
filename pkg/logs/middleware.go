package logs

import (
	"context"
	"net/http"

	"github.com/moonrhythm/parapet"
)

func InjectRecord() parapet.Middleware {
	return parapet.MiddlewareFunc(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, recordKey{}, Record{})
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}
