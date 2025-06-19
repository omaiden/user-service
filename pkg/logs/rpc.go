package logs

import (
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/acoshift/arpc/v2"

	"user-service/pkg/ops"
)

func methodFromPath(p string) string {
	p = strings.TrimSuffix(p, "/")
	if p == "" {
		return p
	}
	i := strings.LastIndex(p, "/")
	if i < 0 {
		i = 0
	}
	return strings.TrimPrefix(p[i:], "/")
}

type redactable interface {
	Redacted() any
}

func redacted(v any) any {
	if v, ok := v.(redactable); ok {
		return v.Redacted()
	}
	return v
}

func ReportRPCOK(w http.ResponseWriter, r *http.Request, req, res any) {
	m := methodFromPath(r.RequestURI)
	Info(r.Context(), S{
		Module:  "logs",
		Method:  m,
		Message: m + " rpc ok",
		Info: I{
			"method": r.Method,
			"path":   r.RequestURI,
			"params": redacted(req),
			"result": redacted(res),
		},
	})
}

func ReportRPCError(w http.ResponseWriter, r *http.Request, req any, err error) {
	m := methodFromPath(r.RequestURI)
	ctx := r.Context()
	s := S{
		Module:  "logs",
		Method:  m,
		Message: m,
		Error:   err.Error(),
		Info: I{
			"method": r.Method,
			"path":   r.RequestURI,
			"params": redacted(req),
		},
	}
	switch err.(type) {
	case arpc.OKError:
		s.Message += " rpc ok error"
		Warn(ctx, s)
	case *arpc.ProtocolError:
		s.Message += " rpc protocol error"
		s.Info["content_type"] = r.Header.Get("Content-Type")
		Warn(ctx, s)
	default:
		s.Message += " rpc internal error"
		Error(ctx, s)
		ops.Report(err, r, "", debug.Stack())
	}
}
