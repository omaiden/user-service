package kctx

import "context"

type ctxKeyUserID struct{}

func NewUserIDContext(parent context.Context, userID string) context.Context {
	return context.WithValue(parent, ctxKeyUserID{}, userID)
}

func GetUserID(ctx context.Context) string {
	v, _ := ctx.Value(ctxKeyUserID{}).(string)
	return v
}
