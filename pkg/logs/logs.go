package logs

import (
	"context"

	"go.opencensus.io/trace"

	"user-service/pkg/ops"
)

type S struct {
	Module  string `json:"module"`
	Method  string `json:"method"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
	// RequestID string `json:"requestId,omitempty"`
	TraceID  string `json:"traceId,omitempty"`
	RealIP   string `json:"realIp,omitempty"`
	ClientIP string `json:"clientIp,omitempty"`
	Info     I      `json:"info"`
}

type I map[string]any

func opsLog(ctx context.Context, severity ops.Severity, s S) {
	// s.RequestID = reqid.Get(ctx)
	s.TraceID = trace.FromContext(ctx).SpanContext().TraceID.String()
	s.RealIP, _ = Get(ctx, "realIp").(string)
	s.ClientIP, _ = Get(ctx, "clientIp").(string)

	ops.Log(severity, s)
}

func Debug(ctx context.Context, s S) {
	opsLog(ctx, ops.Debug, s)
}

func Info(ctx context.Context, s S) {
	opsLog(ctx, ops.Info, s)
}

func Warn(ctx context.Context, s S) {
	opsLog(ctx, ops.Warning, s)
}

func Error(ctx context.Context, s S) {
	opsLog(ctx, ops.Error, s)
}

func Critical(ctx context.Context, s S) {
	opsLog(ctx, ops.Critical, s)
}

type recordKey struct{}

type Record map[string]any

func getRecord(ctx context.Context) Record {
	r, _ := ctx.Value(recordKey{}).(Record)
	return r
}

func Set(ctx context.Context, key string, value any) {
	r := getRecord(ctx)
	if r == nil {
		return
	}
	r[key] = value
}

func Get(ctx context.Context, key string) any {
	return getRecord(ctx)[key]
}
