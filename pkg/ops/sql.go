package ops

import (
	"contrib.go.opencensus.io/integrations/ocsql"
	_ "github.com/lib/pq"
)

func PgDriver() string {
	pgDriver, _ := ocsql.Register("postgres", ocsql.WithQuery(true))
	return pgDriver
}
