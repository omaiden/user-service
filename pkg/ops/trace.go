package ops

import (
	"github.com/moonrhythm/parapet"
	"github.com/moonrhythm/parapet/pkg/stackdriver"
)

func Trace() parapet.Middleware {
	return stackdriver.Trace()
}
