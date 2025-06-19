package ops

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/errorreporting"
)

func Report(err interface{}, r *http.Request, user string, stack []byte) {
	if errorClient == nil {
		return
	}

	var e error
	if err, ok := err.(error); ok {
		e = err
	} else {
		e = fmt.Errorf("%v", err)
	}

	errorClient.Report(errorreporting.Entry{
		Error: e,
		Req:   r,
		User:  user,
		Stack: stack,
	})
}

func Reportf(format string, v ...interface{}) {
	Report(fmt.Errorf(format, v...), nil, "", nil)
}
