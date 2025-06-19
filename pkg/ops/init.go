package ops

import (
	"context"
	"log"
	"strconv"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/logging"
	"github.com/acoshift/configfile"
	"github.com/moonrhythm/parapet/pkg/stackdriver"
	"go.opencensus.io/trace"
	"golang.org/x/oauth2/google"
	mrpb "google.golang.org/genproto/googleapis/api/monitoredres"
)

var (
	projectID   string
	serviceName = "user-service"
	logToStd    = false
	traceSample = 1.0
)

var (
	logClient   *logging.Client
	logWriter   *logging.Logger
	errorClient *errorreporting.Client
)

func Init(ctx context.Context) {
	cfg := configfile.NewReader("config.yaml")

	projectID = cfg.String("ops_project")
	serviceName = cfg.StringDefault("ops_service", serviceName)
	logToStd = cfg.Bool("ops_log_to_std")
	traceSample, _ = strconv.ParseFloat(cfg.StringDefault("ops_trace_sample", "1.0"), 64)

	if logToStd {
		StartLogFlusher()
	}

	if projectID == "" {
		cred, _ := google.FindDefaultCredentials(ctx)
		if cred != nil {
			projectID = cred.ProjectID
		}
	}

	if projectID == "" {
		log.Println("ops: module disabled")
		return
	}

	log.Println("ops: module enabled")

	logClient, _ = logging.NewClient(ctx, "projects/"+projectID)
	if logClient != nil {
		log.Println("ops: logging enabled")
		logWriter = logClient.Logger(
			serviceName,
			logging.CommonResource(&mrpb.MonitoredResource{
				Type: "global",
				Labels: map[string]string{
					"project_id": projectID,
				},
			}),
			logging.ContextFunc(func() (ctx context.Context, afterCall func()) {
				ctx, span := trace.StartSpan(context.Background(), "", trace.WithSampler(trace.NeverSample()))
				return ctx, span.End
			}),
		)
	} else {
		log.Println("ops: logging disabled")
	}

	errorClient, _ = errorreporting.NewClient(ctx, projectID, errorreporting.Config{
		ServiceName: serviceName,
	})
	if errorClient != nil {
		log.Println("ops: error report enabled")
	} else {
		log.Println("ops: error report disabled")
	}

	if cfg.Bool("ops_profiler") {
		err := StartProfiler()
		if err != nil {
			log.Println("ops: profiler disabled")
			Logf(Warning, "ops: cannot start profiler")
		} else {
			log.Println("ops: profiler enabled")
		}
	}

	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.ProbabilitySampler(traceSample),
	})
	stackdriver.Register(stackdriver.Options{
		ProjectID:                projectID,
		TraceSpansBufferMaxBytes: 32 * 1024 * 1024,
	})
}

func Close() {
	if errorClient != nil {
		errorClient.Close()
	}
	if logClient != nil {
		logClient.Close()
	}
}
