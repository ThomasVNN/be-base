// +build !windows
// +build !plan9
// +build !nacl

package syslog_test

import (
	"fmt"

	gosyslog "log/syslog"

	"git-codecommit.ap-southeast-1.amazonaws.com/v1/repos/be-base/log"
	"git-codecommit.ap-southeast-1.amazonaws.com/v1/repos/be-base/log/level"
	"git-codecommit.ap-southeast-1.amazonaws.com/v1/repos/be-base/log/syslog"
)

func ExampleNewSyslogLogger_defaultPrioritySelector() {
	// Normal syslog writer
	w, err := gosyslog.New(gosyslog.LOG_INFO, "experiment")
	if err != nil {
		fmt.Println(err)
		return
	}

	// syslog logger with logfmt formatting
	logger := syslog.NewSyslogLogger(w, log.NewLogfmtLogger)
	logger.Log("msg", "info because of default")
	logger.Log(level.Key(), level.DebugValue(), "msg", "debug because of explicit level")
}
