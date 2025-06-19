package ops

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/acoshift/arpc/v2"
	"github.com/moonrhythm/parapet"
	"go.opencensus.io/trace"
)

func Recovery() parapet.Middleware {
	am := arpc.New()
	return parapet.MiddlewareFunc(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if e := recover(); e != nil {
					Report(e, r, "", debug.Stack())
					err, ok := e.(error)
					if !ok {
						err = fmt.Errorf("%v", e)
					}
					am.EncodeError(w, r, err)
				}
			}()
			h.ServeHTTP(w, r)
		})
	})
}

func InjectRequestIDToSpan() parapet.Middleware {
	return parapet.MiddlewareFunc(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Header.Get("X-Request-ID")
			span := trace.FromContext(r.Context())
			span.AddAttributes(trace.StringAttribute("request_id", reqID))
			h.ServeHTTP(w, r)
		})
	})
}
