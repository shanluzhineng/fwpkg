package gorminit

import (
	"fmt"

	"github.com/shanluzhineng/fwpkg/system/log"

	"gorm.io/gorm/logger"
)

type writer struct {
	logger.Writer
	LogZap bool
}

func newWriter(logZap bool, w logger.Writer) *writer {
	return &writer{
		Writer: w,
		LogZap: logZap,
	}
}

func (w *writer) printf(message string, data ...interface{}) {
	var logZap bool = w.LogZap
	if logZap {
		log.Logger.Info(fmt.Sprintf(message+"\n", data...))
	} else {
		w.Writer.Printf(message, data...)
	}
}
