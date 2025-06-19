package ops

import (
	"cloud.google.com/go/profiler"
)

func StartProfiler() error {
	return profiler.Start(profiler.Config{
		Service:   serviceName,
		ProjectID: projectID,
	})
}
