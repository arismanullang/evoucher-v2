package x

import (
	"log"
	"os"
	"sync"
)

var (
	loggerMu sync.Mutex
	logger   = log.New(os.Stderr, "", log.LstdFlags)
)

func SetLogger(l *log.Logger) {
	loggerMu.Lock()
	logger = l
	loggerMu.Unlock()
}

func Panic(err error) {
	logger.Panic(err)
}

func PanicOnError(err error) {
	if err != nil {
		logger.Panic(err)
	}
}

func Fatal(err error) {
	logger.Fatal(err)
}

func FatalOnError(err error) {
	if err != nil {
		logger.Fatal(err)
	}
}
