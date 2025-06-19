package hook

import (
	"context"
)

type ctxKey struct{}

type Func func(v any)

type registry struct {
	hooks map[string][]Func
}

func NewContext(parent context.Context) context.Context {
	var reg registry
	reg.hooks = make(map[string][]Func)
	return context.WithValue(parent, ctxKey{}, &reg)
}

func Register(ctx context.Context, event string, hook Func) {
	reg := ctx.Value(ctxKey{}).(*registry)
	reg.hooks[event] = append(reg.hooks[event], hook)
}

func Hook(ctx context.Context, event string, v any) {
	reg := ctx.Value(ctxKey{})
	if reg == nil {
		return
	}
	for _, hook := range reg.(*registry).hooks[event] {
		hook(v)
	}
}
